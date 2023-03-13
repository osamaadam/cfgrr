package ignorefile

import (
	"fmt"
	"path/filepath"

	"github.com/osamaadam/cfgrr/vconfig"
	"github.com/pkg/errors"
)

type IIgnoresContainer interface {
	fmt.Stringer
	AddIgnoreFile(IIgnoreFile)
	ReadLines() ([]string, error)
	Paths() []string
}

type IgnoresContainer struct {
	ignFiles []IIgnoreFile
}

func NewIgnoresContainer(names ...string) IIgnoresContainer {
	config := vconfig.GetConfig()
	baseDirs := []string{
		"./",
		config.BackupDir,
	}

	ic := &IgnoresContainer{}
	for _, name := range names {
		for _, baseDir := range baseDirs {
			path := filepath.Join(baseDir, name)
			i := NewIgnoreFile(path)
			ic.AddIgnoreFile(i)
		}
	}
	return ic
}

func (ic *IgnoresContainer) String() string {
	r, _ := ic.ReadLines()
	return fmt.Sprintf("%v", r)
}

func (ic *IgnoresContainer) AddIgnoreFile(ignFile IIgnoreFile) {
	ic.ignFiles = append(ic.ignFiles, ignFile)
}

func (ic *IgnoresContainer) ReadLines() ([]string, error) {
	var igns []string
	for _, ignFile := range ic.ignFiles {
		ignores, err := ignFile.ReadLines()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		igns = append(igns, ignores...)
	}
	return igns, nil
}

func (ic *IgnoresContainer) Paths() []string {
	paths := make([]string, len(ic.ignFiles))
	for i := 0; i < len(ic.ignFiles); i++ {
		paths[i] = ic.ignFiles[i].Path()
	}
	return paths
}
