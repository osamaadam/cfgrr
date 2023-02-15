package cmd

import (
	"os"
	"strings"

	cf "github.com/osamaadam/cfgrr/configfile"
	"github.com/osamaadam/cfgrr/core"
	"github.com/osamaadam/cfgrr/ignorefile"
	"github.com/osamaadam/cfgrr/prompt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:     "backup [root_dir] [...files]",
	Short:   "Backup the configuration files to the backup directory",
	Aliases: []string{"b", "bkp"},
	Args:    cobra.MinimumNArgs(1),
	Example: strings.Join([]string{
		`cfgrr backup /path/to/root/config/dir`,
		`cfgrr b ~/.bashrc`,
		`cfgrr b ~/.bashrc ~/.zshrc`,
		`cfgrr b ~/.config/ ~/.bashrc`,
		`cfgrr b ~/.config ~/.bashrc -a`,
		`cfgrr b ~/`,
		`cfgrr b /path/to/root/config/dir -p "**/.*" -p "**/*config*"`,
		`cfgrr b /path/to/root/config/dir -p "**/.*" -p "**/*config*" -d /path/to/backup/dir -i .cfgrrignore -m cfgrrmap.yaml`,
	}, "\n"),
	RunE: runBackup,
}

func runBackup(cmd *cobra.Command, args []string) error {
	paths := args

	if _, err := ignorefile.InitDefaultIgnoreFile(); err != nil {
		return errors.WithStack(err)
	}
	ignFiles, _ := cmd.Flags().GetStringSlice("ignore_files")
	ignContainer := ignorefile.NewIgnoresContainer(ignFiles...)

	files := make([]*cf.ConfigFile, 0)

	for _, path := range paths {
		stats, err := os.Lstat(path)
		if err != nil {
			return errors.WithStack(err)
		}

		if stats.Mode()&os.ModeSymlink == os.ModeSymlink {
			continue
		}

		if stats.IsDir() {
			fs, err := core.FindFiles(path, ignContainer, configPatterns...)
			if err != nil {
				return errors.WithStack(err)
			}

			files = append(files, fs...)
		} else {
			f, err := cf.NewConfigFile(path)
			if err != nil {
				return errors.WithStack(err)
			}

			files = append(files, f)
		}

	}

	if !all {
		if err := prompt.PromptForFileSelection(&files, "Which files would you like to track? (this will overwrite existing files)"); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := core.BackupFiles(files...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func init() {
	defaultPatterns := []string{`**/.*`, `**/*config*`}
	backupCmd.Flags().StringSliceVarP(&configPatterns, "pattern", "p", defaultPatterns, "backup files matching the given patterns")
	backupCmd.Flags().BoolVarP(&all, "all", "a", false, "backup all matched files (skip prompt)")
}
