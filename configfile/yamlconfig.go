package configfile

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func UpdateYamlFile(path string, files ...*ConfigFile) error {
	m := make(map[string]*ConfigFile, len(files))
	if exists := checkFileExists(path); exists {
		readMap, err := ReadYamlFile(path)
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

	if err := tidyYamlFile(path); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func ReadYamlFile(path string) (map[string]*ConfigFile, error) {
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

func RestoreConfig(path string) error {
	if err := tidyYamlFile(path); err != nil {
		return errors.WithStack(err)
	}

	m, err := ReadYamlFile(path)
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

func writeYamlFileRaw(path string, m interface{}) error {
	marshalledData, err := yaml.Marshal(&m)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := ensureDirExists(filepath.Dir(path)); err != nil {
		return errors.WithStack(err)
	}

	if err := os.WriteFile(correctFileName(path), marshalledData, 0644); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func tidyYamlFile(path string) error {
	m, err := ReadYamlFile(path)
	if err != nil {
		return errors.WithStack(err)
	}

	baseDir := filepath.Dir(path)

	for _, file := range m {
		filePath := filepath.Join(baseDir, file.HashShort())
		if !checkFileExists(filePath) {
			delete(m, file.HashShort())
		}
	}

	return writeYamlFileRaw(path, m)
}

func correctFileName(path string) string {
	path = filepath.Clean(path)
	ext := filepath.Ext(path)
	if ext != ".yaml" && ext != ".yml" {
		path = path + ".yaml"
	}

	return path
}

func ensureDirExists(path string) error {
	dir := filepath.Clean(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
