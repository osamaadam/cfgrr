package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/osamaadam/cfgrr/core"
	"github.com/osamaadam/cfgrr/helpers"
	"github.com/osamaadam/cfgrr/mapfile"
	"github.com/osamaadam/cfgrr/prompt"
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
		`cfgrr restore -a`,
		`cfgrr r -d /path/to/config/dir`,
		`cfgrr r -d /path/to/config/dir -m cfgrrmap.yaml`,
	}, "\n"),
	Args: cobra.NoArgs,
}

func restore(cmd *cobra.Command, args []string) error {
	backupDir := viper.GetString("backup_dir")

	if exists := helpers.CheckFileExists(backupDir); !exists {
		return errors.New("the directory doesn't exist")
	}

	mapFileName := viper.GetString("map_file")
	mapFilePath := filepath.Join(backupDir, mapFileName)

	mapFile := mapfile.NewMapFile(mapFilePath)

	m, err := mapFile.Parse()
	if err != nil {
		return errors.WithStack(err)
	}

	files := helpers.GetMapValues(m)

	if !all {
		if err = prompt.PromptForFileSelection(&files, "Select the files to restore: "); err != nil {
			return errors.WithStack(err)
		}
	}

	if len(files) == 0 {
		fmt.Println("No files selected, terminating...")
		return nil
	}

	if err := core.RestoreFiles(files...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func init() {
	restoreCmd.Flags().BoolVarP(&all, "all", "a", false, "restore all files in the backup directory")
}
