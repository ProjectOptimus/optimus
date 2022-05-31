package cmd

// Note: these tests are relying on the global 'lintFailures' failureMap as
// defined in lint.go

import (
	"fmt"
	"log"
	"testing"
)

var (
	testLintRoot = "../testdata/linters"

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
	case "wantPass":
		msg = fmt.Sprintf("failureMap 'lintFailures' does not reflect a(n) %s linting success as expected - map key '%s' should not exist but might: %%v", lintType, mapKeyName)
	case "wantFail":
		msg = fmt.Sprintf("failureMap 'lintFailures' does not reflect a(n) %s linting failure as expected - map key should '%s' have value 'fail': %%v", lintType, mapKeyName)
	default:
		log.Fatalln("hey you're not supposed to get here")
	}

	return msg
}

func TestLintShell(t *testing.T) {
	lintShell([]string{testLintRoot + "/" + goodBadFiles["shell"][0]})
	if _, oops := lintFailures["lint-shell"]; oops {
		t.Errorf(makeTestMessage("shell", "lint-shell", "wantPass"), lintFailures)
	}

	lintShell([]string{testLintRoot + "/" + goodBadFiles["shell"][1]})
	if _, oops := lintFailures["lint-shell"]; !oops {
		t.Errorf(makeTestMessage("shell", "lint-shell", "wantFail"), lintFailures)
	}

	delete(lintFailures, "lint-shell")
}

func TestLintGo(t *testing.T) {
	lintGo([]string{testLintRoot + "/" + goodBadFiles["go"][0]})
	if _, oops := lintFailures["fmt-diff-check-go"]; oops {
		t.Errorf(makeTestMessage("go", "fmt-diff-check-go", "wantPass"), lintFailures)
	}
	if _, oops := lintFailures["lint-go"]; oops {
		t.Errorf(makeTestMessage("go", "lint-go", "wantPass"), lintFailures)
	}

	lintGo([]string{testLintRoot + "/" + goodBadFiles["go"][1]})
	if lintFailures["fmt-diff-check-go"] != "fail" {
		t.Errorf(makeTestMessage("go", "fmt-diff-check-go", "wantFail"), lintFailures)
	}
	if lintFailures["lint-go"] != "fail" {
		t.Errorf(makeTestMessage("go", "lint-go", "wantFail"), lintFailures)
	}

	delete(lintFailures, "fmt-diff-check-go")
	delete(lintFailures, "lint-go")
}

func TestLintPython(t *testing.T) {
	lintPython([]string{testLintRoot + "/" + goodBadFiles["python"][0]})
	if _, oops := lintFailures["fmt-diff-check-python"]; oops {
		t.Errorf(makeTestMessage("python", "fmt-diff-check-python", "wantPass"), lintFailures)
	}
	if _, oops := lintFailures["typecheck-python"]; oops {
		t.Errorf(makeTestMessage("python", "typecheck-python", "wantPass"), lintFailures)
	}
	if _, oops := lintFailures["lint-python"]; oops {
		t.Errorf(makeTestMessage("python", "line-python", "wantPass"), lintFailures)
	}

	lintPython([]string{testLintRoot + "/" + goodBadFiles["python"][1]})
	if lintFailures["fmt-diff-check-python"] != "fail" {
		t.Errorf(makeTestMessage("python", "fmt-diff-check-python", "wantFail"), lintFailures)
	}
	if lintFailures["typecheck-python"] != "fail" {
		t.Errorf(makeTestMessage("python", "typecheck-python", "wantFail"), lintFailures)
	}
	if lintFailures["lint-python"] != "fail" {
		t.Errorf(makeTestMessage("python", "lint-python", "wantFail"), lintFailures)
	}

	delete(lintFailures, "fmt-diff-check-python")
	delete(lintFailures, "typecheck-python")
	delete(lintFailures, "lint-python")
}

func TestLintMarkdown(t *testing.T) {
	lintMarkdown([]string{testLintRoot + "/" + goodBadFiles["markdown"][0]})
	if _, oops := lintFailures["lint-markdown"]; oops {
		t.Errorf(makeTestMessage("markdown", "lint-markdown", "wantPass"), lintFailures)
	}

	lintMarkdown([]string{testLintRoot + "/" + goodBadFiles["markdown"][1]})
	if lintFailures["lint-markdown"] != "fail" {
		t.Errorf(makeTestMessage("markdown", "lint-markdown", "wantFail"), lintFailures)
	}

	delete(lintFailures, "lint-markdown")
}

func TestLintSQL(t *testing.T) {
	t.Errorf("Need to implement SQL lint test")
}
