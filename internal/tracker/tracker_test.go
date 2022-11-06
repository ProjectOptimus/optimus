package tracker

import (
	"errors"
	"os"
	"regexp"
	"testing"

	osc "github.com/opensourcecorp/go-common"
)

var (
	// Used when no specific TrackerRecord is needed; we just need something to show up in the tracker file
	basicTestTrackerRecord = TrackerRecord{
		Type:     "lint",
		Subtype:  "fmt-diff-check",
		Language: "go",
		Tool:     "gofmt",
		Result:   "fail",
	}
)

func init() {
	osc.IsTesting = true
}

func TestInitTracker(t *testing.T) {
	if _, err := os.Stat(trackerPath); errors.Is(err, os.ErrNotExist) {
		t.Errorf("rhad tracker file does not exist -- it should have been created at init time")
	}
}

func TestWriteTrackerRecord(t *testing.T) {
	WriteTrackerRecord(basicTestTrackerRecord)
	// Raw bytes read instead of getTrackerData(), so we can a) test the bytes
	// written and b) debug if getTrackerData() tests fail
	trackerFileBytes, err := os.ReadFile(trackerPath)
	if err != nil {
		osc.FatalLog(err, "Couldn't read from rhad's tracker file during test")
	}
	hasResult, err := regexp.Match(`"result": ?"fail"`, trackerFileBytes)
	if err != nil {
		osc.FatalLog(err, "Bad regex spec during test")
	}
	if !hasResult {
		t.Errorf(
			"\nWritten tracker record does not have expected contents -- review file contents below:\n%v",
			string(trackerFileBytes),
		)
	}

	InitTracker()
}

func TestGetTrackerData(t *testing.T) {
	WriteTrackerRecord(basicTestTrackerRecord)
	trackerData := GetTrackerData()
	if trackerData[0].Type != "lint" {
		t.Errorf("Expected tracker record field 'Type' to be 'lint', but got '%v'", trackerData[0].Type)
	}
	if trackerData[0].Subtype != "fmt-diff-check" {
		t.Errorf("Expected tracker record field 'Subtype' to be 'fmt-diff-check', but got '%v'", trackerData[0].Subtype)
	}
	if trackerData[0].Language != "go" {
		t.Errorf("Expected tracker record field 'Language' to be 'go', but got '%v'", trackerData[0].Language)
	}
	if trackerData[0].Tool != "gofmt" {
		t.Errorf("Expected tracker record field 'Tool' to be 'gofmt', but got '%v'", trackerData[0].Tool)
	}
	if trackerData[0].Target != "" {
		t.Errorf("Expected tracker record field 'Target' to be an empty string, but got '%v'", trackerData[0].Target)
	}
	if trackerData[0].Result != "fail" {
		t.Errorf("Expected tracker record field 'Result' to be 'fail', but got '%v'", trackerData[0].Result)
	}

	InitTracker()
}

func TestCheckTrackerFailures(t *testing.T) {
	WriteTrackerRecord(basicTestTrackerRecord)
	trackerData := GetTrackerData()
	gotFailures := CheckTrackerFailures(trackerData, "lint")
	wantFailures := 1
	if gotFailures != wantFailures {
		t.Errorf("\nExpected to find %d failures in the tracker file, but found %d -- contents below:\n%s", wantFailures, gotFailures, trackerData)
	}

	InitTracker()
}
