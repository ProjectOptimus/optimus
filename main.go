package main

import (
	"fmt"

	osc "github.com/opensourcecorp/go-common"
	"github.com/opensourcecorp/oscar/cmd"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})
}

func main() {
	banner := osc.Dedent(`
		   ____________________
		 /____________________/|
		|  _   _   _  _   _  |/|
		| | | |_  |  |_| |_| |/|
		| |_|  _| |_ | | | \ |/|
		|____________________|/
	`)
	fmt.Println(banner)

	cmd.Execute()
}
