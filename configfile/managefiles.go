package configfile

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func CopyFile(copyDir string, file *ConfigFile) error {
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

func CopyFiles(copyDir string, files ...*ConfigFile) error {
	for _, file := range files {
		if err := CopyFile(copyDir, file); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := UpdateYamlFile(filepath.Join(copyDir, ".gocfgr.yaml"), files...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func CheckFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
