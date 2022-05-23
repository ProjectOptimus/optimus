package run

import (
	"os"
	"os/exec"
	"regexp"

	"github.com/opensourcecorp/rhad/logging"
)

var (
	// Common vars that I don't want to have to keep redefining
	cmd    *exec.Cmd
	output []byte
	err    error

	// Used to track failures, and then throw them all as errors at the end if its size is non-zero
	failures = make(failureMap)
)

type failureMap map[string]string

func lintShell(cfg rhadConfig) {
	files := getRelevantFiles(".*\\.sh", cfg)
	if len(files) > 0 {
		logging.Info("Running shell linter...")
		var iFiles []string
		for _, file := range files {
			iFiles = append(iFiles, file.path)
		}

		cmd = exec.Command("shellcheck", iFiles...)
		output, err = cmd.CombinedOutput()
		if err != nil {
			// TODO: we don't want to exit here since a failure could be the
			// shell return code -- but, it could be something worse, sure
			logging.Error("Output below:\n" + string(output))
			logging.Error(err.Error())
			logging.Error("Shell linter failed!")
			failures["lint-shell"] = "fail"
		} else {
			logging.Info("Shell linter passed")
		}
	}
}

func lintGo(cfg rhadConfig) {
	files := getRelevantFiles(".*\\.go", cfg)
	if len(files) > 0 {
		logging.Info("Running Go format diff check...")
		cmd = exec.Command("gofmt", "-d", *cfg.cliOpts.path)
		output, err = cmd.CombinedOutput()
		if err != nil {
			logging.Error("Output below:\n" + string(output))
			logging.Error(err.Error())
		}
		if len(output) > 0 {
			logging.Error("Output below:\n" + string(output))
			logging.Error("Go format diff check failed!")
			failures["fmt-diff-check-go"] = "fail"
		} else {
			logging.Error("Go format diff check passed")
		}

		logging.Info("Running Go linter...")
		// staticcheck is such a pain in my ass
		f, err := os.Stat(*cfg.cliOpts.path)
		if err != nil {
			logging.Error("I have no idea what kind of '-path' you just tried to give me")
			logging.Error(err.Error())
		}
		if f.IsDir() {
			cmd = exec.Command("staticcheck", "./...")
			cmd.Dir = *cfg.cliOpts.path
		} else {
			cmd = exec.Command("staticcheck", *cfg.cliOpts.path)
		}
		output, err = cmd.CombinedOutput()
		if err != nil {
			logging.Error("Output below:\n" + string(output))
			logging.Error(err.Error())
			logging.Error("Go linter failed!")
			failures["lint-go"] = "fail"
		} else {
			logging.Error("Go linter passed")
		}
	}
}

func lintPython(cfg rhadConfig) {
	files := getRelevantFiles(".*\\.py", cfg)
	if len(files) > 0 {
		logging.Info("Running Python format diff checker...")
		cmd = exec.Command("black", "--diff", *cfg.cliOpts.path)
		output, err = cmd.CombinedOutput()
		if err != nil {
			logging.Error("Output below:\n" + string(output))
			logging.Error(err.Error())
		}
		if len(output) > 0 {
			regex := regexp.MustCompile("would reformat")
			if regex.MatchString(string(output)) {
				logging.Error("Output below:\n" + string(output))
				logging.Error("Python format diff check failed!")
				failures["fmt-diff-check-python"] = "fail"
			}
		} else {
			logging.Error("Python format diff checker passed")
		}

		logging.Info("Running Python typecheck...")
		cmd = exec.Command("mypy", *cfg.cliOpts.path)
		output, err = cmd.CombinedOutput()
		if err != nil {
			logging.Error("Output below:\n" + string(output))
			logging.Error(err.Error())
			logging.Error("Python type checker failed!")
			failures["typecheck-python"] = "fail"
		} else {
			logging.Error("Python typecheck passed")
		}

		logging.Info("Running Python linter...")
		cmd = exec.Command("pylint", "--recursive=y", "--disable=import-error,invalid-name", *cfg.cliOpts.path)
		output, err = cmd.CombinedOutput()
		if err != nil {
			logging.Error("Output below:\n" + string(output))
			logging.Error(err.Error())
			logging.Error("Python linter failed!")
			failures["lint-python"] = "fail"
		} else {
			logging.Error("Python linter passed")
		}
	}
}

func lintMarkdown(cfg rhadConfig) {
	files := getRelevantFiles(".*\\.(md|markdown)", cfg)
	if len(files) > 0 {
		logging.Info("Running Markdown linter...")
		cmd = exec.Command("mdl", "--style", "../.mdlrc.style.rb", *cfg.cliOpts.path)
		output, err = cmd.CombinedOutput()
		if err != nil {
			logging.Error("Output below:\n" + string(output))
			logging.Error(err.Error())
			logging.Error("Markdown linter failed!")
			failures["lint-markdown"] = "fail"
		} else {
			logging.Error("Markdown linter passed")
		}
	}
}

func Lint(cfg rhadConfig) {
	logging.Info("Running relevant linters...")
	lintShell(cfg)
	lintGo(cfg)
	lintPython(cfg)
	lintMarkdown(cfg)

	if len(failures) > 0 {
		logging.Error("One or more failures occurred during rhad's lint run! Summary:")
		logging.Error("%v", failures)
		os.Exit(1)
	} else {
		logging.Info("All linters passed!")
	}
}
