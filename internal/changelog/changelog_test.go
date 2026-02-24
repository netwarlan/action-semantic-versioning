package changelog

import (
	"strings"
	"testing"

	"github.com/netwarlan/action-semantic-versioning/internal/commit"
)

func TestGenerate(t *testing.T) {
	commits := []commit.ConventionalCommit{
		{Type: "feat", Scope: "api", Description: "add user endpoint", Hash: "abc1234567", Breaking: true},
		{Type: "feat", Scope: "ui", Description: "add dark mode", Hash: "def5678901"},
		{Type: "feat", Description: "add preferences page", Hash: "ghi9012345"},
		{Type: "fix", Description: "correct timezone handling", Hash: "jkl3456789"},
		{Type: "perf", Description: "optimize query", Hash: "prf1234567"},
		{Type: "docs", Description: "update API docs", Hash: "mno7890123"},
	}

	result := Generate(commits, "v1.2.3", "v2.0.0")

	// Check structure
	if !strings.Contains(result, "## What's Changed") {
		t.Error("missing header")
	}
	if !strings.Contains(result, "### Breaking Changes") {
		t.Error("missing Breaking Changes section")
	}
	if !strings.Contains(result, "### Features") {
		t.Error("missing Features section")
	}
	if !strings.Contains(result, "### Bug Fixes") {
		t.Error("missing Bug Fixes section")
	}
	if !strings.Contains(result, "### Performance") {
		t.Error("missing Performance section")
	}
	if !strings.Contains(result, "### Other Changes") {
		t.Error("missing Other Changes section")
	}

	// Check content
	if !strings.Contains(result, "**api**: add user endpoint (abc1234)") {
		t.Error("missing scoped breaking change")
	}
	if !strings.Contains(result, "**ui**: add dark mode (def5678)") {
		t.Error("missing scoped feature")
	}
	if !strings.Contains(result, "- add preferences page (ghi9012)") {
		t.Error("missing unscoped feature")
	}
	if !strings.Contains(result, "**Full Changelog**: v1.2.3...v2.0.0") {
		t.Error("missing full changelog link")
	}
}

func TestGenerateEmpty(t *testing.T) {
	result := Generate(nil, "v1.0.0", "v1.0.1")
	if !strings.Contains(result, "## What's Changed") {
		t.Error("missing header")
	}
	// Should still have the changelog link
	if !strings.Contains(result, "**Full Changelog**: v1.0.0...v1.0.1") {
		t.Error("missing full changelog link")
	}
}

func TestGenerateNonConventional(t *testing.T) {
	commits := []commit.ConventionalCommit{
		{Type: "", Raw: "Update README\nsome body", Hash: "abc1234567"},
	}

	result := Generate(commits, "", "")
	if !strings.Contains(result, "- Update README (abc1234)") {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestGenerateOmitsEmptySections(t *testing.T) {
	commits := []commit.ConventionalCommit{
		{Type: "fix", Description: "a bug", Hash: "abc1234567"},
	}

	result := Generate(commits, "v1.0.0", "v1.0.1")
	if strings.Contains(result, "### Features") {
		t.Error("Features section should be omitted when empty")
	}
	if strings.Contains(result, "### Breaking Changes") {
		t.Error("Breaking Changes section should be omitted when empty")
	}
	if !strings.Contains(result, "### Bug Fixes") {
		t.Error("Bug Fixes section should be present")
	}
}
