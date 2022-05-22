package main

import (
	"os"

	"github.com/opensourcecorp/rhad/help"
	"github.com/opensourcecorp/rhad/logging"
	"github.com/opensourcecorp/rhad/run"
)

func main() {
	if len(os.Args) < 2 {
		logging.Error("You must pass a subcommand to rhad")
		help.ShowHelp("main", 1)
	}

	switch os.Args[1] {
	case "help":
		help.ShowHelp("main", 0)
	case "run":
		run.Run()
	default:
		logging.Error("You must pass a valid subcommand to rhad")
		help.ShowHelp("main", 1)
	}
}
