package utils

import (
	"fmt"
	"os"
	"path/filepath"

	osc "github.com/opensourcecorp/go-common"
)

func GetRhadSrc() string {
	var rhadSrc string
	var err error

	rhadSrc, ok := os.LookupEnv("RHAD_SRC")
	if !ok {
		rhadSrc, err = filepath.Abs(".")
		if err != nil {
			osc.FatalLog(err, "Error trying to determine default absolute filepath to rhad's sourcecode root")
		}
	} else {
		_, err = os.Lstat(rhadSrc)
		if err != nil {
			osc.FatalLog(err, "Env var 'RHAD_SRC' was provided, but set to a nonexistent directory")
		}
	}

	return rhadSrc
}

func CheckIsTesting() bool {
	if os.Getenv("RHAD_TESTING") == "true" {
		fmt.Println("RHAD_TESTING set to 'true', so will surpress further output")
		return true
	} else {
		return false
	}
}

// VerifySysinit can be used to run before functions that are making
// osc.Syscall.Exec()s, so they hopefully catch runtime errors earlier
func VerifySysinit() {
	sc := osc.Syscall{
		CmdLine: []string{"bash", GetRhadSrc() + "/scripts/sysinit.sh", "test"},
	}
	sc.Exec()
	if !sc.Ok {
		osc.FatalLog(nil, "This host has not been properly configured to use rhad -- please run rhad's sysinit first")
	}
}
