package tracker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	osc "github.com/opensourcecorp/go-common"
	"github.com/sirupsen/logrus"
)

var (
	trackerPath = path.Join(os.TempDir(), "oscar-tracker")
)

// Record collects and represents the JSON record structure of the tracker file
// contents
type Record struct {
	Type     string `json:"type"`     // What oscar subcommand this run represents (e.g. 'lint')
	Subtype  string `json:"subtype"`  // A subtype label, if applicable (e.g. 'fmt-diff-check'). If this is zeroed, then Type == Subtype
	Language string `json:"language"` // The language being processed, if applicable, such as during lint runs
	Tool     string `json:"tool"`     // What tool was used during the run (e.g. 'staticcheck')
	Target   string `json:"target"`   // The target, if applicable (e.g. 'aws')
	Result   string `json:"result"`   // Either 'pass' (or something informative) or 'fail'. 'fail' is what checkTrackerFailures looks for
}

func init() {
	InitTracker()
}

// String will let TrackerRecord print formatting be controllable, by satisying
// the built-in Stringer interface
func (t Record) String() string {
	record := fmt.Sprintf(
		osc.Dedent(`
			========
			Type: %s
			Subtype: %s
			Language: %s
			Tool: %s
			Target: %s
			Result: %s
			========
		`),
		t.Type,
		t.Subtype,
		t.Language,
		t.Tool,
		t.Target,
		t.Result,
	)
	return record
}

func InitTracker() {
	os.Remove(trackerPath)
	_, err := os.Create(trackerPath)
	if err != nil {
		logrus.Fatalf("Could not open or create oscar's tracker file at '%s': %v", trackerPath, err)
	}
}

func WriteRecord(t Record) {
	f, err := os.OpenFile(trackerPath, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatalf("Could not open oscar's tracker file at '%s' for writing: %v", trackerPath, err)
	}
	defer f.Close()

	recordBytes, err := json.Marshal(t)
	if err != nil {
		logrus.Fatalf("Could not marshal the following TrackerRecord to JSON for writing: %v -- error = %v", t, err)
	}

	if _, err = f.Write(recordBytes); err != nil {
		logrus.Fatalf("Could not write to oscar's tracker file: %v", err)
	}
}

func GetTrackerData() []Record {
	trackerFileBytes, err := os.ReadFile(trackerPath)
	if err != nil {
		logrus.Fatalf("Could not open oscar's tracker file at '%s': %v", trackerPath, err)
	}

	var trackerData []Record
	dec := json.NewDecoder(bytes.NewReader(trackerFileBytes))
	for {
		var record Record
		if err := dec.Decode(&record); err == io.EOF {
			break
		} else if err != nil {
			// TODO: *which* line couldn't be processed?
			logrus.Fatalf("Could not process one or more lines of oscar's tracker file: %v", err)
		}
		trackerData = append(trackerData, record)
	}
	return trackerData
}

func CheckTrackerFailures(trackerData []Record, trackerType string) int {
	var failures int
	for _, record := range trackerData {
		if record.Result == "fail" && record.Type == trackerType {
			failures++
		}
	}
	return failures
}