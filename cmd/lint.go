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

	// Used to track failures, and then throw them all as errors at the end if its size is non-zero
	lintFailures = make(failureMap)
)

func init() {
	rootCmd.AddCommand(lintCmd)
}

func lintExecute(cmd *cobra.Command, args []string) {
	testSysinit()

	if len(args) == 0 {
		args = []string{"."}
	}

	logging.Info("Running relevant linters...")
	lintShell(args)
	lintGo(args)
	lintPython(args)
	lintMarkdown(args)

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

func lintShell(args []string) {
	files := getRelevantFiles(args[0], ".*\\.sh")
	if len(files) > 0 {
		logging.Info("Running shell linter...")
		// Shellcheck can take multiple individual file paths in a single run
		var iFiles []string
		for _, file := range files {
			iFiles = append(iFiles, file.path)
		}
		s = syscallCfg{
			append([]string{"shellcheck"}, iFiles...),
			"nonZeroExit",
			"",
		}
		syscallOk = syscall(s)
		if !syscallOk {
			logging.Error("Shell linter failed!")
			lintFailures["lint-shell"] = "fail"
		} else {
			logging.Info("Shell linter passed")
		}
	}
}

func lintGo(args []string) {
	files := getRelevantFiles(args[0], ".*\\.go")
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
			logging.Error("Go format diff check passed")
		}

		logging.Info("Running Go linter...")
		// I hate staticcheck
		f, err := os.Lstat(args[0])
		if err != nil {
			logging.Error("Error looking at file provided to lintGo()")
			logging.Error(err.Error())
			os.Exit(1)
		}
		var staticcheckArg string
		if f.Mode().IsRegular() {
			staticcheckArg = args[0]
		} else {
			staticcheckArg = "./..."
		}

		s = syscallCfg{
			[]string{"staticcheck", staticcheckArg},
			"nonZeroExit",
			"",
		}
		syscallOk = syscall(s)
		if !syscallOk {
			logging.Error("Go linter failed!")
			lintFailures["lint-go"] = "fail"
		} else {
			logging.Error("Go linter passed")
		}
	}
}

func lintPython(args []string) {
	files := getRelevantFiles(args[0], ".*\\.py")
	if len(files) > 0 {
		logging.Info("Running Python format diff checker...")
		s = syscallCfg{
			[]string{"black", "--diff", args[0]},
			"outputGTZero",
			"would reformat",
		}
		syscallOk = syscall(s)
		if !syscallOk {
			logging.Error("Python format diff check failed!")
			lintFailures["fmt-diff-check-python"] = "fail"
		} else {
			logging.Error("Python format diff checker passed")
		}

		logging.Info("Running Python typecheck...")
		s = syscallCfg{
			[]string{"mypy", "--no-incremental", args[0]},
			"nonZeroExit",
			"would reformat",
		}
		syscallOk = syscall(s)
		if !syscallOk {
			logging.Error("Python type checker failed!")
			lintFailures["typecheck-python"] = "fail"
		} else {
			logging.Error("Python typecheck passed")
		}

		logging.Info("Running Python linter...")
		s = syscallCfg{
			[]string{"pylint", "--recursive=y", "--disable=import-error,invalid-name", args[0]},
			"nonZeroExit",
			"",
		}
		syscallOk = syscall(s)
		if !syscallOk {
			logging.Error("Python linter failed!")
			lintFailures["lint-python"] = "fail"
		} else {
			logging.Error("Python linter passed")
		}
	}
}

func lintMarkdown(args []string) {
	files := getRelevantFiles(args[0], ".*\\.(md|markdown)")
	if len(files) > 0 {
		logging.Info("Running Markdown linter...")
		s = syscallCfg{
			[]string{"mdl", "--style", rhadSrc + "/.mdlrc.style.rb", args[0]},
			"nonZeroExit",
			"",
		}
		syscallOk = syscall(s)
		if !syscallOk {
			logging.Error("Markdown linter failed!")
			lintFailures["lint-markdown"] = "fail"
		} else {
			logging.Error("Markdown linter passed")
		}
	}
}
