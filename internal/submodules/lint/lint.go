// Package lint provides code linting functionality that is either easier to
// configure or has unavailable behavior in external linter aggregators
package lint

import (
	"os"

	osc "github.com/opensourcecorp/go-common"
	"github.com/opensourcecorp/rhad/internal/fsutils"
	"github.com/opensourcecorp/rhad/internal/tracker"
	"github.com/opensourcecorp/rhad/internal/utils"
	"github.com/spf13/cobra"
)

var (
	sc osc.Syscall

	LintCmd = &cobra.Command{
		Use:   "lint",
		Short: "Run rhad's heuristic linter aggregator",
		Run:   lintExecute,
	}
)

func init() {
	// // TODO: this works, but isn't being implemented because of how the linters
	// // operate -- they target directories, not individual files, and so even
	// // though the files are ignored, the linter calls hit the whole tree. So,
	// // find a way to get the linters to do better. Maybe copy the tree to /tmp
	// // excluding the ignored files?
	// lintCmd.PersistentFlags().StringVarP(&ignorePattern, "ignore-pattern", "i", `^\b$`, "(NOTE: NOT CURRENTLY WORKING) Valid regex pattern of paths to ignore") // default can never be matched in a regex
}

func lintExecute(cmd *cobra.Command, args []string) {
	utils.VerifySysinit()

	// If no args provided, we'll make an assumption that most callables will
	// take the current directory as an arg. We can control this override within
	// each linter function call, though
	if len(args) == 0 {
		args = []string{"."}
	}

	osc.InfoLog("Running relevant linters that the GitHub Super-Linter didn't already run...")
	lintGo(args)

	trackerData := tracker.GetTrackerData()
	failures := tracker.CheckTrackerFailures(trackerData, "lint")

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
	files := fsutils.GetRelevantFiles(args[0], `.*\.go`)
	if len(files) > 0 {
		// If more than a single file, might as well just use go's package tree
		// syntax to look for packages if the tool supports it
		var packageTree string
		if len(files) > 1 {
			packageTree = "./..."
		} else {
			// But linter unit test runs usually are just given a single file, so fall back
			packageTree = args[0]
		}

		// Vetter
		osc.InfoLog("Running Go vetter...")
		sc = osc.Syscall{
			CmdLine: []string{"go", "vet", packageTree},
		}
		sc.Exec()
		if !sc.Ok {
			tracker.WriteTrackerRecord(tracker.TrackerRecord{
				Type:     "lint",
				Subtype:  "vet",
				Language: "go",
				Tool:     "go vet",
				Result:   "fail",
			})
			osc.ErrorLog(nil, "Go vetter failed!")
		} else {
			tracker.WriteTrackerRecord(tracker.TrackerRecord{
				Type:     "lint",
				Subtype:  "vet",
				Language: "go",
				Tool:     "go vet",
				Result:   "pass",
			})
			osc.InfoLog("Go vetter passed")
		}

		// Format diff checker
		osc.InfoLog("Running Go format diff check...")
		sc = osc.Syscall{
			CmdLine:      []string{"gofmt", "-d", args[0]}, // gofmt doesn't support the package tree syntax
			ErrCheckType: "outputGTZero",
		}
		sc.Exec()
		if !sc.Ok {
			tracker.WriteTrackerRecord(tracker.TrackerRecord{
				Type:     "lint",
				Subtype:  "go-fmt-diff-check",
				Language: "go",
				Tool:     "gofmt",
				Result:   "fail",
			})
			osc.ErrorLog(nil, "Go format diff check failed!")
		} else {
			tracker.WriteTrackerRecord(tracker.TrackerRecord{
				Type:     "lint",
				Subtype:  "fmt-diff-check",
				Language: "go",
				Tool:     "gofmt",
				Result:   "pass",
			})
			osc.InfoLog("Go format diff check passed")
		}

		// Linter
		osc.InfoLog("Running Go linter...")
		sc = osc.Syscall{
			CmdLine: []string{"golangci-lint", "run", packageTree},
		}
		sc.Exec()
		if !sc.Ok {
			tracker.WriteTrackerRecord(tracker.TrackerRecord{
				Type:     "lint",
				Language: "go",
				Tool:     "golangci-lint",
				Result:   "fail",
			})
			osc.ErrorLog(nil, "Go linter failed!")
		} else {
			tracker.WriteTrackerRecord(tracker.TrackerRecord{
				Type:     "lint",
				Language: "go",
				Tool:     "golangci-lint",
				Result:   "pass",
			})
			osc.InfoLog("Go linter passed")
		}
	}
}
