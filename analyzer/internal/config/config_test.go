package config_test

import (
	"testing"

	"github.com/sergeyslonimsky/intflocal/analyzer/internal/config"
)

func TestShouldCheckPackage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		packages   []string
		pkgPath    string
		modulePath string
		want       bool
	}{
		{
			name:       "no filter - check everything",
			packages:   nil,
			pkgPath:    "github.com/user/repo/internal/service",
			modulePath: "github.com/user/repo",
			want:       true,
		},
		{
			name:       "wildcard matches all",
			packages:   []string{"./..."},
			pkgPath:    "github.com/user/repo/anything",
			modulePath: "github.com/user/repo",
			want:       true,
		},
		{
			name:       "recursive pattern matches direct child",
			packages:   []string{"./internal/..."},
			pkgPath:    "github.com/user/repo/internal/service",
			modulePath: "github.com/user/repo",
			want:       true,
		},
		{
			name:       "recursive pattern matches nested package",
			packages:   []string{"./internal/..."},
			pkgPath:    "github.com/user/repo/internal/service/sub",
			modulePath: "github.com/user/repo",
			want:       true,
		},
		{
			name:       "recursive pattern does not match sibling",
			packages:   []string{"./internal/..."},
			pkgPath:    "github.com/user/repo/pkg/service",
			modulePath: "github.com/user/repo",
			want:       false,
		},
		{
			name:       "exact pattern matches",
			packages:   []string{"./internal/service"},
			pkgPath:    "github.com/user/repo/internal/service",
			modulePath: "github.com/user/repo",
			want:       true,
		},
		{
			name:       "exact pattern does not match subpackage",
			packages:   []string{"./internal/service"},
			pkgPath:    "github.com/user/repo/internal/service/sub",
			modulePath: "github.com/user/repo",
			want:       false,
		},
		{
			name:       "v2 module - recursive pattern matches",
			packages:   []string{"./internal/..."},
			pkgPath:    "github.com/user/repo/v2/internal/service",
			modulePath: "github.com/user/repo/v2",
			want:       true,
		},
		{
			name:       "v2 module - does not match wrong module",
			packages:   []string{"./internal/..."},
			pkgPath:    "github.com/user/repo/internal/service",
			modulePath: "github.com/user/repo/v2",
			want:       false,
		},
		{
			name:       "gopkg.in module - recursive pattern matches",
			packages:   []string{"./service/..."},
			pkgPath:    "gopkg.in/user/repo.v3/service/handler",
			modulePath: "gopkg.in/user/repo.v3",
			want:       true,
		},
		{
			name:       "multiple patterns - first matches",
			packages:   []string{"./internal/...", "./pkg/..."},
			pkgPath:    "github.com/user/repo/internal/service",
			modulePath: "github.com/user/repo",
			want:       true,
		},
		{
			name:       "multiple patterns - second matches",
			packages:   []string{"./internal/...", "./pkg/..."},
			pkgPath:    "github.com/user/repo/pkg/handler",
			modulePath: "github.com/user/repo",
			want:       true,
		},
		{
			name:       "multiple patterns - none match",
			packages:   []string{"./internal/...", "./pkg/..."},
			pkgPath:    "github.com/user/repo/domain/service",
			modulePath: "github.com/user/repo",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := config.New(config.Settings{Packages: tt.packages})
			got := cfg.ShouldCheckPackage(tt.pkgPath, tt.modulePath)
			if got != tt.want {
				t.Errorf("ShouldCheckPackage(%q, %q) = %v, want %v", tt.pkgPath, tt.modulePath, got, tt.want)
			}
		})
	}
}
