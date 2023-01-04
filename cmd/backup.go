package cmd

import (
	"path/filepath"
	"strings"

	cf "github.com/osamaadam/cfgrr/configfile"
	"github.com/osamaadam/cfgrr/core"
	"github.com/osamaadam/cfgrr/ignorefile"
	"github.com/osamaadam/cfgrr/prompt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configPatterns []string
)

var backupCmd = &cobra.Command{
	Use:     "backup [root_dir]",
	Short:   "Backup the configuration files to the backup directory",
	Aliases: []string{"b", "bkp"},
	Args:    cobra.ExactArgs(1),
	Example: strings.Join([]string{
		`cfgrr backup /path/to/root/config/dir`,
		`cfgrr b ~/`,
		`cfgrr b /path/to/root/config/dir -p "**/.*" -p "**/*config*"`,
		`cfgrr b /path/to/root/config/dir -p "**/.*" -p "**/*config*" -d /path/to/backup/dir -i .cfgrrignore -m cfgrrmap.yaml`,
	}, "\n"),
	RunE: runBackup,
}

func runBackup(cmd *cobra.Command, args []string) error {
	root := args[0]

	mapFile := viper.GetString("map_file")
	ignFile := viper.GetString("ignore_file")
	backupDir := viper.GetString("backup_dir")

	ignFilePath := filepath.Join(backupDir, ignFile)

	if exists := cf.CheckFileExists(ignFilePath); !exists {
		ignorefile.InitIgnoreFile(ignFilePath)
	}

	files, err := core.FindFiles(root, ignFilePath, backupDir, configPatterns...)
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
	defaultPatterns := []string{`**/.*`, `**/*config*`}
	backupCmd.Flags().StringSliceVarP(&configPatterns, "pattern", "p", defaultPatterns, "backup files matching the given pattern .")
}
