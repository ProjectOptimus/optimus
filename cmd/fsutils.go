package cmd

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/opensourcecorp/rhad/logging"
	"gopkg.in/ini.v1"
)

type fileData struct {
	path  string
	isDir bool
}

// rhadConfig represents each section's contents of a Rhadfile
type rhadConfig struct {
	Version string `ini:"version"`
	Farts   string `ini:"farts"`
}

// rhadFile represents all sections (rhadConfigs) together in a Rhadfile, with
// section names as keys
type rhadFile map[string]rhadConfig

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

func cleanSectionName(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "./", "")
	return s
}

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
			path:  path,
			isDir: fileInfo.Mode().IsDir(),
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

func getRelevantFiles(root, pattern string) []fileData {
	allFiles := getAllFiles(root)
	regex := regexp.MustCompile(pattern)
	var files []fileData
	for _, file := range allFiles {
		if regex.MatchString(file.path) {
			files = append(files, file)
		}
	}
	return files
}
