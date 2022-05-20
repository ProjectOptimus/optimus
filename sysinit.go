package main

import (
	"flag"
	"os"
	"os/exec"
)

var (
	sysinitCommand = flag.NewFlagSet("sysinit", flag.ExitOnError)
	sysinitOSPtr   = sysinitCommand.String("os", "debian:unstable", `OS family & version tag of rhad host. Specified the same as OCI image tags. (default "debian:unstable")`)

	sysinitCLIConfig = sysinitCLIConfigStruct{
		*sysinitOSPtr,
	}
)

type sysinitCLIConfigStruct struct {
	os string
}

func sysinit(cfg sysinitCLIConfigStruct) {
	switch os.Args[2] {
	case "-h", "-help", "--help", "help":
		showHelp("sysinit", 0)
	}

	cmd := exec.Command("bash", "/usr/local/rhad/scripts/sysinit.sh")

	// Might need to set these so that the user has interactive access to any
	// prompts
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		ErrorLog.Println("Error running rhad's sysinit step -- see below for details")
		ErrorLog.Fatalln(err)
	}
}
