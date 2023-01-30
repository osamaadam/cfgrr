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
	Perm os.FileMode
}

/*
Tidies the path before initializing the object.

	cf, _ := NewConfigFile("~/path/../path/.config")
	// cf.Path = "path/.config"
*/
func NewConfigFile(path string) (file *ConfigFile, err error) {
	if path == "" {
		return nil, errors.New("path can't be empty")
	}

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

	file.SavePerm()

	return file, nil
}

/*
Returns the base name if the file.

For example if the Path = "/some/loc/somewhere/.config" ->
Name = ".config"
*/
func (cf *ConfigFile) Name() string {
	return filepath.Base(cf.Path)
}

// Returns the absolute path of the file.
// Relies on there being a $HOME environment variable.
func (cf *ConfigFile) PathAbs() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homedir, cf.Path)
}

// Returns the hash of the Path.
func (cf *ConfigFile) Hash() string {
	hasher := sha1.New()
	hasher.Write([]byte(cf.Path))
	hash := hex.EncodeToString(hasher.Sum(nil))

	return hash
}

// Returns a truncated hash of the Path.
func (cf *ConfigFile) HashShort() string {
	return cf.Hash()[:8]
}

// Makes it printable, functions like fmt.Println know to call this automatically.
func (cf *ConfigFile) String() string {
	return cf.Name() + " - " + "(" + filepath.Join("~", cf.Path) + ")"
}

// Check if the file exists.
func (cf *ConfigFile) Exists() bool {
	_, err := os.Stat(cf.PathAbs())
	return !os.IsNotExist(err)
}

// Save file permissions.
func (cf *ConfigFile) SavePerm() error {
	info, err := os.Stat(cf.PathAbs())
	if err != nil {
		if os.IsNotExist(err) {
			cf.Perm = os.FileMode(0644)
			return nil
		}
		return errors.WithStack(err)
	}

	cf.Perm = info.Mode()

	return nil
}
