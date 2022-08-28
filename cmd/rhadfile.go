package cmd

import (
	"os"
	"strings"

	osc "github.com/opensourcecorp/go-common"
	"gopkg.in/ini.v1"
)

// rhadConfig represents each section's contents of a Rhadfile
type rhadConfig struct {
	Version string `ini:"version"`
}

// rhadFile represents all sections (rhadConfigs) together in a Rhadfile, with
// section names as keys
type rhadFile map[string]rhadConfig

// readRhadfile reads in the Rhadfile provided, and unmarshalls it into
// something usable elsewhere. Currently, Rhadfiles are in the INI format.
func readRhadfile() rhadFile {
	var err error
	var rawContents *ini.File

	rhadfilePath := "./Rhadfile"
	_, err = os.Lstat(string(rhadfilePath))
	if err != nil {
		osc.WarnLog("No Rhadfile found, so just processing this root directory")
		rawContents, err = ini.Load([]byte("[.]"))
		if err != nil {
			osc.FatalLog(err, "Error while loading default Rhadfile contents")
		}
	} else {
		rawContents, err = ini.Load(rhadfilePath)
		if err != nil {
			osc.FatalLog(err, "Error while reading Rhadfile")
		}
	}

	var cfg rhadConfig
	rf := make(rhadFile)

	for _, section := range rawContents.SectionStrings() {
		err = rawContents.Section(section).MapTo(&cfg)
		if err != nil {
			osc.FatalLog(err, "Could not map Rhadfile contents in section '%s' to struct", section)
		}

		sectionClean := cleanSectionName(section)

		rf[sectionClean] = cfg
	}

	return rf
}

// cleanSectionName takes INI section headers, and cleans them into a consistent
// format for reuse (the INI specification doesn't enforce spacing, etc. in
// section names)
func cleanSectionName(s string) string {
	for _, e := range []string{"[", "]", "./"} {
		s = strings.ReplaceAll(s, e, "")
	}
	s = strings.TrimSpace(s)
	s = "[" + s + "]" // gotta put them back for the INI package
	return s
}
