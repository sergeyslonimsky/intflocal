package intflocal

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/sergeyslonimsky/intflocal/analyzer"
)

func init() {
	register.Plugin("intflocal", newPlugin)
}

func newPlugin(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[analyzer.Settings](settings)
	if err != nil {
		return nil, err
	}

	return &plugin{settings: s}, nil
}

type plugin struct {
	settings analyzer.Settings
}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.New(p.settings)}, nil
}

func (p *plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
