package main

import (
	"fmt"

	osc "github.com/opensourcecorp/go-common"
	"github.com/opensourcecorp/rhad/cmd"
)

func main() {
	osc.SetLoggerPrefixName("rhad")

	osc.InfoLog("Firing up rhad!")

	banner := `--------------------------
  __   |      __      __ |
 /  \  |--_  /  \    /  \|
|      |  | |    \  |    |
|      |  |  \__/\   \__/|
--------------------------`
	fmt.Println(banner)

	cmd.Execute()
}
