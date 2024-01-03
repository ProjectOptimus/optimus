package osc

import (
	"os/exec"
	"regexp"
	"strings"
)

// Syscall contains data necessary to build, run, and validate external calls.
// Callers can leverage the default ("zeroed") fields of this struct -- e.g. for
// a simple "run this and see if it fails" call, the struct literal need only
// have CmdLine defined.
type Syscall struct {
	CmdLine                 []string // The actual command and its args to pass to Exec()
	ErrCheckType            string   // What kind of logical check to make against Exec() returns
	OutputErrorPatternMatch string   // Used to match against Exec() output to throw an error
	Ok                      bool     // Was Exec() successful based on conditions?
	Output                  string   // Used to surface Exec() combined stdout/stderr to the caller, in case it needs to look at it
}

// Exec invokes the command & arguments as specified by Syscall.CmdLine. It then
// collects some useful information for the caller to check or inspect against.
func (sc *Syscall) Exec() {
	var cmd *exec.Cmd
	if len(sc.CmdLine) == 1 {
		cmd = exec.Command(sc.CmdLine[0])
	} else if len(sc.CmdLine) > 1 {
		cmd = exec.Command(sc.CmdLine[0], sc.CmdLine[1:]...)
	} else {
		FatalLog(nil, "how tf u gonna give me a zero-length command")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		ErrorLog(err, "Exec() error, output below:\n"+string(output))
		sc.Ok = false
		return
	}
	sc.Output = strings.TrimSpace(string(output))

	switch sc.ErrCheckType {
	case "nonZeroExit", "":
		// This should have failed early above, but if somehow we get here
		// anyway, we can just return after setting a success status here
		sc.Ok = true
		return
	case "outputGTZero":
		if len(output) > 0 {
			if sc.OutputErrorPatternMatch == "" {
				ErrorLog(nil, "Output below:\n"+string(output))
				sc.Ok = false
				return
			} else {
				regex := regexp.MustCompile(sc.OutputErrorPatternMatch)
				if regex.MatchString(string(output)) {
					ErrorLog(nil, "Output below:\n"+string(output))
					sc.Ok = false
					return
				}
			}
		}
	default:
		// If it was a nonzero exit syscall, they should never get here anyway
		FatalLog(nil, "Unhandled Syscall.ErrCheckType '%s'", sc.ErrCheckType)
	}

	sc.Ok = true
}
