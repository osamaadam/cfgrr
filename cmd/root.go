package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/osamaadam/cfgrr/vconfig"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cfgrr [sub_command]",
	Short: `A one-hit solution for your configuration trouble`,
	Long: `cfgrr is a tool for managing config files inspired by GNU Stow.
Essentially, what cfgrr enables you to do is to centralize your config files, creating symlinks of them wherever necessary.
This enables the user to backup their config files to say Git, and restore the files easily.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute(version, tagdate string) error {
	if version != "" && tagdate != "" {
		rootCmd.SetVersionTemplate(fmt.Sprintf("cfgrr %s (published on %s)\n", version, tagdate))
		rootCmd.Version = version
	}

	if err := rootCmd.Execute(); err != nil {
		if tedious {
			fmt.Fprintf(os.Stderr, "ERROR: %+v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		}
		return err
	}

	return nil
}

func init() {
	cobra.OnInitialize(initConfig)
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	c := vconfig.GetConfig()

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", filepath.Join(homedir, ".cfgrr.yaml"), "config file")
	rootCmd.PersistentFlags().StringP("backup_dir", "d", c.BackupDir, "backup directory")
	rootCmd.PersistentFlags().StringSliceP("ignore_files", "i", []string{".cfgrrignore", ".gitignore"}, "ignore file")
	rootCmd.PersistentFlags().StringP("map_file", "m", c.MapFile, "map file")
	rootCmd.PersistentFlags().BoolVarP(&tedious, "tedious", "t", false, "print verbose errors")

	rootCmd.MarkFlagDirname("backup_dir")
	rootCmd.MarkFlagFilename("map_file", "yaml", "json")
	rootCmd.MarkFlagFilename("config", "yaml", "json")

	rootCmd.MarkFlagDirname("backup_dir")
	rootCmd.MarkFlagFilename("map_file", "yaml", "json")
	rootCmd.MarkFlagFilename("config", "yaml", "json")

	v := vconfig.GetViper()

	v.BindPFlag("backup_dir", rootCmd.PersistentFlags().Lookup("backup_dir"))
	v.BindPFlag("map_file", rootCmd.PersistentFlags().Lookup("map_file"))

	rootCmd.AddCommand(restoreCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(unsetCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(replicateCmd)
}

func initConfig() {
	if cfgFile != "" {
		// The user provided a custom config file.
		c := vconfig.GetConfig()
		c.SetConfigFile(cfgFile)
	}
}
