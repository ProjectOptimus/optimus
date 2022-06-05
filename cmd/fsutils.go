package cmd

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/opensourcecorp/rhad/logging"
	"golang.org/x/mod/semver"
	"gopkg.in/ini.v1"
)

// fileData tracks info about discovered files
type fileData struct {
	Path  string
	IsDir bool
}

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
		logging.Warn("No Rhadfile found, so just processing this root directory")
		rawContents, err = ini.Load([]byte("[.]"))
		if err != nil {
			logging.Error("Error while loading default Rhadfile contents")
			logging.Error(err.Error())
			os.Exit(1)
		}
	} else {
		rawContents, err = ini.Load(rhadfilePath)
		if err != nil {
			logging.Error("Error while reading Rhadfile")
			logging.Error(err.Error())
			os.Exit(1)
		}
	}

	var cfg rhadConfig
	rf := make(rhadFile)

	for _, section := range rawContents.SectionStrings() {
		err = rawContents.Section(section).MapTo(&cfg)
		if err != nil {
			logging.Error("Could not map Rhadfile contents in section '%s' to struct", section)
			logging.Error(err.Error())
			os.Exit(1)
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

// getVersion tries to build a (Go) compliant Semantic Version number out of the
// provided string, regardless of how dirty it is. Despite using the semver
// package in a few  places internally, most of this implementation is custom
// due to limitations in that package -- like not being able to parse out just
// the Patch number, or Pre-Release/Build numbers not being allowed in
// semver.Canonical()
func getVersion(s string) string {
	var v, preRelease, build string

	// Grab the semver parts separately so we can clean them up. Firstly, the
	// Major-Minor-Patch parts, but Patch takes some extra work to suss out --
	// if there's any prerelease or build parts of the version, these show up in
	// the "patch" index if we just split on dots, so we can just grab the whole
	// MMP with a regex
	v = "v" + regexp.MustCompile(`\d+(\.\d+)?(\.\d+)?`).FindStringSubmatch(s)[0]
	v = semver.Canonical(v)

	// Gross
	prSplit := strings.Split(s, "-")
	bSplit := strings.Split(s, "+")
	if len(prSplit) > 1 {
		// could still have a build number
		preRelease = strings.Split(prSplit[1], "+")[0]
	}
	if len(bSplit) > 1 {
		build = bSplit[1]
	}

	if preRelease != "" {
		v += "-" + preRelease
	}
	if build != "" {
		v += "+" + build
	}

	if !semver.IsValid(v) {
		logging.Error("Could not understand the semantic version you provided in your Rhadfile: '%s'", s)
		os.Exit(1)
	}

	if v != s {
		logging.Warn("The Semantic Version string built was different from the one provided -- please edit your version to match the correct format: %s --> %s", s, v)
	}

	return v
}

// getAllFiles returns a []fileData after traversing some root path. It will
// include directories as well as files. This is expected to be called within
// getRelevantFiles, but can be used as needed.
func getAllFiles(root string) []fileData {
	var files []fileData

	// This is such an ugly way to walk directories and get the files
	// (filepath.Glob doesn't recurse deep enough), but... Go things
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}

		if filepath.Base(path) == ".git" {
			return filepath.SkipDir
		}

		fileInfo, err := os.Stat(path)
		if err != nil {
			logging.Error("Could not process filepath %s for some reason; error specifics below\n", path)
			logging.Error(err.Error())
		}

		file := fileData{
			Path:  path,
			IsDir: fileInfo.Mode().IsDir(),
		}

		files = append(files, file)

		return nil
	})

	if err != nil {
		logging.Error("There was an error processing the directory")
		logging.Error(err.Error())
		os.Exit(1)
	}

	return files
}

// getRelevantFiles extracts files and/or directories by traversing some root
// path, based on some regular expression pattern provided.
func getRelevantFiles(root, pattern string) []fileData {
	allFiles := getAllFiles(root)
	regex := regexp.MustCompile(pattern)
	ignoreRegex := regexp.MustCompile(ignorePattern)
	var files []fileData

	for _, file := range allFiles {
		if regex.MatchString(file.Path) && !ignoreRegex.MatchString(file.Path) {
			files = append(files, file)
		}
	}
	return files
}
