package configfile

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Updates the map file with the new files.
func UpdateYamlMapFile(path string, files ...*ConfigFile) error {
	m := make(map[string]*ConfigFile, len(files))
	if exists := CheckFileExists(path); exists {
		readMap, err := ReadYamlMapFile(path)
		if err != nil {
			return errors.WithStack(err)
		}
		m = readMap
	}

	for _, file := range files {
		m[file.HashShort()] = file
	}

	if err := writeYamlFileRaw(path, m); err != nil {
		return errors.WithStack(err)
	}

	if err := tidyYamlMapFile(path); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Reads the map file and returns a map of the files.
func ReadYamlMapFile(path string) (map[string]*ConfigFile, error) {
	m := make(map[string]*ConfigFile)

	file, err := os.Open(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&m); err != nil {
		return nil, errors.WithStack(err)
	}

	return m, nil
}

// Restores files from the backup directory to their original locations.
// This will run tidyYamlMapFile() to remove any files that are no longer present.
func RestoreConfig(path string) error {
	if err := tidyYamlMapFile(path); err != nil {
		return errors.WithStack(err)
	}

	m, err := ReadYamlMapFile(path)
	if err != nil {
		return errors.WithStack(err)
	}

	fileArray := make([]*ConfigFile, 0, len(m))

	for _, file := range m {
		fileArray = append(fileArray, file)
	}

	if err := RestoreSymLinks(filepath.Dir(path), fileArray...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Writes a yaml file to the specified path.
func writeYamlFileRaw(path string, m interface{}) error {
	marshalledData, err := yaml.Marshal(&m)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := EnsureDirExists(filepath.Dir(path)); err != nil {
		return errors.WithStack(err)
	}

	if err := os.WriteFile(correctYamlFileName(path), marshalledData, 0644); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Deletes files from the map file that are no longer present.
func tidyYamlMapFile(path string) error {
	m, err := ReadYamlMapFile(path)
	if err != nil {
		return errors.WithStack(err)
	}

	baseDir := filepath.Dir(path)

	for _, file := range m {
		filePath := filepath.Join(baseDir, file.HashShort())
		if !CheckFileExists(filePath) {
			delete(m, file.HashShort())
		}
	}

	return writeYamlFileRaw(path, m)
}

// Ensures the yaml file has the correct extension.
func correctYamlFileName(path string) string {
	path = filepath.Clean(path)
	ext := filepath.Ext(path)
	if ext != ".yaml" && ext != ".yml" {
		path = path + ".yaml"
	}

	return path
}
