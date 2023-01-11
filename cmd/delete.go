package cmd

import (
	"path/filepath"
	"strings"

	cf "github.com/osamaadam/cfgrr/configfile"
	"github.com/osamaadam/cfgrr/prompt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	force   bool
	replace bool
)

var deleteCmd = &cobra.Command{
	Use:     "delete [...paths]",
	Short:   "Delete the configuration files from the backup directory",
	Long:    `Delete the configuration files from the backup directory, also replacing symlinks with the original target.`,
	Aliases: []string{"d", "del"},
	RunE:    deleteRun,
	Example: strings.Join([]string{
		"cfgrr delete",
		"cfgrr delete ~/.vimrc",
		"cfgrr delete ~/.vimrc ~/.zshrc",
		"cfgrr delete -r ~/.vimrc",
		"cfgrr delete -rf ~/.vimrc",
	}, "\n"),
}

func deleteRun(cmd *cobra.Command, args []string) (err error) {
	files := make([]*cf.ConfigFile, 0)

	for _, path := range args {
		file, err := cf.NewConfigFile(path)
		if err != nil {
			return err
		}
		files = append(files, file)
	}

	backupDir := viper.GetString("backup_dir")
	mapFileName := viper.GetString("map_file")

	if len(files) == 0 {
		// User didn't specify any files, so we'll prompt them to select some
		// from the map file.
		files, err = cf.FindFilesToRestore(filepath.Join(backupDir, mapFileName))
		if err != nil {
			return errors.WithStack(err)
		}

		files, err = prompt.PromptForFileSelection(files, "Select the files to delete: ")
		if err != nil {
			return errors.WithStack(err)
		}
	}

	if len(files) == 0 {
		return nil
	}

	if replace {
		// User wants to replace the symlinks with the original target.
		if err := cf.RemoveFilesAndRevert(backupDir, mapFileName, force, files...); err != nil {
			return errors.WithStack(err)
		}
	} else {
		// User just wants to remove the files from the backup directory.
		if err := cf.RemoveFiles(backupDir, mapFileName, files...); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func init() {
	deleteCmd.Flags().BoolVarP(&force, "force", "f", false, "force replace the original files even if they're not symlinks")
	deleteCmd.Flags().BoolVarP(&replace, "replace", "r", false, "replace the symlinks with the original target")
}
