package config

import "strings"

// Config holds parsed and validated linter configuration.
type Config struct {
	packages        []string
	excludeTypes    map[string]struct{}
	excludePackages map[string]struct{}
}

// Settings is the raw settings from golangci-lint configuration.
type Settings struct {
	Packages        []string `json:"packages"`
	ExcludeTypes    []string `json:"excludeTypes"`
	ExcludePackages []string `json:"excludePackages"`
}

// New creates a Config from raw Settings.
func New(s Settings) *Config {
	excludeTypes := make(map[string]struct{}, len(s.ExcludeTypes))
	for _, t := range s.ExcludeTypes {
		excludeTypes[t] = struct{}{}
	}

	excludePackages := make(map[string]struct{}, len(s.ExcludePackages))
	for _, p := range s.ExcludePackages {
		excludePackages[p] = struct{}{}
	}

	return &Config{
		packages:        s.Packages,
		excludeTypes:    excludeTypes,
		excludePackages: excludePackages,
	}
}

// ShouldCheckPackage reports whether the package with the given path should be analyzed.
// modulePath is the module path from go.mod (used to resolve relative package patterns).
// If no packages are configured, all packages are checked.
func (c *Config) ShouldCheckPackage(pkgPath, modulePath string) bool {
	if len(c.packages) == 0 {
		return true
	}

	for _, pattern := range c.packages {
		if pattern == "./..." {
			return true
		}

		// Convert relative pattern to absolute module path.
		// "./internal/services/..." → "module/internal/services"
		rel := strings.TrimPrefix(pattern, "./")
		recursive := strings.HasSuffix(rel, "/...")
		rel = strings.TrimSuffix(rel, "/...")

		abs := modulePath + "/" + rel

		if recursive {
			if pkgPath == abs || strings.HasPrefix(pkgPath, abs+"/") {
				return true
			}
		} else {
			if pkgPath == abs {
				return true
			}
		}
	}

	return false
}

// IsExcludedType reports whether a fully qualified type name (e.g. "github.com/x/pkg.MyInterface")
// is in the exclude list.
func (c *Config) IsExcludedType(fullName string) bool {
	_, ok := c.excludeTypes[fullName]
	return ok
}

// IsExcludedPackage reports whether the given package path is in the exclude list.
func (c *Config) IsExcludedPackage(pkgPath string) bool {
	_, ok := c.excludePackages[pkgPath]
	return ok
}
