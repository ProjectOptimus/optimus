// Package run handles the logic of the `rhad run` processes
package run

import (
	"flag"
	"os"

	"github.com/opensourcecorp/rhad/help"
	"github.com/opensourcecorp/rhad/logging"
)

type cliOptsStruct struct {
	path          *string
	ignorePattern *string
}

type cfgFileData map[string]any

type rhadConfig struct {
	cfgFileData cfgFileData
	cliOpts     cliOptsStruct
}

func Run() {
	if len(os.Args) < 3 {
		logging.Error("You must pass a subcommand to rhad run")
		help.ShowHelp("run", 1)
	}

	command := flag.NewFlagSet("run", flag.ExitOnError)
	path := command.String("path", ".", "Path for rhad to act against")
	ignorePattern := command.String("ignore-pattern", "", "Regex pattern of files and/or directories to ignore when running rhad")
	command.Parse(os.Args[2:])

	if _, err := os.Lstat(*path); err != nil {
		logging.Error("Provided path does not exist!")
		os.Exit(1)
	}

	cliOpts := cliOptsStruct{
		path,
		ignorePattern,
	}

	cfgFileData := readConfig()

	cfg := rhadConfig{
		cfgFileData,
		cliOpts,
	}

	if os.Args[2] == "help" {
		help.ShowHelp("run", 0)
	}

	switch os.Args[len(os.Args)-1] {
	case "all":
		Lint(cfg)
		Test(cfg)
		Build(cfg)
		Push(cfg)
		Deploy(cfg)
	case "lint":
		Lint(cfg)
	case "test":
		Test(cfg)
	case "build":
		Build(cfg)
	case "push":
		Push(cfg)
	case "deploy":
		Deploy(cfg)
	default:
		logging.Error("Invalid subcommand provided to rhad run")
	}
}
