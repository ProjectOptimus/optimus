// Package cmd implements the CLI logic for rhad
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	osc "github.com/opensourcecorp/go-common"
	"github.com/spf13/cobra"
)

var (
	// rhadSrc is the root of the rhad source code itself -- useful for looking
	// up paths relative to the binary, etc
	rhadSrc string

	rf Rhadfile

	rootCmd = &cobra.Command{
		Use:   "rhad",
		Short: "CI/CD task runner for OpenSourceCorp",
	}

	// When running rhad's tests
	isTesting bool
)

func init() {
	var ok bool
	var err error

	rhadSrc, ok = os.LookupEnv("RHAD_SRC")
	if !ok {
		rhadSrc, err = filepath.Abs(".")
		if err != nil {
			osc.FatalLog(err, "Error trying to determine default absolute filepath to rhad's sourcecode root")
		}
	} else {
		_, err = os.Lstat(rhadSrc)
		if err != nil {
			osc.FatalLog(err, "Env var 'RHAD_SRC' was provided, but set to a nonexistant directory")
		}
	}

	if os.Getenv("RHAD_TESTING") == "true" {
		isTesting = true
		fmt.Printf("RHAD_TESTING set to '%v', so will surpress further output\n", isTesting)
	}
}

func Execute() {
	var err error

	rf, _ = readRhadfile()
	for module := range rf.Modules {
		// Rhadfile section
		// Rhadfile section keys are subdirectory names in the tree, so we can
		// also use them as path names where we need to
		modulePath, err := filepath.Abs(module)
		if err != nil {
			osc.FatalLog(err, "Error when construction absolute filepath to provided module '%s'", module)
		}
		err = os.Chdir(modulePath)
		if err != nil {
			osc.FatalLog(err, "Could not set working directory to '%s' for rhad on startup", module)
		}

		err = rootCmd.Execute()
		if err != nil {
			osc.FatalLog(err, "Unhandled error when executing rhad subcommands, caught at top-level")
		}
	}
	err = os.Chdir(rhadSrc)
	if err != nil {
		osc.FatalLog(err, "Could not reset working directory to rhad root '%s' for rhad on finish", rhadSrc)
	}
}

// testSysinit can be used to run before functions that are making
// osc.Syscall.Exec()s, so they hopefully catch runtime errors earlier
func testSysinit() {
	if _, err := os.Stat(rhadSrc); errors.Is(err, os.ErrNotExist) {
		rhadSrc = "."
	}

	sc := osc.Syscall{
		CmdLine: []string{"bash", rhadSrc + "/scripts/sysinit.sh", "test"},
	}
	sc.Exec()
	if !sc.Ok {
		os.Exit(1)
	}
}
