package commit

import (
	"regexp"
	"strings"
)

var (
	subjectRegex = regexp.MustCompile(`^(\w+)(\(([^)]*)\))?(!)?:\s*(.+)$`)
	footerRegex  = regexp.MustCompile(`^([\w-]+|BREAKING CHANGE)\s*:\s*(.+)$`)
)

// Parse parses a commit message into a ConventionalCommit.
// Non-conventional messages return a ConventionalCommit with an empty Type.
func Parse(hash, message string) ConventionalCommit {
	cc := ConventionalCommit{
		Hash: hash,
		Raw:  message,
	}

	lines := strings.Split(message, "\n")
	if len(lines) == 0 {
		return cc
	}

	matches := subjectRegex.FindStringSubmatch(strings.TrimSpace(lines[0]))
	if matches == nil {
		return cc
	}

	cc.Type = strings.ToLower(matches[1])
	cc.Scope = matches[3]
	if matches[4] == "!" {
		cc.Breaking = true
	}
	cc.Description = matches[5]

	// Parse body and footers from remaining lines.
	if len(lines) > 1 {
		cc.Body, cc.Footers = parseBodyAndFooters(lines[1:])
	}

	// Check footers for BREAKING CHANGE.
	for _, f := range cc.Footers {
		token := strings.ToUpper(f.Token)
		if token == "BREAKING CHANGE" || token == "BREAKING-CHANGE" {
			cc.Breaking = true
		}
	}

	return cc
}

func parseBodyAndFooters(lines []string) (string, []Footer) {
	// Find the start of the footer section by scanning from the end.
	// Footers are contiguous lines at the end that match the footer pattern.
	footerStart := len(lines)
	for i := len(lines) - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" {
			// Blank line before footers — this is the separator.
			break
		}
		if footerRegex.MatchString(trimmed) {
			footerStart = i
		} else {
			break
		}
	}

	var footers []Footer
	for _, line := range lines[footerStart:] {
		trimmed := strings.TrimSpace(line)
		if m := footerRegex.FindStringSubmatch(trimmed); m != nil {
			footers = append(footers, Footer{Token: m[1], Value: m[2]})
		}
	}

	// Body is everything between the blank line after the subject and the footer section.
	bodyLines := lines[:footerStart]
	body := strings.TrimSpace(strings.Join(bodyLines, "\n"))

	return body, footers
}

// DetermineBump determines the highest bump type from a list of commits.
func DetermineBump(commits []ConventionalCommit, bumpPatchOnUnknown bool) BumpType {
	bump := BumpNone

	for _, c := range commits {
		if c.Breaking {
			return BumpMajor
		}

		switch c.Type {
		case "feat":
			if bump < BumpMinor {
				bump = BumpMinor
			}
		case "fix", "perf":
			if bump < BumpPatch {
				bump = BumpPatch
			}
		default:
			if bumpPatchOnUnknown && c.Type != "" && bump < BumpPatch {
				bump = BumpPatch
			}
			// Non-conventional commits (empty Type) only bump if bumpPatchOnUnknown.
			if bumpPatchOnUnknown && c.Type == "" && bump < BumpPatch {
				bump = BumpPatch
			}
		}
	}

	return bump
}
