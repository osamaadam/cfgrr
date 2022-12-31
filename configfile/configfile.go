package configfile

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type ConfigFile struct {
	Name string
	Path string
	Hash string
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

	sanitizedPath, err := filepath.Rel(homedir, absPath)
	if err != nil {
		return nil, errors.WithMessage(err, "couldn't get a path relative to home dir")
	}

	file = &ConfigFile{
		Path: sanitizedPath,
		Name: filepath.Base(sanitizedPath),
	}

	if err := file.GenHash(); err != nil {
		return nil, errors.WithStack(err)
	}

	return file, nil
}

func (cf *ConfigFile) GenHash() error {
	hasher := sha1.New()
	if _, err := hasher.Write([]byte(cf.Path)); err != nil {
		return errors.WithMessage(err, "couldn't generate a hash of file path")
	}

	hash := hex.EncodeToString(hasher.Sum(nil))

	cf.Hash = hash

	return nil
}

func (cf *ConfigFile) String() string {
	return cf.Name + " - " + "(" + cf.Path + ")"
}

func ArrToString(cfs []*ConfigFile) []string {
	names := []string{}

	for _, cf := range cfs {
		names = append(names, cf.String())
	}

	return names
}
