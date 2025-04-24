package mapfile

import (
	"path/filepath"

	"github.com/osamaadam/cfgrr/vconfig"
)

// Returns a new IMapFile based on the file extension
func NewMapFile(optPath ...string) IMapFile {
	var path string
	if len(optPath) == 0 {
		config := vconfig.GetConfig()
		path = config.GetMapFilePath()
	} else {
		path = optPath[0]
	}
	ext := filepath.Ext(path)

	switch ext {
	case ".json":
		return NewJsonMapFile(path)
	case ".yml", ".yaml":
		return NewYamlMapFile(path)
	default:
		return NewYamlMapFile(path)
	}
}
