package cmd

import (
	"fmt"
	"strings"

	"github.com/osamaadam/cfgrr/core"
	"github.com/osamaadam/cfgrr/helpers"
	"github.com/osamaadam/cfgrr/mapfile"
	"github.com/osamaadam/cfgrr/prompt"
	"github.com/osamaadam/cfgrr/vconfig"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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
	config := vconfig.GetConfig()
	backupDir := config.BackupDir

	if exists := helpers.CheckFileExists(backupDir); !exists {
		return errors.New("the directory doesn't exist")
	}

	mapFile := mapfile.NewMapFile(config.GetMapFilePath())

	m, err := mapFile.Parse()
	if err != nil {
		return errors.WithStack(err)
	}

	files := helpers.GetMapValues(m)

	if !all {
		files, err = prompt.PromptForFileSelection(files, "Select the files to restore: ")
		if err != nil {
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
	restoreCmd.Flags().BoolVarP(&all, "all", "a", false, "restore all files in the backup directory (skip prompt)")
}
