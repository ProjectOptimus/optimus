package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path"

	osc "github.com/opensourcecorp/go-common"
)

var (
	trackerPath = path.Join(os.TempDir(), "rhad-tracker")
)

// Tracker collects and represents the JSON record structure of the tracker file
// contents
type TrackerRecord struct {
	Type     string `json:"type"`     // What rhad subcommand this run represents (e.g. 'lint')
	Subtype  string `json:"subtype"`  // A subtype label, if applicable (e.g. 'fmt-diff-check'). If this is zeroed, then Type == Subtype
	Language string `json:"language"` // The language being processed, if applicable, such as during lint runs
	Target   string `json:"target"`   // The target, if applicable (e.g. 'aws')
	Tool     string `json:"tool"`     // What tool was used during the run (e.g. 'staticcheck')
	Result   string `json:"result"`   // Either 'pass' (or something informative) or 'fail'. 'fail' is what checkTrackerFailures looks for
}

func init() {
	initTracker()
}

func initTracker() {
	os.Remove(trackerPath)
	_, err := os.Create(trackerPath)
	if err != nil {
		osc.FatalLog(err, "Could not open or create rhad's tracker file at '%s'", trackerPath)
	}
}

func writeTrackerRecord(t TrackerRecord) {
	f, err := os.OpenFile(trackerPath, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		osc.FatalLog(err, "Could not open rhad's tracker file at '%s' for writing", trackerPath)
	}
	defer f.Close()

	recordBytes, err := json.Marshal(t)
	if err != nil {
		osc.FatalLog(err, "Could not marshal the following TrackerRecord to JSON for writing: %v", t)
	}

	if _, err = f.Write(recordBytes); err != nil {
		osc.FatalLog(err, "Could not write to rhad's tracker file")
	}
}

func getTrackerData() []TrackerRecord {
	trackerFileBytes, err := os.ReadFile(trackerPath)
	if err != nil {
		osc.FatalLog(err, "Could not open rhad's tracker file at '%s'", trackerPath)
	}

	var trackerData []TrackerRecord
	dec := json.NewDecoder(bytes.NewReader(trackerFileBytes))
	for {
		var record TrackerRecord
		if err := dec.Decode(&record); err == io.EOF {
			break
		} else if err != nil {
			// TODO: *which* line couldn't be processed?
			osc.FatalLog(err, "Could not process one or more lines of rhad's tracker file")
		}
		trackerData = append(trackerData, record)
	}
	return trackerData
}

func checkTrackerFailures(trackerData []TrackerRecord, trackerType string) int {
	var failures int
	for _, record := range trackerData {
		if record.Result == "fail" && record.Type == trackerType {
			failures++
		}
	}
	return failures
}
