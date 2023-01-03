package configfile

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Copies files to a directory, and updates the map file.
func CopyFiles(copyDir, mapFile string, files ...*ConfigFile) error {
	for _, file := range files {
		if err := copyFile(copyDir, file); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := UpdateYamlMapFile(filepath.Join(copyDir, mapFile), files...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Copies files to a directory, replaces the old file with a symlink, and updates the map file.
func CopyAndReplaceFiles(copyDir, mapFile string, files ...*ConfigFile) error {
	for _, file := range files {
		if err := copyAndReplaceFile(copyDir, file); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := UpdateYamlMapFile(filepath.Join(copyDir, mapFile), files...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Ensures symlinks to the backup files are created at the original file locations.
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

// Creates a backup of the file in the backup directory, and replaces the file with a symlink.
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

// Creates a backup of the file in the backup directory.
func copyFile(copyDir string, file *ConfigFile) error {
	if err := EnsureDirExists(copyDir); err != nil {
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

// Creates a symlink to the file in the backup directory at the file's path (file.Path).
func restoreSymLink(backupDir string, file *ConfigFile) error {
	if exists := CheckFileExists(file.PathAbs()); exists {
		if err := os.Remove(file.PathAbs()); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := EnsureDirExists(filepath.Dir(file.PathAbs())); err != nil {
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

// Creates a directory (recursively) if it does not exist.
// Similar to `mkdir -p`.
func EnsureDirExists(path string) error {
	dir := filepath.Clean(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
