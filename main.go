package main

import (
	"os"

	"github.com/opensourcecorp/rhadamanthus/help"
	"github.com/opensourcecorp/rhadamanthus/logging"
	"github.com/opensourcecorp/rhadamanthus/run"
	"github.com/opensourcecorp/rhadamanthus/sysinit"
)

func main() {
	if len(os.Args) < 2 {
		logging.Error("You must pass a subcommand to rhad")
		help.ShowHelp("main", 1)
	}

	switch os.Args[1] {
	case "-h", "-help", "--help", "help":
		help.ShowHelp("main", 0)
	case "sysinit":
		sysinit.Sysinit()
	case "run":
		run.Run()
	default:
		logging.Error("You must pass a valid subcommand to rhad")
		help.ShowHelp("main", 1)
	}
}
