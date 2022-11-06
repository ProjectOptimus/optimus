package utils

import (
	"os"
	"testing"
)

func TestGetRhadSrc(t *testing.T) {
	want, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	os.Setenv("RHAD_SRC", want)
	got := GetRhadSrc()

	if want != got {
		t.Errorf("Want: %s, Got: %s\n", want, got)
	}
}

func TestCheckIsTesting(t *testing.T) {
	var isTesting bool

	t.Run("rhad knows it's being tested", func(t *testing.T) {
		os.Setenv("RHAD_TESTING", "true")
		isTesting = CheckIsTesting()

		if !isTesting {
			t.Errorf("isTesting should have been TRUE based on set env var")
		}
	})

	t.Run("rhad knows it's being tested", func(t *testing.T) {
		os.Setenv("RHAD_TESTING", "")
		isTesting = CheckIsTesting()

		if isTesting {
			t.Errorf("isTesting should have been FALSE based on unset/invalid env var")
		}
	})
}
