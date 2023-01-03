package core

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/mattn/go-zglob"
	"github.com/osamaadam/cfgrr/configfile"
	"github.com/pkg/errors"
)

// FindFiles finds files in the given rootPath that match the given patterns.
func FindFiles(rootPath, ignoreFilePath, backupDir string, patterns ...string) (files []*configfile.ConfigFile, err error) {
	if len(patterns) == 0 {
		return nil, errors.New("no patterns given")
	}

	ignoreGlobs, err := ReadFileLines(ignoreFilePath)

	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	ignoreGlobs = append(ignoreGlobs, backupDir)

	err = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}

		if d.IsDir() {
			// Check if the directory is ignored.
			if len(ignoreGlobs) > 0 {
				if ignored := checkIfMatches(path, ignoreGlobs...); ignored {
					return filepath.SkipDir
				}
			}

			return nil
		}

		// Check if file matches any of the given patterns.
		if matches := checkIfMatches(path, patterns...); matches {
			// Check if file is ignored.
			if ignored := checkIfMatches(path, ignoreGlobs...); !ignored {
				file, err := configfile.InitFile(path)
				if err != nil {
					return errors.WithStack(err)
				}

				files = append(files, file)
			}
		}

		return nil
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return files, nil
}

// Checks if the given file matches any of the given patterns.
func checkIfMatches(file string, patterns ...string) bool {
	for _, pattern := range patterns {
		if matches, _ := zglob.Match(pattern, file); matches {
			return true
		}
	}

	return false
}
