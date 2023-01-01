package configfile

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func CreateYamlFile(path string, files ...*ConfigFile) error {
	m := make(map[string]*ConfigFile, len(files))

	for _, file := range files {
		m[file.HashShort()] = file
	}

	marshalledData, err := yaml.Marshal(&m)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := ensureDirExists(path); err != nil {
		return errors.WithStack(err)
	}

	if err := os.WriteFile(correctFileName(path), marshalledData, 0644); err != nil {
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

func correctFileName(path string) string {
	path = filepath.Clean(path)
	ext := filepath.Ext(path)
	if ext != ".yaml" && ext != ".yml" {
		path = path + ".yaml"
	}

	return path
}

func ensureDirExists(path string) error {
	path = filepath.Clean(path)
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
