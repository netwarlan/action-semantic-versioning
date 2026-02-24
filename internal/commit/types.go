package commit

// BumpType represents the kind of version bump.
type BumpType int

const (
	BumpNone  BumpType = iota
	BumpPatch
	BumpMinor
	BumpMajor
)

func (b BumpType) String() string {
	switch b {
	case BumpPatch:
		return "patch"
	case BumpMinor:
		return "minor"
	case BumpMajor:
		return "major"
	default:
		return "none"
	}
}

// ConventionalCommit represents a parsed conventional commit message.
type ConventionalCommit struct {
	Type        string
	Scope       string
	Description string
	Body        string
	Footers     []Footer
	Breaking    bool
	Raw         string
	Hash        string
}

// Footer represents a git trailer / conventional commit footer.
type Footer struct {
	Token string
	Value string
}
