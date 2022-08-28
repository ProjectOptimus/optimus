// Package cmd's lint submodule provides linting functionality that is either
// easier to configure or has unavailable behavior in external linter
// aggregators
package cmd

import (
	"os"

	osc "github.com/opensourcecorp/go-common"
	"github.com/spf13/cobra"
)

var (
	sc osc.Syscall

	lintCmd = &cobra.Command{
		Use:   "lint",
		Short: "Run rhad's heuristic linter aggregator",
		Run:   lintExecute,
	}
	ignorePattern string
)

func init() {
	rootCmd.AddCommand(lintCmd)
	// TODO: this works, but isn't being implemented because of how the linters
	// operate -- they target directories, not individual files, and so even
	// though the files are ignored, the linter calls hit the whole tree. So,
	// find a way to get the linters to do better. Maybe copy the tree to /tmp
	// excluding the ignored files?
	lintCmd.PersistentFlags().StringVarP(&ignorePattern, "ignore-pattern", "i", `^\b$`, "(NOTE: NOT CURRENTLY WORKING) Valid regex pattern of paths to ignore") // default can never be matched in a regex
}

func lintExecute(cmd *cobra.Command, args []string) {
	testSysinit()

	if len(args) == 0 {
		args = []string{"."}
	}

	osc.InfoLog("Running relevant linters that the GitHub Super-Linter didn't already run...")
	lintGo(args)

	trackerData := getTrackerData()
	failures := checkTrackerFailures(trackerData, "lint")

	if failures > 0 {
		osc.ErrorLog(nil, "One or more failures occurred during rhad's lint run! Summary:")
		for _, record := range trackerData {
			if record.Result == "fail" {
				osc.ErrorLog(nil, "%s", record)
			}
		}
		os.Exit(1)
	} else {
		osc.InfoLog("All linters passed!")
	}
}

func lintGo(args []string) {
	files := getRelevantFiles(args[0], `.*\.go`)
	if len(files) > 0 {
		osc.InfoLog("Running Go format diff check...")
		sc = osc.Syscall{
			CmdLine:      []string{"gofmt", "-d", args[0]},
			ErrCheckType: "outputGTZero",
		}
		sc.Exec()
		if !sc.Ok {
			writeTrackerRecord(TrackerRecord{
				Type:     "lint",
				Subtype:  "go-fmt-diff-check",
				Language: "go",
				Tool:     "gofmt",
				Result:   "fail",
			})
			osc.ErrorLog(nil, "Go format diff check failed!")
		} else {
			writeTrackerRecord(TrackerRecord{
				Type:     "lint",
				Subtype:  "fmt-diff-check",
				Language: "go",
				Tool:     "go fmt",
				Result:   "pass",
			})
			osc.InfoLog("Go format diff check passed")
		}

		osc.InfoLog("Running Go linter...")
		sc = osc.Syscall{
			CmdLine: []string{"golangci-lint", "run", args[0]},
		}
		sc.Exec()
		if !sc.Ok {
			writeTrackerRecord(TrackerRecord{
				Type:     "lint",
				Language: "go",
				Tool:     "golangci-lint",
				Result:   "fail",
			})
			osc.ErrorLog(nil, "Go linter failed!")
		} else {
			writeTrackerRecord(TrackerRecord{
				Type:     "lint",
				Language: "go",
				Tool:     "golangci-lint",
				Result:   "pass",
			})
			osc.InfoLog("Go linter passed")
		}
	}
}
