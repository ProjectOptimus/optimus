package oscarfile

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

// A Oscarfile represents all sections (oscarModules) together in a real Oscarfile
type Oscarfile struct {
	Modules map[string]oscarModule `toml:"module"`
}

// oscarModule defines each tree branch to traverse in the codebase
type oscarModule struct {
	Version string `toml:"version"`
}

// ReadOscarfile reads in the Oscarfile provided, and unpacks it into something
// usable elsewhere. Currently, Oscarfiles are in the TOML format. The hacky
// `oscarfilePath` parameter allows for internal setting of the Oscarfile path
// (such as during tests), but will default to searching the current directory
// if not provided.
func ReadOscarfile(customOscarfilePath ...string) (Oscarfile, toml.MetaData) {
	var err error
	var oscarfileData Oscarfile
	var metadata toml.MetaData

	var oscarfilePath string
	if len(customOscarfilePath) == 0 {
		oscarfilePath = "./Oscarfile"
	} else {
		oscarfilePath = customOscarfilePath[0]
	}

	oscarfilePath, err = filepath.Abs(oscarfilePath)
	if err != nil {
		logrus.Fatalf("determining the absolute path of ./Oscarfile: %v", err)
	}

	_, err = os.Lstat(oscarfilePath)
	if err != nil {
		logrus.Warn("No Oscarfile found, so oscar will only process this root directory, and without configuration options. You should create a top-level 'Oscarfile' with a '[module.root]' TOML table.")
		// This is a "default" Oscarfile -- we need to return one or else the
		// outer loop in the caller will fail to run without any error. The
		// module path for a single-level Oscarfile can either be "root" or ".",
		// both of which we will translate to the latter
		oscarfileData = Oscarfile{
			Modules: map[string]oscarModule{
				".": {Version: "v0.0.1"},
			},
		}
	} else {
		metadata, err = toml.DecodeFile(oscarfilePath, &oscarfileData)
		if err != nil {
			logrus.Fatalf("reading or parsing Oscarfile: %v", err)
		}

		// Another catch for when oscar may fail entirely silently -- typos like
		// specifying module blocks in TOML as 'modules.X' instead of
		// 'module.X', etc.
		if len(metadata.Undecoded()) > 0 {
			logrus.Warnf("Undecoded field in Oscarfile detected, and oscar may break -- you might have made a typo somewhere. Undecoded fields: %q", metadata.Undecoded())
		}

		// Here's where we need to replace a possible 'root' module path key
		// with a real resolvable directory name
		if _, ok := oscarfileData.Modules["root"]; ok {
			logrus.Info("Found 'root' module name in Oscarfile -- will treat that as the top-level directory '.'")
			oscarfileData.Modules["."] = oscarfileData.Modules["root"]
			delete(oscarfileData.Modules, "root")
		}
	}

	return oscarfileData, metadata
}
