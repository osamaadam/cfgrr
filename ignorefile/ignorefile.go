package ignorefile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/osamaadam/cfgrr/helpers"
	"github.com/osamaadam/cfgrr/vconfig"
	"github.com/pkg/errors"
)

// TODO: Flesh out the ignore file component of the app.

type IIgnoreFile interface {
	fmt.Stringer
	Write(...string) error
	Read() ([]string, error)
	Path() string
}

type IgnoreFile struct {
	path string
}

func NewIgnoreFile(path string) IIgnoreFile {
	return &IgnoreFile{path: path}
}

func (i *IgnoreFile) String() string {
	r, err := os.ReadFile(i.path)
	if err != nil {
		return ""
	}
	return string(r)
}

func (i *IgnoreFile) Write(lines ...string) error {
	readLines, err := i.Read()
	if err != nil && os.IsNotExist(err) {
		return errors.WithStack(err)
	}

	lines = append(lines, readLines...)
	sort.Strings(lines)

	if err := helpers.EnsureDirExists(filepath.Dir(i.path)); err != nil {
		return errors.WithStack(err)
	}

	file, err := os.OpenFile(i.Path(), os.O_RDWR|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	if _, err := file.Write([]byte(strings.Join(lines, "\n"))); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (i *IgnoreFile) Read() ([]string, error) {
	lines, err := helpers.ReadFileLines(i.path)
	if err != nil {
		return []string{}, errors.WithStack(err)
	}

	return lines, nil
}

func (i *IgnoreFile) Path() string {
	return i.path
}

func InitDefaultIgnoreFile() (IIgnoreFile, error) {
	c := vconfig.GetConfig()
	ign := NewIgnoreFile(c.GetIgnoreFilePath())

	lines, err := ign.Read()
	if err != nil && os.IsNotExist(err) {
		return nil, errors.WithStack(err)
	}

	if len(lines) == 0 {
		if err := ign.Write(defaultIgnores...); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return ign, nil
}
