package cmd

// Note: these tests are relying on the global 'lintFailures' failureMap as
// defined in lint.go

import (
	"testing"
)

var (
	testLintRoot = "../testdata/linters"

	testTrackerData  []TrackerRecord
	testLintFailures int

	goodBadFiles = map[string][]string{
		"go": {"go_good.go", "go_bad.go"},
	}
)

func TestLintGo(t *testing.T) {
	lintGo([]string{testLintRoot + "/" + goodBadFiles["go"][0]})
	testTrackerData = getTrackerData()
	testLintFailures = checkTrackerFailures(testTrackerData, "lint")
	if testLintFailures > 0 {
		t.Errorf(
			"\nlintGo failed on either the format diff-check or the lint itself, but should have succeeded -- tracker data below:\n%v",
			testTrackerData,
		)
	}
	// resets the tracker file on disk
	initTracker()

	lintGo([]string{testLintRoot + "/" + goodBadFiles["go"][1]})
	testTrackerData = getTrackerData()
	testLintFailures = checkTrackerFailures(testTrackerData, "lint")
	if testLintFailures == 0 {
		t.Errorf(
			"\nlintGo succeeded on either the format diff-check or the lint itself, but should have failed -- tracker data below:\n%v",
			testTrackerData,
		)
	}
	initTracker()
}
