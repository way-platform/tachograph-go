//go:build mage

package main

import (
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	return cmd(root("cmd/tachograph"), "go", "install", ".").Run()
}

// VHS records the CLI GIF using VHS.
func VHS() error {
	mg.Deps(CLI)
	return tool(root("docs"), "vhs", "cli.tape").Run()
}

// RegulationPDF fetches the EU regulation document.
func RegulationPDF() error {
	targetFile := root("docs", "regulation", "regulation.pdf")
	const url = "https://eur-lex.europa.eu/legal-content/EN/TXT/PDF/?uri=CELEX:02016R0799-20230821"
	if _, err := os.Stat(targetFile); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(targetFile), 0o600); err != nil {
		return err
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

// RegulationChapters splits the regulation PDF into chapters.
func RegulationChapters() error {
	if err := os.RemoveAll(root("docs", "regulation", "chapters")); err != nil {
		return err
	}
	if err := os.MkdirAll(root("docs", "regulation", "chapters"), 0o700); err != nil {
		return err
	}
	if err := tool(
		root("docs/regulation"),
		"pdfcpu",
		"split",
		"-m", "page",
		"regulation.pdf",
		"chapters",
		"7",
		"111",
		"230",
		"278",
		"303",
		"307",
		"321",
		"322",
		"326",
		"341",
		"350",
		"356",
		"391",
		"423",
		"424",
		"502",
		"524",
		"534",
		"587",
		"593",
		"598",
		"605",
		"612",
	).Run(); err != nil {
		return err
	}
	for _, op := range []struct {
		from string
		to   string
	}{
		{from: "regulation_1-6.pdf", to: "01-articles.pdf"},
		{from: "regulation_7-110.pdf", to: "02-requirements.pdf"},
		{from: "regulation_111-229.pdf", to: "03-data-dictionary.pdf"},
		{from: "regulation_230-277.pdf", to: "04-tachograph-cards-specification.pdf"},
		{from: "regulation_278-302.pdf", to: "05-tachograph-cards-file-structure.pdf"},
		{from: "regulation_303-306.pdf", to: "06-pictograms.pdf"},
		{from: "regulation_307-320.pdf", to: "07-printouts.pdf"},
		{from: "regulation_321.pdf", to: "08-display.pdf"},
		{from: "regulation_322-325.pdf", to: "09-front-connector.pdf"},
		{from: "regulation_326-340.pdf", to: "10-data-downloading-protocols.pdf"},
		{from: "regulation_341-349.pdf", to: "11-response-message-content.pdf"},
		{from: "regulation_350-355.pdf", to: "12-card-downloading.pdf"},
		{from: "regulation_356-390.pdf", to: "13-calibration-protocol.pdf"},
		{from: "regulation_391-422.pdf", to: "14-type-approval.pdf"},
		{from: "regulation_423.pdf", to: "15-security-requirements.pdf"},
		{from: "regulation_424-501.pdf", to: "16-common-security-mechanisms.pdf"},
		{from: "regulation_502-523.pdf", to: "17-gnss-positioning.pdf"},
		{from: "regulation_524-533.pdf", to: "18-its-interface.pdf"},
		{from: "regulation_534-586.pdf", to: "19-remote-communication.pdf"},
		{from: "regulation_587-592.pdf", to: "20-computation-of-driving-time.pdf"},
		{from: "regulation_593-597.pdf", to: "21-migration.pdf"},
		{from: "regulation_598-604.pdf", to: "22-adaptor.pdf"},
		{from: "regulation_605-611.pdf", to: "23-osnma-galileo.pdf"},
		{from: "regulation_612-616.pdf", to: "24-approval-mark-and-certificate.pdf"},
	} {
		slog.Info("renaming", "from", op.from, "to", op.to)
		if err := os.Rename(
			root("docs", "regulation", "chapters", op.from),
			root("docs", "regulation", "chapters", op.to),
		); err != nil {
			return err
		}
	}
	return nil
}

// RegulationOCR runs OCR on the regulation PDF chapter files.
func RegulationOCR() error {
	for _, target := range []string{
		"02-requirements.pdf",
		"03-data-dictionary.pdf",
		"04-tachograph-cards-specification.pdf",
		"05-tachograph-cards-file-structure.pdf",
		"10-data-downloading-protocols.pdf",
		"11-response-message-content.pdf",
		"12-card-downloading.pdf",
		"15-security-requirements.pdf",
		"16-common-security-mechanisms.pdf",
		"17-gnss-positioning.pdf",
		"18-its-interface.pdf",
		"20-computation-of-driving-time.pdf",
	} {
		targetDir := root("docs", "regulation", "chapters", strings.TrimSuffix(target, ".pdf"))
		if _, err := os.Stat(targetDir); err == nil {
			slog.Info("skipping OCR", "target", target, "reason", "target dir exists")
			continue
		}
		slog.Info("running OCR", "target", target)
		if err := cmd(
			root("docs", "regulation", "chapters"),
			"uvx",
			"--from", "marker-pdf",
			"marker_single",
			target,
			"--output_dir", ".",
			"--output_format=markdown",
			"--page_separator=''",
			"--use_llm",
			"--llm_service=marker.services.vertex.GoogleVertexService",
			"--vertex_project_id=way-local-dev",
		).Run(); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
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
