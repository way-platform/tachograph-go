# Way Platform Magefile Conventions

**Table of Contents:**

1. [Directory Structure](#1-directory-structure)
2. [Running Mage](#2-running-mage)
3. [Dependency Management](#3-dependency-management-toolsgomod)
4. [Standard Magefile Structure](#4-standard-magefile-structure)
5. [Documentation](#5-documentation)
6. [Linter Configuration](#6-linter-configuration)

Way Platform follows specific patterns for Magefiles to ensure consistency and reliable dependency management.

## 1. Directory Structure

Place build tools and the Magefile in a `tools/` subdirectory:

```text
my-project/
├── tools/
│   ├── go.mod        # Dependencies for build tools (including mage)
│   ├── mage          # Entrypoint script (optional)
│   └── magefile.go   # The actual Magefile
├── go.mod            # Project dependencies
└── ...
```

## 2. Running Mage

Use Go 1.24+ `tool` support to run Mage with pinned versions and correct module context.

**Native Usage (Recommended):**

```sh
cd tools
go tool mage -l       # List targets
go tool mage Build    # Run Build target
```

**Wrapper Script (Optional):**

Create `tools/mage` to run from project root:

```sh
#!/bin/sh
root="$(dirname "$(realpath "$0")")/.."
go tool -modfile "${root}/tools/go.mod" mage -d "${root}/tools" "$@"
```

## 3. Dependency Management (`tools/go.mod`)

Use Go 1.24+ `tool` directive:

```go
module github.com/way-platform/my-project/tools

go 1.24

tool (
    github.com/magefile/mage
    github.com/golangci/golangci-lint/v2/cmd/golangci-lint
    // ... other tools (e.g. github.com/bufbuild/buf/cmd/buf)
)

require github.com/magefile/mage v1.15.0 // ...
```

## 4. Standard Magefile Structure

**Path**: `tools/magefile.go`

### Header and Imports

```go
//go:build mage

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/mg"
)

var Default = Build
```

### Standard Helpers

```go
// root returns the absolute path to the project root.
func root(subdirs ...string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to get current file path")
	}
	rootDir := filepath.Dir(filepath.Dir(filename))
	return filepath.Join(append([]string{rootDir}, subdirs...)...)
}

// cmd runs a command in a specific directory.
func cmd(dir string, command string, args ...string) *exec.Cmd {
	return cmdWith(nil, dir, command, args...)
}

// cmdWith runs a command with environment variables.
func cmdWith(env map[string]string, dir string, command string, args ...string) *exec.Cmd {
	c := exec.Command(command, args...)
	c.Env = os.Environ()
	for key, value := range env {
		c.Env = append(c.Env, fmt.Sprintf("%s=%s", key, value))
	}
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c
}

// tool runs a go tool command using tools/go.mod.
func tool(dir string, toolName string, args ...string) *exec.Cmd {
	return toolWith(nil, dir, toolName, args...)
}

// toolWith runs a go tool command with environment variables.
func toolWith(env map[string]string, dir string, toolName string, args ...string) *exec.Cmd {
	cmdArgs := []string{"tool", "-modfile", filepath.Join(root(), "tools", "go.mod"), toolName}
	return cmdWith(env, dir, "go", append(cmdArgs, args...)...)
}
```

### Standard Targets

Top-level targets (invoked directly) should include a log line at the start. Targets called only via `mg.SerialDeps`/`mg.Deps` don't need log lines.

```go
// Build runs the full build pipeline.
func Build() {
	mg.SerialDeps(
		Download,
		Format,
		Lint,
		Generate,
		Tidy,
		Diff,
	)
}

// Download downloads dependencies.
func Download() error {
	log.Println("downloading dependencies")
	return cmd(root(), "go", "mod", "download").Run()
}

// Format formats code.
func Format() error {
	log.Println("formatting code")
	return tool(root(), "buf", "format", "-w").Run()
}

// Lint runs linters and fixes code style issues.
func Lint() error {
	log.Println("linting and fixing code")
	return tool(root(), "golangci-lint", "run", "--fix").Run()
}

// Generate runs code generation.
func Generate() error {
	log.Println("generating code")
	return cmd(root(), "go", "generate", "./...").Run()
}

// Tidy runs go mod tidy.
func Tidy() error {
	log.Println("tidying Go mod files")
	return cmd(root(), "go", "mod", "tidy").Run()
}

// Diff checks if there are uncommitted changes (useful for CI).
func Diff() error {
	log.Println("checking for git diffs")
	return cmd(root(), "git", "diff", "--exit-code").Run()
}
```

## 5. Documentation

All exported targets MUST have a documentation comment:

- First sentence is the short description (shown in `mage -l`)
- Subsequent lines provide detailed help

```go
// Build runs the full CI build.
func Build() {
    // ...
}
```

## 6. Linter Configuration

Use `golangci-lint` v2 with standard formatters enabled to maintain consistent code style.

**Path**: `.golangci.yml`

```yaml
version: "2"

formatters:
  enable:
    - gci
    - gofumpt
    - goimports
    - golines
```
