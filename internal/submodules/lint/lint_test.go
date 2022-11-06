package lint

import (
	"testing"

	osc "github.com/opensourcecorp/go-common"
	"github.com/opensourcecorp/rhad/internal/tracker"
)

var (
	testLintRoot = "../../../testdata/linters"

	testTrackerData  []tracker.TrackerRecord
	testLintFailures int

	goodBadFiles = map[string][]string{
		"go": {"go_good.go", "go_bad.go"},
	}
)

func init() {
	osc.IsTesting = true
}

func TestLint(t *testing.T) {

	t.Run("Lint Go and pass", func(t *testing.T) {
		lintGo([]string{testLintRoot + "/" + goodBadFiles["go"][0]})
		testTrackerData = tracker.GetTrackerData()
		testLintFailures = tracker.CheckTrackerFailures(testTrackerData, "lint")
		if testLintFailures > 0 {
			t.Errorf(
				"\nlintGo failed on either the format diff-check or the lint itself, but should have succeeded -- tracker data below:\n%v",
				testTrackerData,
			)
		}
		// resets the tracker file on disk
		tracker.InitTracker()
	})

	t.Run("Lint Go and fail", func(t *testing.T) {
		lintGo([]string{testLintRoot + "/" + goodBadFiles["go"][1]})
		testTrackerData = tracker.GetTrackerData()
		testLintFailures = tracker.CheckTrackerFailures(testTrackerData, "lint")
		if testLintFailures == 0 {
			t.Errorf(
				"\nlintGo succeeded on either the format diff-check or the lint itself, but should have failed -- tracker data below:\n%v",
				testTrackerData,
			)
		}
		tracker.InitTracker()
	})
}
