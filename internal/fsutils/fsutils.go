package fsutils

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/sirupsen/logrus"
)

// fileData tracks info about discovered files
type fileData struct {
	Path  string
	IsDir bool
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
			logrus.Fatalf("could not process filepath '%s': %v", path, err)
		}

		file := fileData{
			Path:  path,
			IsDir: fileInfo.Mode().IsDir(),
		}

		files = append(files, file)

		return nil
	})

	if err != nil {
		logrus.Fatalf("processing root directory '%s': %v", root, err)
	}

	return files
}

// GetRelevantFiles extracts files and/or directories by traversing some root
// path, based on some regular expression pattern provided.
func GetRelevantFiles(root, pattern string) []fileData {
	allFiles := getAllFiles(root)
	regex := regexp.MustCompile(pattern)
	var files []fileData

	for _, file := range allFiles {
		if regex.MatchString(file.Path) {
			files = append(files, file)
		}
	}
	return files
}
