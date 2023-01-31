package core

import (
	cf "github.com/osamaadam/cfgrr/configfile"
	"github.com/osamaadam/cfgrr/mapfile"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Backs up the files to the backup directory.
// And creates a symlink to the backup files at the original file locations.
func BackupFiles(mapFilePath string, files ...*cf.ConfigFile) error {
	for _, file := range files {
		if err := file.Backup(); err != nil {
			return errors.WithStack(err)
		}
	}

	mapFile := mapfile.NewMapFile(mapFilePath)

	if err := mapFile.AddFiles(files...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Restores the files from the backup directory.
func RestoreFiles(files ...*cf.ConfigFile) error {
	for _, file := range files {
		if err := file.Restore(); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// Deletes the files from the backup directory.
func DeleteFiles(restore bool, files ...*cf.ConfigFile) error {
	for _, file := range files {
		if err := file.DeleteBackup(restore); err != nil {
			return errors.WithStack(err)
		}
	}

	yamlFilePath := viper.GetString("map_file")

	mapFile := mapfile.NewYamlMapFile(yamlFilePath)

	if err := mapFile.RemoveFiles(files...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
