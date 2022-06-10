package cmd

import (
	"os"

	"github.com/opensourcecorp/rhad/logging"
	"github.com/spf13/cobra"
)

type failureMap map[string]string

var (
	s         syscallCfg
	syscallOk bool

	lintCmd = &cobra.Command{
		Use:   "lint",
		Short: "Run rhad's heuristic linter aggregator",
		Run:   lintExecute,
	}
	ignorePattern string

	// Used to track failures, and then throw them all as errors at the end if its size is non-zero
	lintFailures = make(failureMap)
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

	logging.Info("Running relevant linters that the GitHub Super-Linter didn't already run...")
	// lintShell(args)
	lintGo(args)
	// lintPython(args)
	// lintMarkdown(args)
	// lintSQL(args)
	// lintTerraform(args)

	if len(lintFailures) > 0 {
		logging.Error("One or more failures occurred during rhad's lint run! Summary:")
		for k, v := range lintFailures {
			logging.Error("%v: %v", k, v)
		}
		os.Exit(1)
	} else {
		logging.Info("All linters passed!")
	}
}

// func lintShell(args []string) {
// 	files := getRelevantFiles(args[0], `.*\.sh`)
// 	if len(files) > 0 {
// 		logging.Info("Running shell linter...")
// 		// Shellcheck can take multiple individual file paths in a single run
// 		var iFiles []string
// 		for _, file := range files {
// 			iFiles = append(iFiles, file.Path)
// 		}
// 		s = syscallCfg{
// 			append([]string{"shellcheck"}, iFiles...),
// 			"nonZeroExit",
// 			"",
// 		}
// 		syscallOk = syscall(s)
// 		if !syscallOk {
// 			logging.Error("Shell linter failed!")
// 			lintFailures["lint-shell"] = "fail"
// 		} else {
// 			logging.Info("Shell linter passed")
// 		}
// 	}
// }

func lintGo(args []string) {
	files := getRelevantFiles(args[0], `.*\.go`)
	if len(files) > 0 {
		logging.Info("Running Go format diff check...")
		s = syscallCfg{
			[]string{"gofmt", "-d", args[0]},
			"outputGTZero",
			"",
		}
		syscallOk = syscall(s)
		if !syscallOk {
			logging.Error("Go format diff check failed!")
			lintFailures["fmt-diff-check-go"] = "fail"
		} else {
			logging.Info("Go format diff check passed")
		}

		logging.Info("Running Go linter...")
		s = syscallCfg{
			// []string{"staticcheck", staticcheckArg},
			[]string{"golangci-lint", "run", args[0]},
			"nonZeroExit",
			"",
		}
		syscallOk = syscall(s)
		if !syscallOk {
			logging.Error("Go linter failed!")
			lintFailures["lint-go"] = "fail"
		} else {
			logging.Info("Go linter passed")
		}
	}
}

// func lintPython(args []string) {
// 	files := getRelevantFiles(args[0], `.*\.py`)
// 	if len(files) > 0 {
// 		logging.Info("Running Python format diff checker...")
// 		s = syscallCfg{
// 			[]string{"black", "--diff", args[0]},
// 			"outputGTZero",
// 			"would reformat",
// 		}
// 		syscallOk = syscall(s)
// 		if !syscallOk {
// 			logging.Error("Python format diff check failed!")
// 			lintFailures["fmt-diff-check-python"] = "fail"
// 		} else {
// 			logging.Info("Python format diff checker passed")
// 		}

// 		logging.Info("Running Python typecheck...")
// 		s = syscallCfg{
// 			[]string{"mypy", "--no-incremental", args[0]},
// 			"nonZeroExit",
// 			"would reformat",
// 		}
// 		syscallOk = syscall(s)
// 		if !syscallOk {
// 			logging.Error("Python type checker failed!")
// 			lintFailures["typecheck-python"] = "fail"
// 		} else {
// 			logging.Info("Python typecheck passed")
// 		}

// 		logging.Info("Running Python linter...")
// 		s = syscallCfg{
// 			[]string{"pylint", "--recursive=y", "--disable=import-error,invalid-name", args[0]},
// 			"nonZeroExit",
// 			"",
// 		}
// 		syscallOk = syscall(s)
// 		if !syscallOk {
// 			logging.Error("Python linter failed!")
// 			lintFailures["lint-python"] = "fail"
// 		} else {
// 			logging.Info("Python linter passed")
// 		}
// 	}
// }

// func lintMarkdown(args []string) {
// 	files := getRelevantFiles(args[0], `.*\.(md|markdown)`)
// 	if len(files) > 0 {
// 		logging.Info("Running Markdown linter...")
// 		s = syscallCfg{
// 			[]string{"mdl", "--style", rhadSrc + "/.mdlrc.style.rb", args[0]},
// 			"nonZeroExit",
// 			"",
// 		}
// 		syscallOk = syscall(s)
// 		if !syscallOk {
// 			logging.Error("Markdown linter failed!")
// 			lintFailures["lint-markdown"] = "fail"
// 		} else {
// 			logging.Info("Markdown linter passed")
// 		}
// 	}
// }

// func lintSQL(args []string) {
// 	files := getRelevantFiles(args[0], `.*\.sql`)
// 	if len(files) > 0 {
// 		logging.Info("Running SQL linter...")
// 		s = syscallCfg{
// 			[]string{"sqlfluff", "lint", "--dialect", "postgres", args[0]},
// 			"nonZeroExit",
// 			"",
// 		}
// 		syscallOk = syscall(s)
// 		if !syscallOk {
// 			logging.Error("SQL linter failed!")
// 			lintFailures["lint-sql"] = "fail"
// 		} else {
// 			logging.Info("SQL linter passed")
// 		}
// 	}
// }

// TODO: Find a Terraform linter that doesn't suck
// func lintTerraform(args []string) {
// 	files := getRelevantFiles(`.*\.tf(vars)?`)
// 	if len(files) > 0 {
// 		logging.Info("Running Terraform linter...")
// 		s = syscallCfg{
// 			[]string{},
// 			"nonZeroExit",
// 			"",
// 		}
// 		syscallOk = syscall(s)
// 		if !syscallOk {
// 			logging.Error("Terraform linter failed!")
// 			lintFailures["lint-terraform"] = "fail"
// 		} else {
// 			logging.Info("Terraform linter passed")
// 		}
// 	}
// }
