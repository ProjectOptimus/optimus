package fsutils

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	osc "github.com/opensourcecorp/go-common"
)

var testingRoot = "../../testdata/fsutils"

func init() {
	osc.IsTesting = true
}

func TestGetAllFiles(t *testing.T) {
	var want, got []fileData

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

	want = []fileData{
		{filepath.Join(testingRoot, "test.go"), false},
		{filepath.Join(testingRoot, "test.json"), false},
	}
	got = GetRelevantFiles(testingRoot, `(\.go|\.json)`)
	if !cmp.Equal(want, got) {
		t.Errorf("\nWant: %v\nGot:  %v\n", want, got)
	}
}
