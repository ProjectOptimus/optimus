package oscarfile

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
)

var testOscarfilePath, _ = filepath.Abs("../../testdata/config_file/Oscarfile")

func TestRead(t *testing.T) {
	t.Run("Valid Oscarfile can be read", func(t *testing.T) {
		want := Oscarfile{
			Modules: map[string]oscarModule{
				".":          {Version: "v0.1.0"},
				"imgbuilder": {Version: "v1.0.0"},
				"datastore":  {Version: "v0.2.0"},
			},
		}
		got, _ := Read(testOscarfilePath)

		if !cmp.Equal(want, got) {
			t.Errorf("\nDid not match expected Oscarfile contents:\nwant: %v\ngot: %v", want, got)
		}
	})

	t.Run("Invalid Oscarfile fields throw a warning", func(t *testing.T) {
		var err error

		var buf bytes.Buffer

		badfilePath, err := filepath.Abs("../testdata/config_file/Oscarfile_invalid")
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Have to temporarily redirect output so the logs get emitted to the
		// right place
		logrus.SetOutput(&buf)
		Read(badfilePath)
		logrus.SetOutput(os.Stderr)

		logs := buf.String()
		if err != nil {
			t.Fatalf(err.Error())
		}

		match, err := regexp.MatchString(`warning`, string(logs))
		if err != nil {
			t.Fatalf(err.Error())
		}
		if !match {
			t.Errorf("\nExpected a Oscarfile with invalid fields to throw a warning in the logs, but didn't see that. Log contents:\n%s", string(logs))
		}
	})
}
