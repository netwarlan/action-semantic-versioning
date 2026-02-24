package action

import (
	"fmt"
	"os"
	"strings"
)

// Inputs holds the parsed GitHub Action inputs.
type Inputs struct {
	Token              string
	DefaultVersion     string
	TagPrefix          string
	CreateRelease      bool
	ReleaseDraft       bool
	ReleasePrerelease  bool
	BumpPatchOnUnknown bool
	DryRun             bool
}

// ParseInputs reads action inputs from INPUT_* environment variables.
func ParseInputs() (Inputs, error) {
	token := getInput("TOKEN")
	if token == "" {
		return Inputs{}, fmt.Errorf("input 'token' is required")
	}

	return Inputs{
		Token:              token,
		DefaultVersion:     getInputDefault("DEFAULT-VERSION", "v0.1.0"),
		TagPrefix:          getInputDefault("TAG-PREFIX", "v"),
		CreateRelease:      parseBool(getInput("CREATE-RELEASE")),
		ReleaseDraft:       parseBool(getInput("RELEASE-DRAFT")),
		ReleasePrerelease:  parseBool(getInput("RELEASE-PRERELEASE")),
		BumpPatchOnUnknown: parseBool(getInput("BUMP-PATCH-ON-UNKNOWN")),
		DryRun:             parseBool(getInput("DRY-RUN")),
	}, nil
}

func getInput(name string) string {
	return strings.TrimSpace(os.Getenv("INPUT_" + name))
}

func getInputDefault(name, defaultVal string) string {
	v := getInput(name)
	if v == "" {
		return defaultVal
	}
	return v
}

func parseBool(s string) bool {
	return strings.EqualFold(s, "true")
}
