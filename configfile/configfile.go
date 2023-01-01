package configfile

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type ConfigFile struct {
	Path string
}

func InitFile(path string) (file *ConfigFile, err error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.WithMessage(err, "couldn't find user's home dir")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.WithMessage(err, "couldn't get an absolute path")
	}

	relPath, err := filepath.Rel(homedir, absPath)
	if err != nil {
		return nil, errors.WithMessage(err, "couldn't get a path relative to home dir")
	}

	file = &ConfigFile{
		Path: relPath,
	}

	return file, nil
}

func (cf *ConfigFile) Name() string {
	return filepath.Base(cf.Path)
}

func (cf *ConfigFile) Hash() string {
	hasher := sha1.New()
	hasher.Write([]byte(cf.Path))
	hash := hex.EncodeToString(hasher.Sum(nil))

	return hash
}

func (cf *ConfigFile) String() string {
	return cf.Name() + " - " + "(" + filepath.Join("~", cf.Path) + ")"
}
