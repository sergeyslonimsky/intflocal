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

func TestAnalyzer(t *testing.T) {
	a := analyzer.New(config.Settings{})
	analysistest.Run(t, testdataDir(), a, "example")
}
