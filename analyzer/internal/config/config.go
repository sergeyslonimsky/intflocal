package config

import "strings"

// Config holds parsed and validated linter configuration.
type Config struct {
	packages        []string
	skipPackages    []string
	excludeTypes    map[string]struct{}
	excludePackages map[string]struct{}
}

// Settings is the raw settings from golangci-lint configuration.
type Settings struct {
	Packages        []string `json:"packages"        mapstructure:"packages"`
	SkipPackages    []string `json:"skipPackages"    mapstructure:"skipPackages"`
	ExcludeTypes    []string `json:"excludeTypes"    mapstructure:"excludeTypes"`
	ExcludePackages []string `json:"excludePackages" mapstructure:"excludePackages"`
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
		skipPackages:    s.SkipPackages,
		excludeTypes:    excludeTypes,
		excludePackages: excludePackages,
	}
}

// ShouldCheckPackage reports whether the package with the given path should be analyzed.
// modulePath is the module path from go.mod (used to resolve relative package patterns).
// If no packages are configured, all packages are checked (minus any skipPackages).
func (c *Config) ShouldCheckPackage(pkgPath, modulePath string) bool {
	// Blacklist: skip explicitly excluded packages regardless of whitelist.
	for _, pattern := range c.skipPackages {
		if matchPattern(pkgPath, modulePath, pattern) {
			return false
		}
	}

	if len(c.packages) == 0 {
		return true
	}

	for _, pattern := range c.packages {
		if matchPattern(pkgPath, modulePath, pattern) {
			return true
		}
	}

	return false
}

// matchPattern reports whether pkgPath matches the given relative package pattern.
func matchPattern(pkgPath, modulePath, pattern string) bool {
	if pattern == "./..." {
		return true
	}

	rel := strings.TrimPrefix(pattern, "./")
	recursive := strings.HasSuffix(rel, "/...")
	rel = strings.TrimSuffix(rel, "/...")
	abs := modulePath + "/" + rel

	if recursive {
		return pkgPath == abs || strings.HasPrefix(pkgPath, abs+"/")
	}

	return pkgPath == abs
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
