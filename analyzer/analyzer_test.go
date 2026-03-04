package analyzer_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/sergeyslonimsky/intflocal/analyzer"
	"github.com/sergeyslonimsky/intflocal/analyzer/internal/config"
)

func testdataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "testdata")
}

// TestAnalyzer verifies all bad patterns are reported and good patterns are not.
func TestAnalyzer(t *testing.T) {
	t.Parallel()
	a := analyzer.New(config.Settings{})
	analysistest.Run(t, testdataDir(), a, "example")
}

// TestAnalyzer_StdlibAllowed verifies that stdlib interfaces are never reported.
func TestAnalyzer_StdlibAllowed(t *testing.T) {
	t.Parallel()
	a := analyzer.New(config.Settings{})
	analysistest.Run(t, testdataDir(), a, "stdlib")
}

// TestAnalyzer_ExcludePackage verifies that interfaces from an excluded package are not reported.
func TestAnalyzer_ExcludePackage(t *testing.T) {
	t.Parallel()
	a := analyzer.New(config.Settings{
		ExcludePackages: []string{"example.com/extiface"},
	})
	analysistest.Run(t, testdataDir(), a, "excluded")
}

// TestAnalyzer_ExcludeType verifies that individually excluded types are not reported.
func TestAnalyzer_ExcludeType(t *testing.T) {
	t.Parallel()
	a := analyzer.New(config.Settings{
		ExcludeTypes: []string{
			"example.com/extiface.MyInterface",
			"example.com/extiface.AnotherInterface",
		},
	})
	analysistest.Run(t, testdataDir(), a, "excluded")
}
