package configfile

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

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
