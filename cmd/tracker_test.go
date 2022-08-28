package cmd

import (
	"errors"
	"os"
	"regexp"
	"testing"

	osc "github.com/opensourcecorp/go-common"
)

func TestInitTracker(t *testing.T) {
	if _, err := os.Stat(trackerPath); errors.Is(err, os.ErrNotExist) {
		t.Errorf("rhad tracker file does not exist -- it should have been created at init time")
	}
}

func TestWriteTrackerRecord(t *testing.T) {
	writeTrackerRecord(
		TrackerRecord{
			Type:     "lint",
			Subtype:  "fmt-diff-check",
			Language: "go",
			Tool:     "go fmt",
			Result:   "fail",
		},
	)
	trackerData, err := os.ReadFile(trackerPath)
	if err != nil {
		osc.FatalLog(err, "Couldn't read from rhad's tracker file during test")
	}
	hasResult, err := regexp.Match(`"result": ?"fail"`, trackerData)
	if err != nil {
		osc.FatalLog(err, "Bad regex spec during test")
	}
	if !hasResult {
		t.Errorf(
			"\nWritten tracker record does not have expected contents -- review file contents below:\n%v",
			string(trackerData),
		)
	}
}
