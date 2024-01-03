package utils

import (
	"os"
	"testing"
)

func TestGetOscarSrc(t *testing.T) {
	want, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	os.Setenv("OSCAR_SRC", want)
	got := GetOscarSrc()

	if want != got {
		t.Errorf("Want: %s, Got: %s\n", want, got)
	}
}
