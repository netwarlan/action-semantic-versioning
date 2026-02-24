package semver

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input   string
		want    Version
		wantErr bool
	}{
		{"v1.2.3", Version{Major: 1, Minor: 2, Patch: 3, Prefix: "v"}, false},
		{"1.2.3", Version{Major: 1, Minor: 2, Patch: 3}, false},
		{"v0.0.0", Version{Prefix: "v"}, false},
		{"v1.0.0-alpha.1", Version{Major: 1, Prefix: "v", Prerelease: "alpha.1"}, false},
		{"v1.0.0+build.123", Version{Major: 1, Prefix: "v", Metadata: "build.123"}, false},
		{"v1.0.0-beta+build", Version{Major: 1, Prefix: "v", Prerelease: "beta", Metadata: "build"}, false},
		{"v10.20.30", Version{Major: 10, Minor: 20, Patch: 30, Prefix: "v"}, false},
		// Invalid
		{"v1.2", Version{}, true},
		{"v1.2.3.4", Version{}, true},
		{"abc", Version{}, true},
		{"", Version{}, true},
		{"v-1.2.3", Version{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := Parse(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q, got %v", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("Parse(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}

func TestBump(t *testing.T) {
	v := Version{Major: 1, Minor: 2, Patch: 3, Prefix: "v"}

	major := v.BumpMajor()
	if major.String() != "v2.0.0" {
		t.Errorf("BumpMajor() = %s, want v2.0.0", major.String())
	}

	minor := v.BumpMinor()
	if minor.String() != "v1.3.0" {
		t.Errorf("BumpMinor() = %s, want v1.3.0", minor.String())
	}

	patch := v.BumpPatch()
	if patch.String() != "v1.2.4" {
		t.Errorf("BumpPatch() = %s, want v1.2.4", patch.String())
	}
}

func TestBumpClearsPrerelease(t *testing.T) {
	v := Version{Major: 1, Minor: 2, Patch: 3, Prefix: "v", Prerelease: "alpha"}

	if got := v.BumpMajor().String(); got != "v2.0.0" {
		t.Errorf("BumpMajor() = %s, want v2.0.0", got)
	}
	if got := v.BumpMinor().String(); got != "v1.3.0" {
		t.Errorf("BumpMinor() = %s, want v1.3.0", got)
	}
	if got := v.BumpPatch().String(); got != "v1.2.4" {
		t.Errorf("BumpPatch() = %s, want v1.2.4", got)
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		v    Version
		want string
	}{
		{Version{Major: 1, Minor: 2, Patch: 3, Prefix: "v"}, "v1.2.3"},
		{Version{Major: 1, Minor: 2, Patch: 3}, "1.2.3"},
		{Version{Major: 1, Prefix: "v", Prerelease: "alpha"}, "v1.0.0-alpha"},
		{Version{Major: 1, Prefix: "v", Metadata: "build"}, "v1.0.0+build"},
		{Version{Major: 1, Prefix: "v", Prerelease: "rc.1", Metadata: "build"}, "v1.0.0-rc.1+build"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"v1.0.0", "v1.0.0", 0},
		{"v1.0.0", "v2.0.0", -1},
		{"v2.0.0", "v1.0.0", 1},
		{"v1.0.0", "v1.1.0", -1},
		{"v1.1.0", "v1.0.0", 1},
		{"v1.0.0", "v1.0.1", -1},
		{"v1.0.1", "v1.0.0", 1},
		// Prerelease has lower precedence
		{"v1.0.0", "v1.0.0-alpha", 1},
		{"v1.0.0-alpha", "v1.0.0", -1},
		{"v1.0.0-alpha", "v1.0.0-beta", -1},
		{"v1.0.0-beta", "v1.0.0-alpha", 1},
		{"v1.0.0-alpha", "v1.0.0-alpha", 0},
	}

	for _, tt := range tests {
		t.Run(tt.a+"_vs_"+tt.b, func(t *testing.T) {
			a, _ := Parse(tt.a)
			b, _ := Parse(tt.b)
			if got := a.Compare(b); got != tt.want {
				t.Errorf("Compare(%s, %s) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestPrefixPreserved(t *testing.T) {
	v, _ := Parse("v1.2.3")
	bumped := v.BumpMinor()
	if bumped.Prefix != "v" {
		t.Errorf("prefix not preserved: got %q", bumped.Prefix)
	}
	if bumped.String() != "v1.3.0" {
		t.Errorf("String() = %s, want v1.3.0", bumped.String())
	}
}
