package cmd

import (
	"path/filepath"
	"strings"

	cf "github.com/osamaadam/cfgrr/configfile"
	"github.com/osamaadam/cfgrr/prompt"
	"github.com/osamaadam/cfgrr/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configPatterns []string
)

var backupCmd = &cobra.Command{
	Use:     "backup",
	Short:   "Backup the configuration files to the backup directory",
	Aliases: []string{"b", "bkp"},
	Args:    cobra.ExactArgs(1),
	Example: strings.Join([]string{
		`cfgrr backup /path/to/root/config/dir`,
		`cfgrr b ~/`,
	}, "\n"),
	RunE: runBackup,
}

func runBackup(cmd *cobra.Command, args []string) error {
	// TODO: Default ignore file should be initialized on init if it doesn't exist
	// This ignore file should include files like .git, node_modules, etc.
	root := args[0]

	mapFile := viper.GetString("map-file")
	ignFile := viper.GetString("ignore-file")
	backupDir := viper.GetString("backup-dir")

	ignFilePath := filepath.Join(backupDir, ignFile)

	files, err := util.FindFiles(root, ignFilePath, configPatterns...)
	if err != nil {
		return errors.WithStack(err)
	}

	selectedFiles, err := prompt.PromptForFileSelection(files)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := cf.CopyAndReplaceFiles(backupDir, mapFile, selectedFiles...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func init() {
	defaultPatterns := []string{".*", "*config*"}
	backupCmd.Flags().StringSliceVarP(&configPatterns, "pattern", "p", defaultPatterns, "backup files matching the given pattern .")
}
