package util

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/osamaadam/gocfgr/configfile"
	"github.com/pkg/errors"
)

func FindFiles(rootPath, ignoreFilePath string, patterns ...string) (files []*configfile.ConfigFile, err error) {
	ignoreGlobs, err := ReadFileLines(ignoreFilePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	err = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if len(ignoreGlobs) > 0 {
				for _, ign := range ignoreGlobs {
					fileName := filepath.Base(d.Name())
					if matched, err := filepath.Match(ign, fileName); err != nil {
						return err
					} else if matched {
						return filepath.SkipDir
					}
				}
				return nil
			}
		}

		for _, pattern := range patterns {
			if matched, err := filepath.Match(pattern, filepath.Base(d.Name())); err != nil {
				return err
			} else if matched {
				if isIgnored, err := checkIfIgnored(d.Name(), ignoreGlobs...); err != nil {
					return err
				} else if !isIgnored {
					configFile, err := configfile.InitFile(path)
					if err != nil {
						return errors.WithStack(err)
					}
					files = append(files, configFile)
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func checkIfIgnored(file string, patterns ...string) (bool, error) {
	for _, pattern := range patterns {
		if matches, err := filepath.Match(pattern, file); err != nil {
			return false, err
		} else if matches {
			return true, nil
		}
	}

	return false, nil
}
