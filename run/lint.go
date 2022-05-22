package run

import (
	"fmt"
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
		if !cfg.cliOpts.testing {
			logging.Info("Running Shellcheck...")
		}
		for _, file := range files {
			cmd = exec.Command("shellcheck", file.path)
			output, err = cmd.CombinedOutput()
			if err != nil {
				// TODO: we don't want to exit here since a failure could be the
				// shell return code -- but, it could be something worse, sure
				if !cfg.cliOpts.testing {
					fmt.Println(string(output))
					fmt.Println(err)
					logging.Error("Shellcheck failed!")
				}
				failures["lint-shell"] = "fail"
			}
		}
	}
}

func lintGo(cfg rhadConfig) {
	files := getRelevantFiles(".*\\.go", cfg)
	if len(files) > 0 {
		if !cfg.cliOpts.testing {
			logging.Info("Running Go format diff check...")
		}
		cmd = exec.Command("gofmt", "-d", *cfg.cliOpts.path)
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(output))
			fmt.Println(err)
		}
		if len(output) > 0 {
			if !cfg.cliOpts.testing {
				fmt.Println(string(output))
				logging.Error("Go format diff check failed!")
			}
			failures["fmt-diff-check-go"] = "fail"
		}

		if !cfg.cliOpts.testing {
			logging.Info("Running Go linter...")
		}

		// staticcheck is such a pain in my ass
		f, err := os.Stat(*cfg.cliOpts.path)
		if err != nil {
			logging.Error("I have no idea what kind of '-path' you just tried to give me")
			fmt.Println(err)
		}
		if f.IsDir() {
			cmd = exec.Command("staticcheck", "./...")
			cmd.Dir = *cfg.cliOpts.path
		} else {
			cmd = exec.Command("staticcheck", *cfg.cliOpts.path)
		}
		output, err = cmd.CombinedOutput()
		if err != nil {
			if !cfg.cliOpts.testing {
				fmt.Println(string(output))
				fmt.Println(err)
				logging.Error("Go linter failed!")
			}
			failures["lint-go"] = "fail"
		}
	}
}

func lintPython(cfg rhadConfig) {
	files := getRelevantFiles(".*\\.py", cfg)
	if len(files) > 0 {
		if !cfg.cliOpts.testing {
			logging.Info("Running Python format diff checker...")
		}
		cmd = exec.Command("black", "--diff", *cfg.cliOpts.path)
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(output))
			fmt.Println(err)
		}
		if len(output) > 0 {
			regex := regexp.MustCompile("would reformat")
			if regex.MatchString(string(output)) {
				if !cfg.cliOpts.testing {
					fmt.Println(string(output))
					logging.Error("Python format diff check failed!")
				}
				failures["fmt-diff-check-python"] = "fail"
			}
		}

		if !cfg.cliOpts.testing {
			logging.Info("Running Python typechecker...")
		}
		cmd = exec.Command("mypy", *cfg.cliOpts.path)
		output, err = cmd.CombinedOutput()
		if err != nil {
			if !cfg.cliOpts.testing {
				fmt.Println(string(output))
				fmt.Println(err)
				logging.Error("Python type checker failed!")
			}
			failures["typecheck-python"] = "fail"
		}

		if !cfg.cliOpts.testing {
			logging.Info("Running Python linter...")
		}
		for _, file := range files {
			cmd = exec.Command("pylint", "--disable=import-error,invalid-name", file.path)
			output, err = cmd.CombinedOutput()
			if err != nil {
				if !cfg.cliOpts.testing {
					fmt.Println(string(output))
					fmt.Println(err)
					logging.Error("Python linter failed!")
				}
				failures["lint-python"] = "fail"
			}
		}
	}
}

func lintMarkdown(cfg rhadConfig) {
	files := getRelevantFiles(".*\\.(md|markdown)", cfg)
	if len(files) > 0 {
		if !cfg.cliOpts.testing {
			logging.Info("Running Markdown linter...")
		}
		for _, file := range files {
			cmd = exec.Command("mdl", "--style", "../.mdlrc.style.rb", file.path)
			output, err = cmd.CombinedOutput()
			if err != nil {
				if !cfg.cliOpts.testing {
					fmt.Println(string(output))
					fmt.Println(err)
					logging.Error("Markdown linter failed!")
				}
				failures["lint-markdown"] = "fail"
			}
		}
	}
}

func Lint(cfg rhadConfig) {
	if !cfg.cliOpts.testing {
		logging.Info("Running relevant linters...")
	}
	lintShell(cfg)
	lintGo(cfg)
	lintPython(cfg)
	lintMarkdown(cfg)

	if len(failures) > 0 {
		if !cfg.cliOpts.testing {
			logging.Error("One or more failures occurred during rhad's lint run! Summary:")
			fmt.Println(failures)
		}
		os.Exit(1)
	} else {
		if !cfg.cliOpts.testing {
			logging.Info("All linters passed!")
		}
	}
}
