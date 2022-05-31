package cmd

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetVersion(t *testing.T) {
	var want, got, s string // 's' is the version string from Rhadfile

	s = "1.0.0.0"
	want = "v1.0.0"
	got = getVersion(s)
	if want != got {
		t.Errorf("Expected version string '%s' to become '%s', but got '%s'\n", s, want, got)
	}

	s = "1.0-alpha"
	want = "v1.0.0-alpha"
	got = getVersion(s)
	if want != got {
		t.Errorf("Expected version string '%s' to become '%s', but got '%s'\n", s, want, got)
	}

	s = "1.0.2+abc"
	want = "v1.0.2+abc"
	got = getVersion(s)
	if want != got {
		t.Errorf("Expected version string '%s' to become '%s', but got '%s'\n", s, want, got)
	}

	// No change on ideally-provided semver format
	s = "v1.1.9-prebeta1+abc"
	want = s
	got = getVersion(s)
	if want != got {
		t.Errorf("Expected version string '%s' to become '%s', but got '%s'\n", s, want, got)
	}
}

func TestGetAllFiles(t *testing.T) {
	var want, got []fileData

	testingRoot := "../testdata/fsutils"

	// Make sure these are lexically sorted, all dirs & contents and then
	// top-level files
	want = []fileData{
		{testingRoot, true},
		{filepath.Join(testingRoot, "dir"), true},
		{filepath.Join(testingRoot, "dir/idc.txt"), false},
		{filepath.Join(testingRoot, "test.go"), false},
		{filepath.Join(testingRoot, "test.json"), false},
		{filepath.Join(testingRoot, "test.md"), false},
		{filepath.Join(testingRoot, "test.py"), false},
		{filepath.Join(testingRoot, "test.sh"), false},
	}
	got = getAllFiles(testingRoot)
	if !cmp.Equal(want, got) {
		t.Errorf("\nWant: %v\nGot:  %v\n", want, got)
	}
}

func TestGetRelevantFiles(t *testing.T) {
	var want []fileData
	var got []fileData

	testingRoot := "../testdata/fsutils"

	want = []fileData{
		{filepath.Join(testingRoot, "test.go"), false},
		{filepath.Join(testingRoot, "test.json"), false},
	}
	got = getRelevantFiles(testingRoot, `(\.go|\.json)`)
	if !cmp.Equal(want, got) {
		t.Errorf("\nWant: %v\nGot:  %v\n", want, got)
	}
}

func TestCleanSectionName(t *testing.T) {
	var want, got string

	want = "[.]"
	got = cleanSectionName("[.]")
	if want != got {
		t.Errorf("Want: %v, Got: %v\n", want, got)
	}

	want = "[.]"
	got = cleanSectionName("[  .  ]")
	if want != got {
		t.Errorf("Want: %v, Got: %v\n", want, got)
	}

	want = "[dir1]"
	got = cleanSectionName("[  ./dir1    ]")
	if want != got {
		t.Errorf("Want: %v, Got: %v\n", want, got)
	}
}
