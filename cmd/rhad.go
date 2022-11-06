// Package cmd implements the CLI logic for rhad
package cmd

import (
	"os"
	"path/filepath"

	osc "github.com/opensourcecorp/go-common"
	"github.com/spf13/cobra"

	"github.com/opensourcecorp/rhad/internal/rhadfile"
	"github.com/opensourcecorp/rhad/internal/submodules/lint"
	"github.com/opensourcecorp/rhad/internal/utils"
)

var (
	// rhadSrc is the root of the rhad source code itself -- useful for looking
	// up paths relative to the binary, etc
	rhadSrc string

	rf rhadfile.Rhadfile

	rootCmd = &cobra.Command{
		Use:   "rhad",
		Short: "rhad: the CI/CD task runner for OpenSourceCorp",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			osc.FatalLog(nil, "rhad requires a subcommand")
		},
	}
)

func init() {
	rhadSrc = utils.GetRhadSrc()

	rootCmd.AddCommand(lint.LintCmd)
}

func Execute() {
	var err error

	rf, _ = rhadfile.ReadRhadfile()
	for module := range rf.Modules {
		// Rhadfile section keys are relative-path subdirectory names in the
		// tree, so we can also use them as path names where we need to
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
