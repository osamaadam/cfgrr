package helpers

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Checks if a file exists.
func CheckFileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// Creates a directory (recursively) if it does not exist.
// Similar to `mkdir -p`.
func EnsureDirExists(path string) error {
	dir := filepath.Clean(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.WithMessage(err, "couldn't create directory")
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

// Copies file from one place to another creating
// directory structure if needed.
func CopyFile(dest, origin string) error {
	if err := EnsureDirExists(filepath.Dir(dest)); err != nil {
		return errors.WithStack(err)
	}
	destFile, err := os.Create(dest)
	if err != nil {
		return errors.WithStack(err)
	}
	defer destFile.Close()

	originFile, err := os.Open(origin)
	if err != nil {
		return errors.WithStack(err)
	}
	defer originFile.Close()

	if _, err = io.Copy(destFile, originFile); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Creates a hard link of a file.
func LinkFile(dest, origin string) error {
	if err := EnsureDirExists(filepath.Dir(dest)); err != nil {
		return errors.WithStack(err)
	}

	if err := os.Remove(dest); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return errors.WithStack(err)
		}
	}

	if err := os.Link(origin, dest); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
