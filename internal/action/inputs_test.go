package action

import (
	"testing"
)

func TestParseInputs(t *testing.T) {
	t.Setenv("INPUT_TOKEN", "ghp_test123")
	t.Setenv("INPUT_DEFAULT-VERSION", "v1.0.0")
	t.Setenv("INPUT_TAG-PREFIX", "v")
	t.Setenv("INPUT_CREATE-RELEASE", "true")
	t.Setenv("INPUT_RELEASE-DRAFT", "false")
	t.Setenv("INPUT_RELEASE-PRERELEASE", "True")
	t.Setenv("INPUT_BUMP-PATCH-ON-UNKNOWN", "TRUE")
	t.Setenv("INPUT_DRY-RUN", "false")

	inputs, err := ParseInputs()
	if err != nil {
		t.Fatal(err)
	}

	if inputs.Token != "ghp_test123" {
		t.Errorf("Token = %q", inputs.Token)
	}
	if inputs.DefaultVersion != "v1.0.0" {
		t.Errorf("DefaultVersion = %q", inputs.DefaultVersion)
	}
	if !inputs.CreateRelease {
		t.Error("CreateRelease should be true")
	}
	if inputs.ReleaseDraft {
		t.Error("ReleaseDraft should be false")
	}
	if !inputs.ReleasePrerelease {
		t.Error("ReleasePrerelease should be true (case-insensitive)")
	}
	if !inputs.BumpPatchOnUnknown {
		t.Error("BumpPatchOnUnknown should be true (case-insensitive)")
	}
	if inputs.DryRun {
		t.Error("DryRun should be false")
	}
}

func TestParseInputsDefaults(t *testing.T) {
	t.Setenv("INPUT_TOKEN", "test")

	inputs, err := ParseInputs()
	if err != nil {
		t.Fatal(err)
	}
	if inputs.DefaultVersion != "v0.1.0" {
		t.Errorf("DefaultVersion = %q, want v0.1.0", inputs.DefaultVersion)
	}
	if inputs.TagPrefix != "v" {
		t.Errorf("TagPrefix = %q, want v", inputs.TagPrefix)
	}
}

func TestParseInputsMissingToken(t *testing.T) {
	// t.Setenv not called for INPUT_TOKEN, so it's unset.
	_, err := ParseInputs()
	if err == nil {
		t.Error("expected error for missing token")
	}
}
