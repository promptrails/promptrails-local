package handler

import (
	"strings"
	"testing"
)

func TestNormalizeVFSPath_Valid(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"/a/b/c", "/a/b/c"},
		{"foo", "/foo"},          // relative paths are made absolute
		{"/", "/"},               // root
		{"//a", "/a"},            // duplicate slashes cleaned
		{"/a/", "/a"},            // trailing slash cleaned
		{"/a/./b", "/a/b"},       // "." segments cleaned away
		{"/a/b.txt", "/a/b.txt"}, // dotted file name kept
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			got, err := normalizeVFSPath(tc.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("normalizeVFSPath(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestNormalizeVFSPath_Rejected(t *testing.T) {
	longName := strings.Repeat("a", vfsMaxNameLen+1)
	deep := "/" + strings.TrimSuffix(strings.Repeat("d/", vfsMaxDepth+1), "/")

	tests := []struct {
		name string
		in   string
	}{
		{"empty", ""},
		{"nul byte", "/a/\x00/b"},
		{"parent traversal", "/a/../b"},
		{"leading traversal", "../etc/passwd"},
		{"windows-style traversal", "/a/..\\b/../c/.."}, // contains ".." segment
		{"too deep", deep},
		{"name too long", "/" + longName},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := normalizeVFSPath(tc.in); err == nil {
				t.Errorf("normalizeVFSPath(%q) = nil error, want rejection", tc.in)
			}
		})
	}
}

func TestNormalizeVFSPath_DepthBoundary(t *testing.T) {
	// Exactly vfsMaxDepth segments must be accepted.
	atLimit := "/" + strings.TrimSuffix(strings.Repeat("d/", vfsMaxDepth), "/")
	if _, err := normalizeVFSPath(atLimit); err != nil {
		t.Errorf("path at max depth should be accepted, got %v", err)
	}
}

func TestSplitVFSPath(t *testing.T) {
	tests := []struct {
		in         string
		wantParent string
		wantName   string
	}{
		{"/", "", ""},
		{"", "", ""},
		{"/foo", "/", "foo"},
		{"/a/b", "/a", "b"},
		{"/a/b/c", "/a/b", "c"},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			parent, name := splitVFSPath(tc.in)
			if parent != tc.wantParent || name != tc.wantName {
				t.Errorf("splitVFSPath(%q) = (%q, %q), want (%q, %q)",
					tc.in, parent, name, tc.wantParent, tc.wantName)
			}
		})
	}
}

func TestGlobMatch(t *testing.T) {
	tests := []struct {
		pattern string
		s       string
		want    bool
	}{
		{"*", "anything", true},
		{"*.txt", "notes.txt", true},
		{"*.txt", "notes.md", false},
		{"src*", "src/main.go", true},
		{"src*", "lib/main.go", false},
		{"*main*", "src/main.go", true},
		{"*main*", "src/util.go", false},
		{"exact.go", "exact.go", true},
		{"exact.go", "other.go", false},
	}
	for _, tc := range tests {
		t.Run(tc.pattern+"_"+tc.s, func(t *testing.T) {
			if got := globMatch(tc.pattern, tc.s); got != tc.want {
				t.Errorf("globMatch(%q, %q) = %v, want %v", tc.pattern, tc.s, got, tc.want)
			}
		})
	}
}
