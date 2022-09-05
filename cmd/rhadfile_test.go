package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	osc "github.com/opensourcecorp/go-common"
)

var testRhadfilePath, _ = filepath.Abs("../testdata/config-file/Rhadfile")

func TestReadRhadfile(t *testing.T) {
	t.Run("Valid Rhadfile can be read", func(t *testing.T) {
		want := Rhadfile{
			Modules: map[string]rhadModule{
				".":          {Version: "v0.1.0"},
				"imgbuilder": {Version: "v1.0.0"},
				"datastore":  {Version: "v0.2.0"},
			},
		}
		got, _ := readRhadfile(testRhadfilePath)

		if !cmp.Equal(want, got) {
			t.Errorf("\nDid not match expected Rhadfile contents:\nwant: %v\ngot: %v", want, got)
		}
	})

	t.Run("Invalid Rhadfile fields throw a warning", func(t *testing.T) {
		var err error

		var buf bytes.Buffer

		badfilePath, err := filepath.Abs("../testdata/config-file/Rhadfile_invalid")
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Have to temporarily turn off testing & redirect output so the logs
		// get emitted, and to the right place
		osc.IsTesting = false
		osc.WarnLogger.SetOutput(&buf)
		readRhadfile(badfilePath)
		osc.WarnLogger.SetOutput(os.Stderr)
		osc.IsTesting = true

		logs := buf.String()
		if err != nil {
			t.Fatalf(err.Error())
		}

		match, err := regexp.MatchString(`WARN`, string(logs))
		if err != nil {
			t.Fatalf(err.Error())
		}
		if !match {
			t.Errorf("\nExpected a Rhadfile with invalid fields to throw a warning in the logs, but didn't see that. Log contents:\n%s", string(logs))
		}
	})
}
