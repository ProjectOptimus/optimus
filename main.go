package main

import (
	osc "github.com/opensourcecorp/go-common"
	"github.com/opensourcecorp/rhad/cmd"
)

func main() {
	osc.SetLoggerPrefixName("rhad")
	osc.InfoLog("Firing up rhad!")
	cmd.Execute()
}
