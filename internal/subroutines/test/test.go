package test

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "test",
		Short: "Evaluate & run test suites across codebase",
		Run:   testExecute,
	}
)

func testExecute(_ *cobra.Command, _ []string) {
	logrus.Fatal("test: Not implemented")
}
