package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "cfgrr",
	Short: `A one-hit solution for your configuration trouble`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.cfgrr.yaml)")
	rootCmd.PersistentFlags().StringP("backup-dir", "d", "", "backup directory (default is $HOME/.config/cfgrr)")
	rootCmd.PersistentFlags().StringP("ignore-file", "i", "", "ignore file (default is .cfgrrignore)")
	rootCmd.PersistentFlags().StringP("map-file", "m", "", "map file (default is cfgrrmap.yaml)")

	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("backup-dir", rootCmd.PersistentFlags().Lookup("backup-dir"))
	viper.BindPFlag("map-file", rootCmd.PersistentFlags().Lookup("map-file"))

	rootCmd.AddCommand(restoreCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(unsetCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		homedir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		viper.AddConfigPath(homedir)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cfgrr")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; creating it.
			// TODO: Flag arguments shouldn't be hardcoded into the config file.
			if err := viper.SafeWriteConfig(); err != nil {
				panic(err)
			}
		} else {
			// Config file was found but another error was produced
			panic(err)
		}
	}
}
