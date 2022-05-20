package main

import (
	"log"
	"os"
)

var (
	// Custom logger funcs -- actually created in init(), so they're visible everywhere else
	DebugLog *log.Logger
	InfoLog  *log.Logger
	WarnLog  *log.Logger
	ErrorLog *log.Logger
)

func init() {
	// The usage of bitwise OR here seems to be called "bitmask flagging", since
	// the log output option needs to be an integer and ORing their named bits
	// gives you a single integer result
	DebugLog = log.New(os.Stderr, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLog = log.New(os.Stderr, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	if len(os.Args) < 2 {
		ErrorLog.Println("You must pass a subcommand to rhad")
		showHelp("main", 1)
	}

	switch os.Args[1] {
	case "-h", "-help", "--help", "help":
		showHelp("main", 0)
	case "sysinit":
		sysinitCommand.Parse(os.Args[2:])
	case "run":
		runCommand.Parse(os.Args[2:])
	default:
		ErrorLog.Println("You must pass a valid subcommand to rhad")
		showHelp("main", 1)
	}

	if sysinitCommand.Parsed() {
		sysinit(sysinitCLIConfig)
	}

	if runCommand.Parsed() {
		run(runCLIConfig)
	}
}
