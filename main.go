package main

import (
	"fmt"

	osc "github.com/opensourcecorp/go-common"
	"github.com/opensourcecorp/oscar/cmd"
)

func main() {
	osc.SetLoggerPrefixName("oscar")

	osc.InfoLog("Firing up oscar!")

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
