// Package cmd implements the CLI logic for rhad
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/opensourcecorp/rhad/logging"
	"github.com/spf13/cobra"
)

var (
	// rhadSrc is the root of the rhad source code itself -- useful for looking
	// up paths relative to the binary, etc
	rhadSrc string

	// rhadFile (defined in fsutils.go) represents the Rhadfile configuration,
	// which are collections of rhadConfigs
	rf rhadFile
	// rc rhadConfig

	rootCmd = &cobra.Command{
		Use:   "rhad",
		Short: "CI/CD task runner for OpenSourceCorp",
	}

	// When running rhad's tests
	isTesting bool
)

func init() {
	var ok bool

	rhadSrc, ok = os.LookupEnv("RHAD_SRC")
	if !ok {
		rhadSrc = "/home/rhad/rhad-src"
	}

	if os.Getenv("RHAD_TESTING") == "true" {
		isTesting = true
		fmt.Printf("RHAD_TESTING set to '%v', so will surpress further output\n", isTesting)
	}
}

func Execute() {
	var err error

	rf = readRhadfile()
	for path := range rf { // cfg := range rf {
		// Get rid of the brackets around INI section names
		for _, e := range []string{"[", "]"} {
			path = strings.ReplaceAll(path, e, "")
		}
		if path == "DEFAULT" { // INI's default, top-level section
			continue
		}
		err = os.Chdir(path)
		if err != nil {
			logging.Error("Could not set working directory to '%s' for rhad on startup", path)
			logging.Error(err.Error())
			os.Exit(1)
		}

		// rc = cfg

		err = rootCmd.Execute()
		if err != nil {
			os.Exit(1)
		}
	}
}
