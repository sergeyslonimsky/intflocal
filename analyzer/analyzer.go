package analyzer

import (
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/sergeyslonimsky/intflocal/analyzer/internal/config"
)

const name = "intflocal"

// Settings is the raw settings from golangci-lint configuration.
// Re-exported from internal config package for use by the plugin entry point.
type Settings = config.Settings

// New creates a new analysis.Analyzer with the given settings.
// When settings are empty (standalone mode), flags are used instead.
func New(s Settings) *analysis.Analyzer {
	a := &runner{
		settings: s,
	}

	analyzer := &analysis.Analyzer{
		Name:     name,
		Doc:      "checks that struct interface dependencies are defined locally as private interfaces",
		Run:      a.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	a.registerFlags(&analyzer.Flags)

	return analyzer
}

type runner struct {
	settings Settings
	once     sync.Once
	cfg      *config.Config

	moduleOnce sync.Once
	modulePath string // cached module path read from go.mod

	// CLI flags (used when running standalone, ignored when settings are provided via plugin).
	flagExcludePackages string
	flagExcludeTypes    string
	flagPackages        string
}

func (r *runner) registerFlags(fs *flag.FlagSet) {
	fs.StringVar(&r.flagExcludePackages, "excludePackages", "", "comma-separated list of packages to exclude (e.g. github.com/some/pkg,github.com/other/pkg)")
	fs.StringVar(&r.flagExcludeTypes, "excludeTypes", "", "comma-separated list of fully qualified type names to exclude (e.g. github.com/some/pkg.MyInterface)")
	fs.StringVar(&r.flagPackages, "packages", "", "comma-separated list of package patterns to check (e.g. ./internal/services/...,./pkg/handlers/...)")
}

func (r *runner) initConfig() {
	r.once.Do(func() {
		s := r.settings

		// Merge CLI flags into settings (flags take precedence when non-empty).
		if r.flagExcludePackages != "" {
			s.ExcludePackages = splitCSV(r.flagExcludePackages)
		}
		if r.flagExcludeTypes != "" {
			s.ExcludeTypes = splitCSV(r.flagExcludeTypes)
		}
		if r.flagPackages != "" {
			s.Packages = splitCSV(r.flagPackages)
		}

		r.cfg = config.New(s)
	})
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func (r *runner) run(pass *analysis.Pass) (any, error) {
	r.initConfig()

	// Check if this package should be analyzed.
	if !r.cfg.ShouldCheckPackage(pass.Pkg.Path(), r.resolveModulePath(pass)) {
		return nil, nil
	}

	ins, _ := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	ins.Preorder([]ast.Node{(*ast.StructType)(nil)}, func(node ast.Node) {
		st, _ := node.(*ast.StructType)
		if st.Fields == nil {
			return
		}

		for _, field := range st.Fields.List {
			r.checkField(pass, field)
		}
	})

	return nil, nil
}

func (r *runner) checkField(pass *analysis.Pass, field *ast.Field) {
	typ := pass.TypesInfo.TypeOf(field.Type)
	if typ == nil {
		return
	}

	// Unwrap all pointer levels: *T, **T, etc.
	for {
		ptr, ok := typ.(*types.Pointer)
		if !ok {
			break
		}
		typ = ptr.Elem()
	}

	// We only care about named types.
	named, ok := typ.(*types.Named)
	if !ok {
		return
	}

	// Check if the underlying type is an interface.
	if _, ok := named.Underlying().(*types.Interface); !ok {
		return
	}

	obj := named.Obj()
	pkg := obj.Pkg()

	// Builtin types (e.g. error) have nil package.
	if pkg == nil {
		return
	}

	// Local interface — defined in the same package.
	if pkg.Path() == pass.Pkg.Path() {
		return
	}

	// Stdlib interfaces are always allowed.
	if isStdlib(pkg.Path()) {
		return
	}

	// Check exclude lists.
	if r.cfg.IsExcludedPackage(pkg.Path()) {
		return
	}

	fullName := pkg.Path() + "." + obj.Name()
	if r.cfg.IsExcludedType(fullName) {
		return
	}

	ifaceName := pkg.Name() + "." + obj.Name()

	// Named fields: report a separate diagnostic for each name so that
	// "a, b extiface.Iface" produces two distinct findings.
	if len(field.Names) > 0 {
		for _, name := range field.Names {
			pass.Report(analysis.Diagnostic{
				Pos:     name.Pos(),
				End:     field.Type.End(),
				Message: fmt.Sprintf("struct field %q uses external interface %q; define it locally as a private interface", name.Name, ifaceName),
			})
		}
		return
	}

	// Embedded field (anonymous).
	pass.Report(analysis.Diagnostic{
		Pos:     field.Type.Pos(),
		End:     field.Type.End(),
		Message: fmt.Sprintf("struct field \"(embedded)\" uses external interface %q; define it locally as a private interface", ifaceName),
	})
}

// isStdlib reports whether pkgPath belongs to the Go standard library.
// Stdlib packages never contain a dot in the first path segment (e.g. "fmt", "net/http", "internal/reflectlite"),
// while third-party packages always do (e.g. "github.com/user/pkg").
func isStdlib(pkgPath string) bool {
	first := pkgPath
	if i := strings.IndexByte(pkgPath, '/'); i >= 0 {
		first = pkgPath[:i]
	}

	return !strings.Contains(first, ".")
}

// resolveModulePath returns the module path for the package being analyzed.
// It reads go.mod from the filesystem (cached after the first call) and falls
// back to a best-effort heuristic when go.mod is not found (e.g. in tests).
func (r *runner) resolveModulePath(pass *analysis.Pass) string {
	r.moduleOnce.Do(func() {
		r.modulePath = findModulePath(pass)
	})
	if r.modulePath != "" {
		return r.modulePath
	}
	return extractModulePath(pass.Pkg.Path())
}

// findModulePath walks up the directory tree from the first source file of the
// package until it finds a go.mod, then returns the declared module path.
func findModulePath(pass *analysis.Pass) string {
	if len(pass.Files) == 0 {
		return ""
	}
	filename := pass.Fset.Position(pass.Files[0].Pos()).Filename
	dir := filepath.Dir(filename)
	for {
		data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
		if err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				line = strings.TrimSpace(line)
				if after, ok := strings.CutPrefix(line, "module "); ok {
					return strings.TrimSpace(after)
				}
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// extractModulePath is a fallback heuristic used when go.mod is not accessible
// (e.g. in analysistest). It assumes a three-segment module path which covers
// the common "host.tld/user/repo" pattern.
func extractModulePath(pkgPath string) string {
	parts := strings.SplitN(pkgPath, "/", 4)
	if len(parts) >= 3 && strings.Contains(parts[0], ".") {
		return strings.Join(parts[:3], "/")
	}
	return pkgPath
}
