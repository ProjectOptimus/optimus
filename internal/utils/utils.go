package utils

import (
	"fmt"
	"os"
	"path/filepath"

	osc "github.com/opensourcecorp/go-common"
	"github.com/sirupsen/logrus"
)

func GetOscarSrc() string {
	var oscarSrc string
	var err error

	oscarSrc, ok := os.LookupEnv("OSCAR_SRC")
	if !ok {
		oscarSrc, err = filepath.Abs(".")
		if err != nil {
			logrus.Fatalf("Error trying to determine default absolute filepath to oscar's sourcecode root: %v", err)
		}
	} else {
		_, err = os.Lstat(oscarSrc)
		if err != nil {
			logrus.Fatalf("Env var 'OSCAR_SRC' was provided, but set to a nonexistent directory: %v", err)
		}
	}

	return oscarSrc
}

// VerifySysinit can be used to run before functions that are making
// osc.Syscall.Exec()s, so they hopefully catch runtime errors earlier
func VerifySysinit() error {
	sc := osc.Syscall{
		CmdLine: []string{"bash", GetOscarSrc() + "/scripts/sysinit.sh", "test"},
	}
	sc.Exec()
	if !sc.Ok {
		return fmt.Errorf("this host has not been properly configured to use oscar -- please run oscar's sysinit first")
	}
	return nil
}
