package main

import (
	"fmt"
	"os"

	"github.com/netwarlan/action-semantic-versioning/internal/action"
	"github.com/netwarlan/action-semantic-versioning/internal/changelog"
	"github.com/netwarlan/action-semantic-versioning/internal/commit"
	"github.com/netwarlan/action-semantic-versioning/internal/git"
	"github.com/netwarlan/action-semantic-versioning/internal/github"
	"github.com/netwarlan/action-semantic-versioning/internal/semver"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	inputs, err := action.ParseInputs()
	if err != nil {
		return err
	}

	// Validate default version is valid semver.
	if _, err := semver.Parse(inputs.DefaultVersion); err != nil {
		return fmt.Errorf("invalid default-version %q: %w", inputs.DefaultVersion, err)
	}

	gitClient := &git.Client{}

	// Check for shallow clone.
	shallow, err := gitClient.IsShallowRepository()
	if err != nil {
		return fmt.Errorf("checking repository depth: %w", err)
	}
	if shallow {
		return fmt.Errorf("shallow clone detected — use 'actions/checkout' with 'fetch-depth: 0' to fetch full history")
	}

	// Find latest semver tag.
	latestTag, err := gitClient.FindLatestSemverTag(inputs.TagPrefix)
	if err != nil {
		return fmt.Errorf("finding latest tag: %w", err)
	}

	isInitial := latestTag == ""
	previousVersion := latestTag
	if isInitial {
		previousVersion = ""
		fmt.Printf("No existing semver tags found. Will use default version: %s\n", inputs.DefaultVersion)
	} else {
		fmt.Printf("Latest version tag: %s\n", latestTag)
	}

	// List commits since last tag.
	rawCommits, err := gitClient.ListCommitsSince(latestTag)
	if err != nil {
		return fmt.Errorf("listing commits: %w", err)
	}

	if len(rawCommits) == 0 {
		fmt.Println("No new commits since last tag.")
		return writeSkippedOutputs(previousVersion)
	}

	fmt.Printf("Found %d commit(s) since last tag.\n", len(rawCommits))

	// Parse commits.
	var commits []commit.ConventionalCommit
	for _, rc := range rawCommits {
		commits = append(commits, commit.Parse(rc.Hash, rc.Message))
	}

	// Determine bump type.
	bumpType := commit.DetermineBump(commits, inputs.BumpPatchOnUnknown)

	if bumpType == commit.BumpNone {
		fmt.Println("No version-bumping commits found.")
		return writeSkippedOutputs(previousVersion)
	}

	// Calculate new version.
	var newVersion semver.Version
	if isInitial {
		// Use the default version directly for the initial release.
		newVersion, _ = semver.Parse(inputs.DefaultVersion)
	} else {
		current, _ := semver.Parse(latestTag)
		switch bumpType {
		case commit.BumpMajor:
			newVersion = current.BumpMajor()
		case commit.BumpMinor:
			newVersion = current.BumpMinor()
		case commit.BumpPatch:
			newVersion = current.BumpPatch()
		}
	}

	newTag := newVersion.String()
	changelogText := changelog.Generate(commits, previousVersion, newTag)

	fmt.Printf("Bump type: %s\n", bumpType)
	fmt.Printf("New version: %s\n", newTag)

	if !inputs.DryRun {
		// Create and push tag.
		fmt.Printf("Creating tag %s...\n", newTag)
		if err := gitClient.CreateTag(newTag); err != nil {
			return fmt.Errorf("creating tag: %w", err)
		}

		fmt.Printf("Pushing tag %s...\n", newTag)
		if err := gitClient.PushTag(newTag); err != nil {
			return fmt.Errorf("pushing tag: %w", err)
		}

		// Create release if requested.
		if inputs.CreateRelease {
			fmt.Println("Creating GitHub release...")
			releaseClient := github.NewReleaseClient(inputs.Token)
			if err := releaseClient.CreateRelease(
				newTag,
				newTag,
				changelogText,
				inputs.ReleaseDraft,
				inputs.ReleasePrerelease,
			); err != nil {
				return fmt.Errorf("creating release: %w", err)
			}
			fmt.Println("Release created successfully.")
		}
	} else {
		fmt.Println("Dry run — no tag or release created.")
	}

	// Write outputs.
	for _, o := range []struct{ k, v string }{
		{"previous-version", previousVersion},
		{"new-version", newTag},
		{"bump-type", bumpType.String()},
		{"changelog", changelogText},
		{"skipped", "false"},
	} {
		if err := action.SetOutput(o.k, o.v); err != nil {
			return fmt.Errorf("setting output %s: %w", o.k, err)
		}
	}

	return nil
}

func writeSkippedOutputs(previousVersion string) error {
	for _, o := range []struct{ k, v string }{
		{"previous-version", previousVersion},
		{"new-version", ""},
		{"bump-type", "none"},
		{"changelog", ""},
		{"skipped", "true"},
	} {
		if err := action.SetOutput(o.k, o.v); err != nil {
			return fmt.Errorf("setting output %s: %w", o.k, err)
		}
	}
	return nil
}
