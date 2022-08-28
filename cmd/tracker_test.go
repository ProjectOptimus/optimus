package cmd

import (
	"errors"
	"os"
	"testing"
)

func TestInitTracker(t *testing.T) {
	if _, err := os.Stat(trackerPath); errors.Is(err, os.ErrNotExist) {
		t.Errorf("rhad tracker file does not exist -- it should have been created at init time")
	}
}
