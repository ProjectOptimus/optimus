// Package help consolidates helptext strings for the CLI commands
package help

import (
	"fmt"
	"os"
	"strings"

	"github.com/opensourcecorp/rhadamanthus/logging"
)

// Columnar spacing is kind of wonky at print-time, so be sure to check it in
// advance
func ShowHelp(subcommand string, exitCode int) {
	var helpText string

	switch subcommand {
	case "main":
		helpText = `rhad -- CI/CD task runner for OpenSourceCorp

Usage:
	rhad [help] <subcommand>

Subcommands:
	run     <[all] | lint | test | build | push | deploy>
	sysinit [-os <name:tag>]`

	case "sysinit":
		helpText = `rhad sysinit -- initialize rhad's host system

Usage:
	[-h|-help|--help]	Show low-detail help
	[help]			Show more detailed help
	[-os <osNametag>]	OS family & version tag of rhad host. Specified the same as OCI image tags. (default "debian:unstable")`

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

	default:
		logging.Error("How did you even get here?")
	}

	fmt.Println(strings.TrimSpace(helpText))
	os.Exit(exitCode)
}
