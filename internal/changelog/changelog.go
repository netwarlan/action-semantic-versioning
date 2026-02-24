package changelog

import (
	"fmt"
	"strings"

	"github.com/netwarlan/action-semantic-versioning/internal/commit"
)

// Generate creates a markdown changelog from parsed commits.
func Generate(commits []commit.ConventionalCommit, previousTag, newTag string) string {
	var breaking, features, fixes, perf, other []string

	for _, c := range commits {
		line := formatCommit(c)

		if c.Breaking {
			breaking = append(breaking, line)
			continue
		}

		switch c.Type {
		case "feat":
			features = append(features, line)
		case "fix":
			fixes = append(fixes, line)
		case "perf":
			perf = append(perf, line)
		case "":
			// Non-conventional commits go to other.
			other = append(other, fmt.Sprintf("- %s (%s)", firstLine(c.Raw), shortHash(c.Hash)))
		default:
			other = append(other, line)
		}
	}

	var sb strings.Builder
	sb.WriteString("## What's Changed\n")

	writeSection(&sb, "Breaking Changes", breaking)
	writeSection(&sb, "Features", features)
	writeSection(&sb, "Bug Fixes", fixes)
	writeSection(&sb, "Performance", perf)
	writeSection(&sb, "Other Changes", other)

	if previousTag != "" && newTag != "" {
		fmt.Fprintf(&sb, "\n**Full Changelog**: %s...%s\n", previousTag, newTag)
	}

	return sb.String()
}

func formatCommit(c commit.ConventionalCommit) string {
	hash := shortHash(c.Hash)
	if c.Scope != "" {
		return fmt.Sprintf("- **%s**: %s (%s)", c.Scope, c.Description, hash)
	}
	return fmt.Sprintf("- %s (%s)", c.Description, hash)
}

func shortHash(hash string) string {
	if len(hash) > 7 {
		return hash[:7]
	}
	return hash
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

func writeSection(sb *strings.Builder, title string, items []string) {
	if len(items) == 0 {
		return
	}
	fmt.Fprintf(sb, "\n### %s\n", title)
	for _, item := range items {
		sb.WriteString(item)
		sb.WriteByte('\n')
	}
}
