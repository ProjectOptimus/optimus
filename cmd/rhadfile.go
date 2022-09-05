package cmd

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	osc "github.com/opensourcecorp/go-common"
)

// A Rhadfile represents all sections (rhadModules) together in a real Rhadfile
type Rhadfile struct {
	Modules map[string]rhadModule `toml:"module"`
}

// rhadModule defines each tree branch to traverse in the codebase
type rhadModule struct {
	Version string `toml:"version"`
}

// readRhadfile reads in the Rhadfile provided, and unpacks it into something
// usable elsewhere. Currently, Rhadfiles are in the TOML format. The hacky
// `rhadfilePath` parameter allows for internal setting of the Rhadfile path
// (such as during tests), but will default to searching the current directory
// if not provided.
func readRhadfile(customRhadfilePath ...string) (Rhadfile, toml.MetaData) {
	var err error
	var rhadfileData Rhadfile
	var metadata toml.MetaData

	var rhadfilePath string
	if len(customRhadfilePath) == 0 {
		rhadfilePath = "./Rhadfile"
	} else {
		rhadfilePath = customRhadfilePath[0]
	}

	rhadfilePath, err = filepath.Abs(rhadfilePath)
	if err != nil {
		osc.FatalLog(err, "Error when determining the absolute path of ./Rhadfile")
	}

	_, err = os.Lstat(rhadfilePath)
	if err != nil {
		osc.WarnLog("No Rhadfile found, so rhad will only process this root directory, and without configuration options. Create a top-level 'Rhadfile' with a '[module.root]' TOML table.")
		// This is a "default" Rhadfile -- we need to return one or else the
		// outer loop in the caller will fail to run without any error. The
		// module path for a single-level Rhadfile can either be "root" or ".",
		// both of which we will translate to the latter
		rhadfileData = Rhadfile{
			Modules: map[string]rhadModule{
				".": {Version: "v0.0.1"},
			},
		}
	} else {
		metadata, err = toml.DecodeFile(rhadfilePath, &rhadfileData)
		if err != nil {
			osc.FatalLog(err, "Error while reading or parsing Rhadfile")
		}

		// Another catch for when rhad may fail entirely silently -- typos like
		// specifying module blocks in TOML as 'modules.X' instead of
		// 'module.X', etc.
		if len(metadata.Undecoded()) > 0 {
			osc.WarnLog("Undecoded field in Rhadfile detected, and rhad may break -- you might have made a typo somewhere. Undecoded fields: %q", metadata.Undecoded())
		}

		// Here's where we need to replace a possible 'root' module path key
		// with a real resolvable directory name
		if _, ok := rhadfileData.Modules["root"]; ok {
			osc.InfoLog("Found 'root' module name in Rhadfile -- will treat that as the top-level directory '.'")
			rhadfileData.Modules["."] = rhadfileData.Modules["root"]
			delete(rhadfileData.Modules, "root")
		}
	}

	return rhadfileData, metadata
}
