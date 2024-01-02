// Package cmd implements the CLI logic for oscar
package cmd

import (
	"os"
	"path/filepath"

	osc "github.com/opensourcecorp/go-common"
	"github.com/spf13/cobra"

	"github.com/opensourcecorp/oscar/internal/oscarfile"
	"github.com/opensourcecorp/oscar/internal/utils"
)

var (
	// oscarSrc is the root of the oscar source code itself -- useful for looking
	// up paths relative to the binary, etc
	oscarSrc string

	rf oscarfile.Oscarfile

	rootCmd = &cobra.Command{
		Use:   "oscar",
		Short: "oscar: the CI/CD task runner for OpenSourceCorp",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			osc.FatalLog(nil, "oscar requires a subcommand")
		},
	}
)

func init() {
	oscarSrc = utils.GetOscarSrc()

	// rootCmd.AddCommand(lint.LintCmd)
}

func Execute() {
	var err error

	rf, _ = oscarfile.ReadOscarfile()
	for module := range rf.Modules {
		// Oscarfile section keys are relative-path subdirectory names in the
		// tree, so we can also use them as path names where we need to
		modulePath, err := filepath.Abs(module)
		if err != nil {
			osc.FatalLog(err, "Error when construction absolute filepath to provided module '%s'", module)
		}
		err = os.Chdir(modulePath)
		if err != nil {
			osc.FatalLog(err, "Could not set working directory to '%s' for oscar on startup", module)
		}

		err = rootCmd.Execute()
		if err != nil {
			osc.FatalLog(err, "Unhandled error when executing oscar subcommands, caught at top-level")
		}
	}
	err = os.Chdir(oscarSrc)
	if err != nil {
		osc.FatalLog(err, "Could not reset working directory to oscar root '%s' for oscar on finish", oscarSrc)
	}
}
