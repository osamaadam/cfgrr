package configfile

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Backs up the files to the backup directory.
// And creates a symlink to the backup files at the original file locations.
func BackupFiles(mapFilePath string, files ...*ConfigFile) error {
	for _, file := range files {
		if err := file.Backup(); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := UpdateYamlMapFile(mapFilePath, files...); err != nil {
		errors.WithStack(err)
	}

	return nil
}

// Restores the files from the backup directory.
func RestoreFiles(files ...*ConfigFile) error {
	for _, file := range files {
		if err := file.RestoreSymlink(); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// Deletes the files from the backup directory.
func DeleteFiles(restore bool, files ...*ConfigFile) error {
	for _, file := range files {
		if err := file.DeleteBackup(restore); err != nil {
			return errors.WithStack(err)
		}
	}

	yamlFilePath := viper.GetString("map_file")

	if err := tidyYamlMapFile(yamlFilePath); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func CheckFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Creates a directory (recursively) if it does not exist.
// Similar to `mkdir -p`.
func EnsureDirExists(path string) error {
	dir := filepath.Clean(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Checks if the given path is a path to a symlink.
func CheckIfSymlink(path string) (bool, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return fi.Mode()&os.ModeSymlink != 0, nil
}
