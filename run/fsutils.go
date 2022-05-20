package run

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/opensourcecorp/rhadamanthus/logging"
)

type fileData struct {
	path  string
	isDir bool
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

func getRelevantFiles(pattern string, cfg rhadConfig) []fileData {
	allFiles := getAllFiles(*cfg.cliOpts.path)
	regex := regexp.MustCompile(pattern)
	var files []fileData
	for _, file := range allFiles {
		if regex.MatchString(file.path) {
			files = append(files, file)
		}
	}
	return files
}

func readConfig() cfgFileData {
	return cfgFileData{"a": "a"}
}
