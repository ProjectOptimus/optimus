package run

// Note: these tests are relying on the global 'failures' failureMap as defined
// in lint.go

import (
	"fmt"
	"log"
	"testing"
)

var (
	rootPath = "../testdata/linters"
	testPath string

	cfg = rhadConfig{
		map[string]any{"": ""},
		cliOptsStruct{
			path: &rootPath,
		},
	}

	goodBadFiles = map[string][]string{
		"shell":    {"shell-good.sh", "shell-bad.sh"},
		"go":       {"go_good.go", "go_bad.go"},
		"python":   {"python_good.py", "python-bad.py"},
		"markdown": {"markdown-good.md", "markdown-bad.md"},
	}
)

func makeTestMessage(lintType, mapKeyName, wantResult string) string {
	var msg string
	switch wantResult {
	case "pass":
		msg = fmt.Sprintf("failureMap 'failures' does not reflect a(n) %s linting success as expected - map key '%s' should not exist but might: %%v", lintType, mapKeyName)
	case "fail":
		msg = fmt.Sprintf("failureMap 'failures' does not reflect a(n) %s linting failure as expected - map key should '%s' have value 'fail': %%v", lintType, mapKeyName)
	default:
		log.Fatalln("hey you're not supposed to get here")
	}

	return msg
}

func TestLintShell(t *testing.T) {
	testPath = rootPath + "/" + goodBadFiles["shell"][0]
	cfg.cliOpts.path = &testPath
	lintShell(cfg)
	if _, oops := failures["lint-shell"]; oops {
		t.Errorf(makeTestMessage("shell", "lint-shell", "pass"), failures)
	}

	testPath = rootPath + "/" + goodBadFiles["shell"][1]
	cfg.cliOpts.path = &testPath
	lintShell(cfg)
	if failures["lint-shell"] != "fail" {
		t.Errorf(makeTestMessage("shell", "lint-shell", "fail"), failures)
	}

	delete(failures, "lint-shell")
}

func TestLintGo(t *testing.T) {
	testPath = rootPath + "/" + goodBadFiles["go"][0]
	cfg.cliOpts.path = &testPath
	lintGo(cfg)
	if _, oops := failures["fmt-diff-check-go"]; oops {
		t.Errorf(makeTestMessage("go", "fmt-diff-check-go", "pass"), failures)
	}
	if _, oops := failures["lint-go"]; oops {
		t.Errorf(makeTestMessage("go", "lint-go", "pass"), failures)
	}

	testPath = rootPath + "/" + goodBadFiles["go"][1]
	cfg.cliOpts.path = &testPath
	lintGo(cfg)
	if failures["fmt-diff-check-go"] != "fail" {
		t.Errorf(makeTestMessage("go", "fmt-diff-check-go", "fail"), failures)
	}
	if failures["lint-go"] != "fail" {
		t.Errorf(makeTestMessage("go", "lint-go", "fail"), failures)
	}

	delete(failures, "fmt-diff-check-go")
	delete(failures, "lint-go")
}

func TestLintPython(t *testing.T) {
	testPath = rootPath + "/" + goodBadFiles["python"][0]
	cfg.cliOpts.path = &testPath
	lintPython(cfg)
	if _, oops := failures["fmt-diff-check-python"]; oops {
		t.Errorf(makeTestMessage("python", "fmt-diff-check-python", "pass"), failures)
	}
	if _, oops := failures["typecheck-python"]; oops {
		t.Errorf(makeTestMessage("python", "typecheck-python", "pass"), failures)
	}
	if _, oops := failures["lint-python"]; oops {
		t.Errorf(makeTestMessage("python", "line-python", "pass"), failures)
	}

	testPath = rootPath + "/" + goodBadFiles["python"][1]
	cfg.cliOpts.path = &testPath
	lintPython(cfg)
	if failures["fmt-diff-check-python"] != "fail" {
		t.Errorf(makeTestMessage("python", "fmt-diff-check-python", "fail"), failures)
	}
	if failures["typecheck-python"] != "fail" {
		t.Errorf(makeTestMessage("python", "typecheck-python", "fail"), failures)
	}
	if failures["lint-python"] != "fail" {
		t.Errorf(makeTestMessage("python", "lint-python", "fail"), failures)
	}

	delete(failures, "fmt-diff-check-python")
	delete(failures, "typecheck-python")
	delete(failures, "lint-python")
}

func TestLintMarkdown(t *testing.T) {
	testPath = rootPath + "/" + goodBadFiles["markdown"][0]
	cfg.cliOpts.path = &testPath
	lintMarkdown(cfg)
	if _, oops := failures["lint-markdown"]; oops {
		t.Errorf(makeTestMessage("markdown", "lint-markdown", "pass"), failures)
	}

	testPath = rootPath + "/" + goodBadFiles["markdown"][1]
	cfg.cliOpts.path = &testPath
	lintMarkdown(cfg)
	if failures["lint-markdown"] != "fail" {
		t.Errorf(makeTestMessage("markdown", "lint-markdown", "fail"), failures)
	}

	delete(failures, "lint-markdown")
}
