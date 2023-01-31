package ignorefile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/osamaadam/cfgrr/helpers"
	"github.com/pkg/errors"
)

type IIgnoreFile interface {
	fmt.Stringer
	Read() ([]string, error)
}

type IgnoreFile struct {
	path string
}

func NewIgnoreFile(path string) *IgnoreFile {
	return &IgnoreFile{path: path}
}

func (i *IgnoreFile) String() string {
	r, err := os.ReadFile(i.path)
	if err != nil {
		return ""
	}
	return string(r)
}

func (i *IgnoreFile) Read() ([]string, error) {
	lines, err := helpers.ReadFileLines(i.path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return lines, nil
}

func InitIgnoreFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := helpers.EnsureDirExists(filepath.Dir(path)); err != nil {
			return errors.WithStack(err)
		}
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := file.Write([]byte(strings.Join(defaultIgnores, "\n"))); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
