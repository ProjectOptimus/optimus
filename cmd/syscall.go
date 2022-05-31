package cmd

import (
	"os"
	"os/exec"
	"regexp"

	"github.com/opensourcecorp/rhad/logging"
)

type syscallCfg struct {
	cmdLine                 []string
	errCheckType            string
	outputErrorPatternMatch string
}

func syscall(s syscallCfg) bool {
	var cmd *exec.Cmd
	if len(s.cmdLine) == 1 {
		cmd = exec.Command(s.cmdLine[0])
	} else if len(s.cmdLine) > 1 {
		cmd = exec.Command(s.cmdLine[0], s.cmdLine[1:]...)
	} else {
		logging.Error("how tf u gonna give me a zero-length command")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		logging.Error("Output below:\n" + string(output))
		logging.Error(err.Error())
		return false
	}

	switch s.errCheckType {
	case "nonZeroExit":
		// this should have failed early above, so we can just return true early
		// here
		return true
	case "outputGTZero":
		if len(output) > 0 {
			if s.outputErrorPatternMatch == "" {
				logging.Error("Output below:\n" + string(output))
				return false
			} else {
				regex := regexp.MustCompile(s.outputErrorPatternMatch)
				if regex.MatchString(string(output)) {
					logging.Error("Output below:\n" + string(output))
					return false
				}
			}
		}
	default:
		// If it was a nonzero exit syscall, they should never get here anyway
		logging.Error("Unhandled syscall() errCheckType")
	}

	return true
}

// This can be used to run before functions that are making syscalls, so they
// hopefully catch runtime errors earlier
func testSysinit() {
	s := syscallCfg{
		[]string{"bash", rhadSrc + "/scripts/sysinit.sh", "test"},
		"nonZeroExit",
		"",
	}
	syscallOk := syscall(s)
	if !syscallOk {
		os.Exit(1)
	}
}
