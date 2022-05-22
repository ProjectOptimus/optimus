// Package help consolidates helptext strings for the CLI commands
package help

import (
	"fmt"
	"os"
	"strings"

	"github.com/opensourcecorp/rhad/logging"
)

// ShowHelp prints CLI help strings and then exits accordingly
func ShowHelp(subcommand string, exitCode int) {
	var helpText string

	// Columnar spacing is kind of wonky at print-time, so be sure to check each
	// of these in advance of publishing
	switch subcommand {

	//
	case "main":
		helpText = `rhad -- CI/CD task runner for OpenSourceCorp

Usage:
	rhad [help] <subcommand>

Subcommands:
	run     <[all] | lint | test | build | push | deploy>`

	//
	case "run":
		helpText = `rhad run -- run rhad CI/CD stages over directory

Usage:
	[-h|-help|--help]	Show low-detail help
	[help]			Show more detailed help

Subcommands:
	all		Run *all* the subcommands. This is the default.
	lint		Run linters.
	test		Run 'make test' target.
	build		Run 'make build' target.
	push		Run 'make push' target.
	deploy		Run deployment(s)`

	//
	default:
		logging.Error("How did you even get here?")
		exitCode = 1
	}

	fmt.Println(strings.TrimSpace(helpText))
	os.Exit(exitCode)
}
