// Package sysinit provides the system initialization management for rhad host
package sysinit

import (
	"flag"
	"os"
	"os/exec"

	"github.com/opensourcecorp/rhadamanthus/help"
	"github.com/opensourcecorp/rhadamanthus/logging"
)

type cliConfig struct {
	osTag         *string
	yesRunSysinit *bool
}

func Sysinit() {
	if len(os.Args) >= 3 {
		if os.Args[2] == "help" {
			help.ShowHelp("sysinit", 0)
		}
	}

	command := flag.NewFlagSet("sysinit", flag.ExitOnError)
	osTag := command.String("os", "debian:unstable", `OS family & version tag of rhad host. Specified the same as OCI image tags. (default "debian:unstable")`)
	yesRunSysinit := command.Bool("yes", false, "Allow sysinit to proceed. Since this can easily clutter up your actual development host, this is required")
	command.Parse(os.Args[2:])

	cfg := cliConfig{
		osTag,
		yesRunSysinit,
	}

	if *cfg.yesRunSysinit {
		cmd := exec.Command("bash", "./scripts/sysinit.sh")

		// Might need to set these so that the user has interactive access to any
		// prompts
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			logging.Error(err.Error())
			logging.Error("Error running rhad's sysinit step -- see above below for details")
			os.Exit(1)
		}
	} else {
		logging.Error("rhad's sysinit can be locally destructive, so you must pass the '-yes' flag to confirm you know what you are doing!")
		os.Exit(1)
	}
}
