package cmd

// Note: these tests are relying on the global 'lintFailures' failureMap as
// defined in lint.go

import (
	"testing"
)

var (
	testLintRoot = "../testdata/linters"

	testTrackerData  []TrackerRecord
	testLintFailures int

	goodBadFiles = map[string][]string{
		// "shell":    {"shell-good.sh", "shell-bad.sh"},
		"go": {"go_good.go", "go_bad.go"},
		// "python":   {"python_good.py", "python-bad.py"},
		// "markdown": {"markdown-good.md", "markdown-bad.md"},
		// "sql":      {"sql_good.sql", "sql_bad.sql"},
	}
)

// func TestLintShell(t *testing.T) {
// 	lintShell([]string{testLintRoot + "/" + goodBadFiles["shell"][0]})
// 	if _, oops := lintFailures["lint-shell"]; oops {
// 		t.Errorf(makeTestMessage("shell", "lint-shell", "wantPass"), lintFailures)
// 	}

// 	lintShell([]string{testLintRoot + "/" + goodBadFiles["shell"][1]})
// 	if _, oops := lintFailures["lint-shell"]; !oops {
// 		t.Errorf(makeTestMessage("shell", "lint-shell", "wantFail"), lintFailures)
// 	}

// 	delete(lintFailures, "lint-shell")
// }

func TestLintGo(t *testing.T) {
	lintGo([]string{testLintRoot + "/" + goodBadFiles["go"][0]})
	testTrackerData = getTrackerData()
	testLintFailures = checkTrackerFailures(testTrackerData, "lint")
	if testLintFailures > 0 {
		t.Errorf(
			"\nlintGo failed on either the format diff-check or the lint itself, but should have succeeded -- tracker data below:\n%v",
			testTrackerData,
		)
	}
	// resets the tracker file on disk
	initTracker()

	lintGo([]string{testLintRoot + "/" + goodBadFiles["go"][1]})
	testTrackerData = getTrackerData()
	testLintFailures = checkTrackerFailures(testTrackerData, "lint")
	if testLintFailures == 0 {
		t.Errorf(
			"\nlintGo succeeded on either the format diff-check or the lint itself, but should have failed -- tracker data below:\n%v",
			testTrackerData,
		)
	}
	initTracker()
}

// func TestLintPython(t *testing.T) {
// 	lintPython([]string{testLintRoot + "/" + goodBadFiles["python"][0]})
// 	if _, oops := lintFailures["fmt-diff-check-python"]; oops {
// 		t.Errorf(makeTestMessage("python", "fmt-diff-check-python", "wantPass"), lintFailures)
// 	}
// 	if _, oops := lintFailures["typecheck-python"]; oops {
// 		t.Errorf(makeTestMessage("python", "typecheck-python", "wantPass"), lintFailures)
// 	}
// 	if _, oops := lintFailures["lint-python"]; oops {
// 		t.Errorf(makeTestMessage("python", "line-python", "wantPass"), lintFailures)
// 	}

// 	lintPython([]string{testLintRoot + "/" + goodBadFiles["python"][1]})
// 	if lintFailures["fmt-diff-check-python"] != "fail" {
// 		t.Errorf(makeTestMessage("python", "fmt-diff-check-python", "wantFail"), lintFailures)
// 	}
// 	if lintFailures["typecheck-python"] != "fail" {
// 		t.Errorf(makeTestMessage("python", "typecheck-python", "wantFail"), lintFailures)
// 	}
// 	if lintFailures["lint-python"] != "fail" {
// 		t.Errorf(makeTestMessage("python", "lint-python", "wantFail"), lintFailures)
// 	}

// 	delete(lintFailures, "fmt-diff-check-python")
// 	delete(lintFailures, "typecheck-python")
// 	delete(lintFailures, "lint-python")
// }

// func TestLintMarkdown(t *testing.T) {
// 	lintMarkdown([]string{testLintRoot + "/" + goodBadFiles["markdown"][0]})
// 	if _, oops := lintFailures["lint-markdown"]; oops {
// 		t.Errorf(makeTestMessage("markdown", "lint-markdown", "wantPass"), lintFailures)
// 	}

// 	lintMarkdown([]string{testLintRoot + "/" + goodBadFiles["markdown"][1]})
// 	if lintFailures["lint-markdown"] != "fail" {
// 		t.Errorf(makeTestMessage("markdown", "lint-markdown", "wantFail"), lintFailures)
// 	}

// 	delete(lintFailures, "lint-markdown")
// }

// func TestLintSQL(t *testing.T) {
// 	lintSQL([]string{testLintRoot + "/" + goodBadFiles["sql"][0]})
// 	if _, oops := lintFailures["lint-sql"]; oops {
// 		t.Errorf(makeTestMessage("sql", "lint-sql", "wantPass"), lintFailures)
// 	}

// 	lintSQL([]string{testLintRoot + "/" + goodBadFiles["sql"][1]})
// 	if lintFailures["lint-sql"] != "fail" {
// 		t.Errorf(makeTestMessage("sql", "lint-sql", "wantFail"), lintFailures)
// 	}

// 	delete(lintFailures, "lint-sql")
// }
