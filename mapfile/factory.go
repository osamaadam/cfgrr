package mapfile

import (
	"path/filepath"
)

func NewMapFile(path string) IMapFile {
	ext := filepath.Ext(path)

	if ext == ".json" {
		// TODO: implement JSON mapfile, maybe
	} else {
		// just default to YAML
		return NewYamlMapFile(path)
	}

	return nil
}
