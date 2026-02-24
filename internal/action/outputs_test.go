package action

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSetOutputSimple(t *testing.T) {
	dir := t.TempDir()
	outputFile := filepath.Join(dir, "output")
	t.Setenv("GITHUB_OUTPUT", outputFile)

	if err := SetOutput("version", "1.2.3"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(data); got != "version=1.2.3\n" {
		t.Errorf("output = %q", got)
	}
}

func TestSetOutputMultiline(t *testing.T) {
	dir := t.TempDir()
	outputFile := filepath.Join(dir, "output")
	t.Setenv("GITHUB_OUTPUT", outputFile)

	if err := SetOutput("changelog", "line1\nline2"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)

	if !strings.HasPrefix(content, "changelog<<EOF_") {
		t.Errorf("expected delimiter syntax, got: %q", content)
	}
	if !strings.Contains(content, "line1\nline2\n") {
		t.Errorf("expected multiline content, got: %q", content)
	}
}

func TestSetOutputMultiple(t *testing.T) {
	dir := t.TempDir()
	outputFile := filepath.Join(dir, "output")
	t.Setenv("GITHUB_OUTPUT", outputFile)

	if err := SetOutput("a", "1"); err != nil {
		t.Fatal(err)
	}
	if err := SetOutput("b", "2"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, "a=1\n") || !strings.Contains(content, "b=2\n") {
		t.Errorf("output = %q", content)
	}
}
