package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/osamaadam/cfgrr/vconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "cfgrr [sub_command]",
	Short: `A one-hit solution for your configuration trouble`,
}

func Execute(version, tagdate string) error {
	if version != "" && tagdate != "" {
		rootCmd.SetVersionTemplate(fmt.Sprintf("cfgrr %s (published on %s)\n", version, tagdate))
		rootCmd.Version = version
	}
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", filepath.Join(homedir, ".cfgrr.yaml"), "config file")
	rootCmd.PersistentFlags().StringP("backup_dir", "d", "", "backup directory (default $HOME/.config/cfgrr)")
	rootCmd.PersistentFlags().StringSliceP("ignore_files", "i", []string{".cfgrrignore", ".gitignore"}, "ignore file")
	rootCmd.PersistentFlags().StringP("map_file", "m", "", "map file (default cfgrrmap.yaml)")

	viper.BindPFlag("backup_dir", rootCmd.PersistentFlags().Lookup("backup_dir"))
	viper.BindPFlag("map_file", rootCmd.PersistentFlags().Lookup("map_file"))

	rootCmd.AddCommand(restoreCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(unsetCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(deleteCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		c := vconfig.GetConfig()
		if err := c.Init(); err != nil {
			panic(err)
		}
	}
}
