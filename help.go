package main

import (
	"fmt"
	"os"
	"strings"
)

// Printing spacing is kind of wonky, so be sure to check it in advance
func showHelp(subcommand string, exitCode int) {
	var helpText string

	switch subcommand {
	case "main":
		helpText = `rhad -- CI/CD task runner for OpenSourceCorp

Usage:
	rhad [-h|-help|--help|help] <subcommand>

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
	all		Run *all* the subcommands. The default subcommand.
	lint		Run linters.
	test		Run 'make test' target.
	build		Run 'make build' target.
	push		Run 'make push' target.
	deploy		Run deployment(s)`

	default:
		ErrorLog.Fatalln("How did you even get here?")
	}

	fmt.Println(strings.TrimSpace(helpText))
	os.Exit(exitCode)
}
