# intflocal

[![CI](https://github.com/sergeyslonimsky/intflocal/actions/workflows/ci.yml/badge.svg)](https://github.com/sergeyslonimsky/intflocal/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sergeyslonimsky/intflocal)](https://goreportcard.com/report/github.com/sergeyslonimsky/intflocal)
[![codecov](https://codecov.io/gh/sergeyslonimsky/intflocal/branch/master/graph/badge.svg)](https://codecov.io/gh/sergeyslonimsky/intflocal)
[![Go Reference](https://pkg.go.dev/badge/github.com/sergeyslonimsky/intflocal.svg)](https://pkg.go.dev/github.com/sergeyslonimsky/intflocal)
[![GitHub release](https://img.shields.io/github/v/release/sergeyslonimsky/intflocal)](https://github.com/sergeyslonimsky/intflocal/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go linter that checks struct fields use locally defined private interfaces instead of importing them from external packages.

Follows the Go best practice of [defining interfaces at the point of use](https://go.dev/wiki/CodeReviewComments#interfaces) and the Interface Segregation Principle.

## Motivation

```go
// Bad: importing interface from another package
type Service struct {
    repo repository.Repository
}

// Good: defining interface locally
type repository interface {
    Get(ctx context.Context, id string) (*Entity, error)
}

type Service struct {
    repo repository
}
```

Defining interfaces locally:
- Decouples packages from each other
- Allows consumers to depend only on the methods they actually use
- Makes testing easier with focused mocks
- Follows the principle of accepting interfaces, returning structs

## Requirements

Go 1.22 or later.

## Installation

```bash
go install github.com/sergeyslonimsky/intflocal/cmd/intflocal@latest
```

## Usage

### Standalone

```bash
# Check all packages
intflocal ./...

# Check specific packages
intflocal ./internal/services/...

# With exclusions
intflocal --excludePackages=github.com/some/pkg,github.com/other/pkg ./...
intflocal --excludeTypes=github.com/some/pkg.SpecialInterface ./...

# Limit scope to specific packages (whitelist)
intflocal --packages=./internal/services/...,./pkg/handlers/... ./...

# Skip specific packages (blacklist)
intflocal --skipPackages=./internal/di/...,./test/... ./...
```

### Flags

| Flag | Description | Example |
|------|-------------|---------|
| `-excludePackages` | Comma-separated list of packages whose interfaces are always allowed | `github.com/some/pkg,github.com/other/pkg` |
| `-excludeTypes` | Comma-separated list of fully qualified type names to allow | `github.com/some/pkg.MyInterface` |
| `-packages` | Whitelist: only check these package patterns | `./internal/services/...,./pkg/handlers/...` |
| `-skipPackages` | Blacklist: skip these package patterns from analysis | `./internal/di/...,./test/...` |

### With golangci-lint (module plugin)

1. Create `.custom-gcl.yml` in your project:

```yaml
version: v2.10.1
plugins:
  - module: 'github.com/sergeyslonimsky/intflocal'
    version: v0.1.0
```

2. Add to `.golangci.yml`:

```yaml
linters-settings:
  custom:
    intflocal:
      type: "module"
      settings:
        packages:
          - "./internal/services/..."
        skipPackages:
          - "./internal/di/..."
        excludePackages:
          - "github.com/some/pkg"
        excludeTypes:
          - "github.com/some/pkg.SpecialInterface"

linters:
  disable-all: true
  enable:
    - intflocal
```

3. Build and run:

```bash
golangci-lint custom
./custom-gcl run ./...
```

Suppress individual findings with `//nolint:intflocal`.

## Local Development

To test the plugin locally with golangci-lint before publishing a new version, create a `.custom-gcl.yml` file in the project root (it is git-ignored):

```yaml
version: v2.10.1
plugins:
  - module: 'github.com/sergeyslonimsky/intflocal'
    path: /absolute/path/to/your/local/intflocal
```

Then build a custom golangci-lint binary and run it:

```bash
golangci-lint custom -c .custom-gcl.yml
./custom-gcl run ./...
```

The `path` field tells golangci-lint to load the plugin from a local directory instead of fetching a published version.

## What is checked

- Struct fields whose type is an interface imported from another package

## What is allowed

- Locally defined private interfaces
- Concrete types (structs, primitives, etc.)
- Standard library interfaces (`io.Reader`, `context.Context`, etc.)
- Builtin interfaces (`error`)
- Types and packages in the exclude lists

## License

MIT
