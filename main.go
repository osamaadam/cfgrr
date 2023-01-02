package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/osamaadam/cfgrr/cmd"
	"github.com/spf13/viper"
)

type Config struct {
	Config     string `mapstructure:"config"`
	BackupDir  string `mapstructure:"backup-dir"`
	MapFile    string `mapstructure:"map-file"`
	IgnoreFile string `mapstructure:"ignore-file"`
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Print(err)
	}
}

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	viper.AddConfigPath(homedir)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".cfgrr")

	userConfig, _ := os.UserConfigDir()
	if userConfig == "" {
		userConfig = filepath.Join(homedir, ".config")
	}

	defaultConfigDir := filepath.Join(userConfig, "cfgrr")

	viper.SetDefault("config", filepath.Join(homedir, ".cfgrr.yaml"))
	viper.SetDefault("backup-dir", defaultConfigDir)
	viper.SetDefault("map-file", "cfgrrmap.yaml")
	viper.SetDefault("ignore-file", ".cfgrrignore")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; creating it.
			if err := viper.SafeWriteConfig(); err != nil {
				panic(err)
			}
		} else {
			// Config file was found but another error was produced
			panic(err)
		}
	}
}
