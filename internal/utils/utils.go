package utils

import (
	"fmt"
	"os"
	"path/filepath"

	osc "github.com/opensourcecorp/go-common"
)

func GetOscarSrc() string {
	var oscarSrc string
	var err error

	oscarSrc, ok := os.LookupEnv("RHAD_SRC")
	if !ok {
		oscarSrc, err = filepath.Abs(".")
		if err != nil {
			osc.FatalLog(err, "Error trying to determine default absolute filepath to oscar's sourcecode root")
		}
	} else {
		_, err = os.Lstat(oscarSrc)
		if err != nil {
			osc.FatalLog(err, "Env var 'RHAD_SRC' was provided, but set to a nonexistent directory")
		}
	}

	return oscarSrc
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
		CmdLine: []string{"bash", GetOscarSrc() + "/scripts/sysinit.sh", "test"},
	}
	sc.Exec()
	if !sc.Ok {
		osc.FatalLog(nil, "This host has not been properly configured to use oscar -- please run oscar's sysinit first")
	}
}
