package configfile

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Copies files to a directory, and updates the map file.
func CopyFiles(copyDir string, files ...*ConfigFile) error {
	for _, file := range files {
		if err := copyFile(copyDir, file); err != nil {
			return errors.WithStack(err)
		}
	}

	mapFile := viper.GetString("map-file")

	if err := UpdateYamlMapFile(filepath.Join(copyDir, mapFile), files...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Copies files to a directory, replaces the old file with a symlink, and updates the map file.
func CopyAndReplaceFiles(copyDir string, files ...*ConfigFile) error {
	for _, file := range files {
		if err := copyAndReplaceFile(copyDir, file); err != nil {
			return errors.WithStack(err)
		}
	}

	mapFile := viper.GetString("map-file")

	if err := UpdateYamlMapFile(filepath.Join(copyDir, mapFile), files...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func RestoreSymLinks(backupDir string, files ...*ConfigFile) error {
	for _, file := range files {
		if err := restoreSymLink(backupDir, file); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func CheckFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func copyAndReplaceFile(copyDir string, file *ConfigFile) error {
	if err := copyFile(copyDir, file); err != nil {
		return errors.WithStack(err)
	}

	if err := os.Remove(file.PathAbs()); err != nil {
		return errors.WithStack(err)
	}

	if err := os.Symlink(filepath.Join(copyDir, file.HashShort()), file.PathAbs()); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func copyFile(copyDir string, file *ConfigFile) error {
	if err := ensureDirExists(copyDir); err != nil {
		return errors.WithStack(err)
	}

	orgFile, err := os.Open(file.PathAbs())
	if err != nil {
		return errors.WithStack(err)
	}
	defer orgFile.Close()

	dstFile, err := os.Create(filepath.Join(copyDir, file.HashShort()))
	if err != nil {
		return errors.WithStack(err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, orgFile); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func restoreSymLink(backupDir string, file *ConfigFile) error {
	if exists := CheckFileExists(file.PathAbs()); exists {
		if err := os.Remove(file.PathAbs()); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := ensureDirExists(filepath.Dir(file.PathAbs())); err != nil {
		return errors.WithStack(err)
	}

	absTargPath, err := filepath.Abs(filepath.Join(backupDir, file.HashShort()))
	if err != nil {
		return errors.WithStack(err)
	}

	if err := os.Symlink(absTargPath, file.PathAbs()); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
