package main

import (
	"fmt"

	osc "github.com/opensourcecorp/go-common"
	"github.com/opensourcecorp/rhad/cmd"
)

func main() {
	osc.SetLoggerPrefixName("rhad")

	osc.InfoLog("Firing up rhad!")

	banner := osc.Dedent(`
		  ______________________________
		 /                             /|
		|/   __   |      __      __ | |/|
		|/  /  \  |__   /  \    /  \| |/|
		|/ |      |  \ |    \  |    | |/|
		|/ |      |  |  \__/\   \__/| |/|
		|/____________________________|/
	`)
	fmt.Println(banner)

	cmd.Execute()
}
