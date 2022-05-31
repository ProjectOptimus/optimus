package cmd

import (
	"testing"
)

var (
	testSyscallCfg syscallCfg
	testSyscallOk  bool
)

func TestSyscall(t *testing.T) {
	// Did it run successfully?
	testSyscallCfg = syscallCfg{
		[]string{"echo", "hello"},
		"nonZeroExit",
		"",
	}
	testSyscallOk = syscall(testSyscallCfg)
	if !testSyscallOk {
		t.Errorf("Expected 'echo hello' to succeed")
	}

	// Did it fail if it had output at all?
	testSyscallCfg = syscallCfg{
		[]string{"echo", "hello"},
		"outputGTZero",
		"",
	}
	testSyscallOk = syscall(testSyscallCfg)
	if testSyscallOk {
		t.Errorf("Expected 'echo hello' to have command output")
	}

	// Did it fail if it had a *specific* output regex match?
	testSyscallCfg = syscallCfg{
		[]string{"echo", "hello", "folks"},
		"outputGTZero",
		"folks",
	}
	testSyscallOk = syscall(testSyscallCfg)
	if testSyscallOk {
		t.Errorf("Expected 'echo hello folks' to have 'folks' in the captured output")
	}
}
