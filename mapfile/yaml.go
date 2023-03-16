package mapfile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	cf "github.com/osamaadam/cfgrr/configfile"
	"github.com/osamaadam/cfgrr/helpers"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type YamlMapFile struct {
	path string
}

func NewYamlMapFile(path string) *YamlMapFile {
	path = filepath.Clean(path)

	ext := filepath.Ext(path)
	base := filepath.Base(path)
	okayExts := []string{".yaml", ".yml"}

	if !helpers.Contains(okayExts, ext) {
		path = filepath.Join(filepath.Dir(path), base+".yaml")
	}

	return &YamlMapFile{path: path}
}

func (yf *YamlMapFile) Path() string {
	return yf.path
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

	if err := helpers.EnsureDirExists(filepath.Dir(yf.path)); err != nil {
		return errors.WithStack(err)
	}

	if err := os.WriteFile(yf.path, marshalledData, os.FileMode(0644)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Get the backupDir.
func (yf *YamlMapFile) backupDir() string {
	return filepath.Dir(yf.path)
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
		if errors.Is(err, io.EOF) {
			return map[string]*cf.ConfigFile{}, nil
		}
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
		fileHash := file.HashShort()
		if f, ok := m[fileHash]; ok {
			file.Browsable = f.Browsable || file.Browsable
		}
		m[fileHash] = file
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

	for _, file := range m {
		if !helpers.CheckFileExists(file.BackupPath()) {
			delete(m, file.HashShort())
		}
	}

	if err := yf.write(m); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
