package configfile

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/osamaadam/cfgrr/helpers"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
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

	if err := file.SavePerm(); err != nil {
		return nil, errors.WithMessage(err, "couldn't save file permissions")
	}

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

// Finds the backup dir from the config.
func (cf *ConfigFile) BackupDir() string {
	return viper.GetString("backup_dir")
}

// Constructs the backup file path.
func (cf *ConfigFile) BackupPath() string {
	return filepath.Join(cf.BackupDir(), cf.HashShort())
}

// Creates a symlink to the backup file.
func (cf *ConfigFile) Restore() error {
	if err := os.Symlink(cf.BackupPath(), cf.PathAbs()); err != nil {
		return errors.WithMessage(err, "couldn't create a symlink to the backup file")
	}

	return nil
}

// Deletes the backup file.
func (cf *ConfigFile) DeleteBackup(restore bool) error {
	if restore {
		if err := os.Rename(cf.BackupPath(), cf.PathAbs()); err != nil {
			return errors.WithMessagef(err, "couldn't move backup file to the original location: %s", cf.PathAbs())
		}
	} else {
		if err := os.Remove(cf.BackupPath()); err != nil {
			return errors.WithMessagef(err, "couldn't remove backup file: %s", cf.BackupPath())
		}
	}

	return nil
}

func (cf *ConfigFile) Backup() error {
	// Save the file permissions
	cf.SavePerm()

	// Ensure the backup dir exists
	if err := helpers.EnsureDirExists(cf.BackupDir()); err != nil {
		return errors.WithMessage(err, "couldn't ensure backup dir exists")
	}

	// Move the file to the backup dir
	if err := os.Rename(cf.PathAbs(), cf.BackupPath()); err != nil {
		return errors.WithMessage(err, "couldn't move file to backup dir")
	}

	// Create a symlink to the backup file
	if err := cf.Restore(); err != nil {
		return errors.WithMessage(err, "couldn't create a symlink to the backup file")
	}

	return nil
}
