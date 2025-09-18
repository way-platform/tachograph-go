//go:build mage

package main

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = Build

// Build runs a full CI build.
func Build() {
	mg.SerialDeps(
		Download,
		Generate,
		Lint,
		Test,
		Tidy,
		CLI,
		Diff,
	)
}

// Lint runs the Go linter.
func Lint() error {
	return forEachGoMod(func(dir string) error {
		return tool(dir, "golangci-lint", "run", "--path-prefix", dir, "--build-tags", "mage").Run()
	})
}

// Test runs the Go tests.
func Test() error {
	return cmd(root(), "go", "test", "-v", "-cover", "./...").Run()
}

// Download downloads the Go dependencies.
func Download() error {
	return forEachGoMod(func(dir string) error {
		return cmd(dir, "go", "mod", "download").Run()
	})
}

// Tidy tidies the Go mod files.
func Tidy() error {
	return forEachGoMod(func(dir string) error {
		return cmd(dir, "go", "mod", "tidy", "-v").Run()
	})
}

// Diff checks for git diffs.
func Diff() error {
	return cmd(root(), "git", "diff", "--exit-code").Run()
}

// Generate runs all code generators.
func Generate() error {
	return forEachGoMod(func(dir string) error {
		return cmd(dir, "go", "generate", "-v", "./...").Run()
	})
}

// CLI builds the CLI.
func CLI() error {
	return cmd(root("cmd/tacho"), "go", "install", ".").Run()
}

// VHS records the CLI GIF using VHS.
func VHS() error {
	mg.Deps(CLI)
	return tool(root("docs"), "vhs", "cli.tape").Run()
}

func forEachGoMod(f func(dir string) error) error {
	return filepath.WalkDir(root(), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() != "go.mod" {
			return nil
		}
		return f(filepath.Dir(path))
	})
}

func cmd(dir string, command string, args ...string) *exec.Cmd {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func tool(dir string, tool string, args ...string) *exec.Cmd {
	cmdArgs := []string{"tool", "-modfile", filepath.Join(root(), "tools", "go.mod"), tool}
	return cmd(dir, "go", append(cmdArgs, args...)...)
}

func root(subdirs ...string) string {
	result, err := sh.Output("git", "rev-parse", "--show-toplevel")
	if err != nil {
		panic(err)
	}
	return filepath.Join(append([]string{result}, subdirs...)...)
}
