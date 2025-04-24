package mapfile

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"

	cf "github.com/osamaadam/cfgrr/configfile"
	"github.com/osamaadam/cfgrr/helpers"
	"github.com/pkg/errors"
)

type JsonMapFile struct {
	path string
}

func NewJsonMapFile(path string) *JsonMapFile {
	path = filepath.Clean(path)

	ext := filepath.Ext(path)
	base := filepath.Base(path)
	okayExts := []string{".json"}

	if !slices.Contains(okayExts, ext) {
		path = filepath.Join(filepath.Dir(path), base+".json")
	}

	return &JsonMapFile{path: path}
}

func (jf *JsonMapFile) Path() string {
	return jf.path
}

// Opens the map file.
func (jf *JsonMapFile) open() (*os.File, error) {
	if err := helpers.EnsureDirExists(filepath.Dir(jf.path)); err != nil {
		return nil, errors.WithMessage(err, "couldn't create base directory for json file")
	}
	return os.OpenFile(jf.path, os.O_RDWR|os.O_CREATE, os.FileMode(0644))
}

// Writes the map to the map file.
func (jf *JsonMapFile) write(m map[string]*cf.ConfigFile) error {
	marshalledData, err := json.MarshalIndent(&m, "", "  ")
	if err != nil {
		return errors.WithStack(err)
	}

	if err := helpers.EnsureDirExists(filepath.Dir(jf.path)); err != nil {
		return errors.WithStack(err)
	}

	if err := os.WriteFile(jf.path, marshalledData, os.FileMode(0644)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Get the backupDir.
func (jf *JsonMapFile) backupDir() string {
	return filepath.Dir(jf.path)
}

// Print in string format.
func (jf *JsonMapFile) String() string {
	m, _ := jf.Parse()
	return fmt.Sprintf("%v", m)
}

// Parses the map file into a `map[string]*cf.ConfigFile`.
func (jf *JsonMapFile) Parse() (mf map[string]*cf.ConfigFile, err error) {
	file, err := jf.open()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer file.Close()

	if err := json.NewDecoder(file).Decode(&mf); err != nil {
		if errors.Is(err, io.EOF) {
			return map[string]*cf.ConfigFile{}, nil
		}
		return nil, errors.WithStack(err)
	}

	return
}

// Adds files to the map file.
func (jf *JsonMapFile) AddFiles(files ...*cf.ConfigFile) error {
	m, err := jf.Parse()
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

	if err := jf.write(m); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Removes files from the map file.
func (jf *JsonMapFile) RemoveFiles(files ...*cf.ConfigFile) error {
	m, err := jf.Parse()
	if err != nil {
		return errors.WithStack(err)
	}

	for _, file := range files {
		delete(m, file.HashShort())
	}

	if err := jf.write(m); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Removes files from the map file that don't exist in the backup directory.
func (jf *JsonMapFile) Tidy() error {
	m, err := jf.Parse()
	if err != nil {
		return errors.WithStack(err)
	}

	for _, file := range m {
		if !helpers.CheckFileExists(file.BackupPath()) {
			delete(m, file.HashShort())
		}
	}

	if err := jf.write(m); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
