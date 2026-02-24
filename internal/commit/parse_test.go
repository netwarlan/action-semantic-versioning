package commit

import (
	"testing"
)

func TestParseSimple(t *testing.T) {
	tests := []struct {
		name    string
		message string
		wantType    string
		wantScope   string
		wantDesc    string
		wantBreak   bool
	}{
		{"fix", "fix: correct typo", "fix", "", "correct typo", false},
		{"feat", "feat: add login", "feat", "", "add login", false},
		{"scoped", "feat(api): add endpoint", "feat", "api", "add endpoint", false},
		{"breaking bang", "feat!: redesign auth", "feat", "", "redesign auth", true},
		{"scoped breaking", "fix(core)!: change return type", "fix", "core", "change return type", true},
		{"docs", "docs: update readme", "docs", "", "update readme", false},
		{"chore", "chore: update deps", "chore", "", "update deps", false},
		{"non-conventional", "Update README", "", "", "", false},
		{"merge commit", "Merge pull request #42 from feature-branch", "", "", "", false},
		{"empty", "", "", "", "", false},
		{"nested scope", "fix(auth/oauth): handle token refresh", "fix", "auth/oauth", "handle token refresh", false},
		{"uppercase type", "FEAT: something", "feat", "", "something", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := Parse("abc123", tt.message)
			if cc.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", cc.Type, tt.wantType)
			}
			if cc.Scope != tt.wantScope {
				t.Errorf("Scope = %q, want %q", cc.Scope, tt.wantScope)
			}
			if cc.Description != tt.wantDesc {
				t.Errorf("Description = %q, want %q", cc.Description, tt.wantDesc)
			}
			if cc.Breaking != tt.wantBreak {
				t.Errorf("Breaking = %v, want %v", cc.Breaking, tt.wantBreak)
			}
			if cc.Hash != "abc123" {
				t.Errorf("Hash = %q, want %q", cc.Hash, "abc123")
			}
		})
	}
}

func TestParseWithBodyAndFooters(t *testing.T) {
	msg := `feat: add user system

This adds the complete user management system.

BREAKING CHANGE: removed legacy auth endpoints
Reviewed-by: Alice`

	cc := Parse("def456", msg)

	if cc.Type != "feat" {
		t.Errorf("Type = %q, want feat", cc.Type)
	}
	if !cc.Breaking {
		t.Error("expected Breaking to be true")
	}
	if cc.Body != "This adds the complete user management system." {
		t.Errorf("Body = %q", cc.Body)
	}
	if len(cc.Footers) != 2 {
		t.Fatalf("expected 2 footers, got %d", len(cc.Footers))
	}
	if cc.Footers[0].Token != "BREAKING CHANGE" || cc.Footers[0].Value != "removed legacy auth endpoints" {
		t.Errorf("footer[0] = %+v", cc.Footers[0])
	}
}

func TestParseBreakingChangeHyphenated(t *testing.T) {
	msg := `refactor: change API

BREAKING-CHANGE: new response format`

	cc := Parse("ghi789", msg)
	if !cc.Breaking {
		t.Error("expected Breaking to be true for BREAKING-CHANGE footer")
	}
}

func TestDetermineBump(t *testing.T) {
	tests := []struct {
		name               string
		commits            []ConventionalCommit
		bumpPatchOnUnknown bool
		want               BumpType
	}{
		{
			"empty",
			nil,
			false,
			BumpNone,
		},
		{
			"fix only",
			[]ConventionalCommit{{Type: "fix"}},
			false,
			BumpPatch,
		},
		{
			"feat only",
			[]ConventionalCommit{{Type: "feat"}},
			false,
			BumpMinor,
		},
		{
			"breaking",
			[]ConventionalCommit{{Type: "feat", Breaking: true}},
			false,
			BumpMajor,
		},
		{
			"mixed fix and feat",
			[]ConventionalCommit{{Type: "fix"}, {Type: "feat"}, {Type: "docs"}},
			false,
			BumpMinor,
		},
		{
			"mixed with breaking",
			[]ConventionalCommit{{Type: "fix"}, {Type: "feat", Breaking: true}},
			false,
			BumpMajor,
		},
		{
			"docs only no bump",
			[]ConventionalCommit{{Type: "docs"}},
			false,
			BumpNone,
		},
		{
			"docs with bump-patch-on-unknown",
			[]ConventionalCommit{{Type: "docs"}},
			true,
			BumpPatch,
		},
		{
			"non-conventional with bump-patch-on-unknown",
			[]ConventionalCommit{{Type: ""}},
			true,
			BumpPatch,
		},
		{
			"perf is patch",
			[]ConventionalCommit{{Type: "perf"}},
			false,
			BumpPatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetermineBump(tt.commits, tt.bumpPatchOnUnknown)
			if got != tt.want {
				t.Errorf("DetermineBump() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBumpTypeString(t *testing.T) {
	tests := []struct {
		b    BumpType
		want string
	}{
		{BumpNone, "none"},
		{BumpPatch, "patch"},
		{BumpMinor, "minor"},
		{BumpMajor, "major"},
	}
	for _, tt := range tests {
		if got := tt.b.String(); got != tt.want {
			t.Errorf("BumpType(%d).String() = %q, want %q", tt.b, got, tt.want)
		}
	}
}
