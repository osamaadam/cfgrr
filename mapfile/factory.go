package mapfile

import (
	"path/filepath"

	"github.com/spf13/viper"
)

// Returns a new IMapFile based on the file extension
func NewMapFile(optPath ...string) IMapFile {
	var path string
	if len(optPath) == 0 {
		path = guessMapFilePath()
	} else {
		path = optPath[0]
	}
	ext := filepath.Ext(path)

	switch ext {
	case ".json":
		// TODO: implement JSON mapfile, maybe
	default:
		return NewYamlMapFile(path)
	}

	return nil
}

func guessMapFilePath() string {
	backupDir := viper.GetString("backup_dir")
	mapFileName := viper.GetString("map_file")
	return filepath.Join(backupDir, mapFileName)
}
