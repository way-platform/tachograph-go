//go:build mage

package main

import (
	"io"
	"io/fs"
	"net/http"
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

// RegulationOriginal fetches the EU regulation document.
func RegulationOriginal() error {
	const url = "https://publications.europa.eu/resource/cellar/50ef99c6-7896-11ec-9136-01aa75ed71a1.0006.02/DOC_2"
	outputPath := root("docs", "regulation.html")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o600); err != nil {
		return err
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

// Regulation runs the regulation-parser on the original document to produce split files.
func Regulation() error {
	return cmd(
		root("tools"),
		"go", "run", "./cmd/regulation-parser",
		"-i", root("docs/regulation.html"),
		"-d", root("docs/regulation"),
	).Run()
}

// RegulationDataDictionary runs the regulation-preprocessor on the data dictionary file with level 2 chunking.
func RegulationDataDictionary() error {
	return cmd(
		root("tools"),
		"go", "run", "./cmd/regulation-preprocessor",
		"-i", root("docs/regulation/09-appendix-1-data-dictionary.html"),
		"-d", root("docs/regulation/data-dictionary"),
		"-chunk-level", "2",
	).Run()
}

// RegulationAnnex1C runs the regulation-preprocessor on the annex 1c file with level 2 chunking.
func RegulationAnnex1C() error {
	return cmd(
		root("tools"),
		"go", "run", "./cmd/regulation-preprocessor",
		"-i", root("docs/regulation/08-annex-1c-requirements-for-construction-testing-installation-and-inspection.html"),
		"-d", root("docs/regulation/annex-1c"),
		"-chunk-level", "2",
	).Run()
}

// RegulationTachographCards runs the regulation-preprocessor on the tachograph cards file with level 2 chunking.
func RegulationTachographCards() error {
	return cmd(
		root("tools"),
		"go", "run", "./cmd/regulation-preprocessor",
		"-i", root("docs/regulation/10-appendix-2-tachograph-cards-specification.html"),
		"-d", root("docs/regulation/tachograph-cards-specification"),
		"-chunk-level", "2",
	).Run()
}

// RegulationDataTypeDefinitions runs the regulation-preprocessor on the data type definitions file with level 3 chunking.
func RegulationDataTypeDefinitions() error {
	return cmd(
		root("tools"),
		"go", "run", "./cmd/regulation-preprocessor",
		"-i", root("docs/regulation/data-dictionary/04-data-type-definitions.html"),
		"-d", root("docs/regulation/data-dictionary/data-type-definitions"),
		"-chunk-level", "3",
	).Run()
}

// RegulationChunked runs all regulation chunking tasks for the large files.
func RegulationChunked() error {
	mg.SerialDeps(
		RegulationDataDictionary,
		RegulationDataTypeDefinitions,
		RegulationAnnex1C,
		RegulationTachographCards,
	)
	return nil
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
