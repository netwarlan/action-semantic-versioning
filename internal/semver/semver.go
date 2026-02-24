package semver

import (
	"fmt"
	"regexp"
	"strconv"
)

var semverRegex = regexp.MustCompile(`^(v)?(\d+)\.(\d+)\.(\d+)(?:-([a-zA-Z0-9][a-zA-Z0-9.]*))??(?:\+([a-zA-Z0-9][a-zA-Z0-9.]*))?$`)

// Version represents a parsed semantic version.
type Version struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
	Metadata   string
	Prefix     string
}

// Parse parses a version string like "v1.2.3", "1.2.3-alpha.1", or "v1.0.0+build.123".
func Parse(s string) (Version, error) {
	matches := semverRegex.FindStringSubmatch(s)
	if matches == nil {
		return Version{}, fmt.Errorf("invalid semver: %q", s)
	}

	major, _ := strconv.Atoi(matches[2])
	minor, _ := strconv.Atoi(matches[3])
	patch, _ := strconv.Atoi(matches[4])

	return Version{
		Prefix:     matches[1],
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: matches[5],
		Metadata:   matches[6],
	}, nil
}

// BumpMajor returns a new version with major incremented and minor/patch reset.
func (v Version) BumpMajor() Version {
	return Version{
		Major:  v.Major + 1,
		Minor:  0,
		Patch:  0,
		Prefix: v.Prefix,
	}
}

// BumpMinor returns a new version with minor incremented and patch reset.
func (v Version) BumpMinor() Version {
	return Version{
		Major:  v.Major,
		Minor:  v.Minor + 1,
		Patch:  0,
		Prefix: v.Prefix,
	}
}

// BumpPatch returns a new version with patch incremented.
func (v Version) BumpPatch() Version {
	return Version{
		Major:  v.Major,
		Minor:  v.Minor,
		Patch:  v.Patch + 1,
		Prefix: v.Prefix,
	}
}

// String returns the version as a string with its original prefix.
func (v Version) String() string {
	s := fmt.Sprintf("%s%d.%d.%d", v.Prefix, v.Major, v.Minor, v.Patch)
	if v.Prerelease != "" {
		s += "-" + v.Prerelease
	}
	if v.Metadata != "" {
		s += "+" + v.Metadata
	}
	return s
}

// Compare returns -1, 0, or 1 comparing v to other per semver precedence.
// Metadata is ignored per the spec.
func (v Version) Compare(other Version) int {
	if v.Major != other.Major {
		return cmpInt(v.Major, other.Major)
	}
	if v.Minor != other.Minor {
		return cmpInt(v.Minor, other.Minor)
	}
	if v.Patch != other.Patch {
		return cmpInt(v.Patch, other.Patch)
	}
	// A version with prerelease has lower precedence than the same version without.
	if v.Prerelease == "" && other.Prerelease != "" {
		return 1
	}
	if v.Prerelease != "" && other.Prerelease == "" {
		return -1
	}
	if v.Prerelease < other.Prerelease {
		return -1
	}
	if v.Prerelease > other.Prerelease {
		return 1
	}
	return 0
}

func cmpInt(a, b int) int {
	if a < b {
		return -1
	}
	return 1
}
