package mapfile

import (
	"fmt"
	"os"
	"path/filepath"

	cf "github.com/osamaadam/cfgrr/configfile"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type YamlMapFile struct {
	path string
}

func NewYamlMapFile(path string) *YamlMapFile {
	path = filepath.Clean(path)

	if filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml" {
		path = filepath.Join(filepath.Dir(path), filepath.Base(path)+".yaml")
	}

	return &YamlMapFile{path: path}
}

// Opens the map file.
func (yf *YamlMapFile) open() (*os.File, error) {
	return os.OpenFile(yf.path, os.O_RDWR|os.O_CREATE, os.FileMode(0644))
}

// Writes the map to the map file.
func (yf *YamlMapFile) write(m map[string]*cf.ConfigFile) error {
	marshalledData, err := yaml.Marshal(&m)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := cf.EnsureDirExists(filepath.Dir(yf.path)); err != nil {
		return errors.WithStack(err)
	}

	if err := os.WriteFile(yf.path, marshalledData, os.FileMode(0644)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Print in string format.
func (yf *YamlMapFile) String() string {
	m, _ := yf.Parse()
	return fmt.Sprintf("%v", m)
}

// Parses the map file into a `map[string]*cf.ConfigFile`.
func (yf *YamlMapFile) Parse() (mf map[string]*cf.ConfigFile, err error) {
	file, err := yf.open()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&mf); err != nil {
		return nil, errors.WithStack(err)
	}

	return
}

// Adds files to the map file.
func (yf *YamlMapFile) AddFiles(files ...*cf.ConfigFile) error {
	m, err := yf.Parse()
	if err != nil {
		return errors.WithStack(err)
	}

	for _, file := range files {
		m[file.HashShort()] = file
	}

	if err := yf.write(m); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Removes files from the map file.
func (yf *YamlMapFile) RemoveFiles(files ...*cf.ConfigFile) error {
	m, err := yf.Parse()
	if err != nil {
		return errors.WithStack(err)
	}

	for _, file := range files {
		delete(m, file.HashShort())
	}

	if err := yf.write(m); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Removes files from the map file that don't exist in the backup directory.
func (yf *YamlMapFile) Tidy() error {
	m, err := yf.Parse()
	if err != nil {
		return errors.WithStack(err)
	}

	backupDir := viper.GetString("backup_dir")

	for _, file := range m {
		filePath := filepath.Join(backupDir, file.HashShort())
		if !cf.CheckFileExists(filePath) {
			delete(m, file.HashShort())
		}
	}

	if err := yf.write(m); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
