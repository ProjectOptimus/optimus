// Package oscar implements the CLI logic for oscar
package oscar

import (
	"os"
	"path/filepath"

	"github.com/opensourcecorp/oscar/internal/oscarfile"
	"github.com/opensourcecorp/oscar/internal/subroutines/test"
	"github.com/opensourcecorp/oscar/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// oscarSrc is the root of the oscar source code itself -- useful for looking
	// up paths relative to the binary, etc
	oscarSrc string

	oscfile oscarfile.Oscarfile

	rootCmd = &cobra.Command{
		Use:   "oscar",
		Short: "oscar: the OpenSourceCorp Automation Runner",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			logrus.Fatal("oscar requires a subcommand")
		},
	}
)

func init() {
	oscarSrc = utils.GetOscarSrc()

	// rootCmd.AddCommand(lint.LintCmd)
	rootCmd.AddCommand(test.Cmd)
}

func Execute() {
	var err error

	oscfile, _ = oscarfile.Read()
	for module := range oscfile.Modules {
		// Oscarfile section keys are relative-path subdirectory names in the
		// tree, so we can also use them as path names where we need to
		modulePath, err := filepath.Abs(module)
		if err != nil {
			logrus.Fatalf("constructing absolute filepath to provided module '%s': %v", module, err)
		}
		err = os.Chdir(modulePath)
		if err != nil {
			logrus.Fatalf("could not set working directory to '%s' for oscar on startup: %v", module, err)
		}

		err = rootCmd.Execute()
		if err != nil {
			logrus.Fatalf("unhandled error when executing oscar subcommands, caught at top-level: %v", err)
		}
	}
	err = os.Chdir(oscarSrc)
	if err != nil {
		logrus.Fatalf("could not reset working directory to oscar root '%s' for oscar on finish: %v", oscarSrc, err)
	}
}
