package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	runCommand = flag.NewFlagSet("run", flag.ExitOnError)

	runPathPtr          = runCommand.String("path", ".", "Path for rhad to act against")
	runIgnorePatternPtr = runCommand.String("ignore-pattern", "", "Regex pattern of files and/or directories to ignore when running rhad")

	runCLIConfig = runCLIConfigStruct{
		*runPathPtr,
		*runIgnorePatternPtr,
	}
)

type runCLIConfigStruct struct {
	path          string
	ignorePattern string
}

func run(cfg runCLIConfigStruct) {
	switch os.Args[2] {
	case "-h", "-help", "--help", "help":
		showHelp("run", 0)
	}

	WarnLog.Println("run not yet implemented")
	fmt.Println(cfg)
}
