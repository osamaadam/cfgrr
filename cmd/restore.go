package cmd

import (
	"path/filepath"
	"strings"

	cf "github.com/osamaadam/cfgrr/configfile"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var restoreCmd = &cobra.Command{
	Use:     "restore",
	Short:   "Restore the configuration files from the backup directory",
	Aliases: []string{"r", "res"},
	RunE:    restore,
	Example: strings.Join([]string{
		`cfgrr restore`,
		`cfgrr r -d /path/to/config/dir`,
		`cfgrr r -d /path/to/config/dir -m cfgrrmap.yaml`,
	}, "\n"),
	Args: cobra.NoArgs,
}

func restore(cmd *cobra.Command, args []string) error {
	dir := viper.GetString("backup_dir")

	if exists := cf.CheckFileExists(dir); !exists {
		return errors.New("the directory doesn't exist")
	}

	mapFile := viper.GetString("map_file")

	if err := cf.RestoreConfig(filepath.Join(dir, mapFile)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
