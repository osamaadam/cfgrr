package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/osamaadam/cfgrr/cmd"
	"github.com/spf13/viper"
)

type Config struct {
	BackupDir  string `mapstructure:"backup_dir"`
	MapFile    string `mapstructure:"map_file"`
	IgnoreFile string `mapstructure:"ignore_file"`
}

var (
	version   string
	builddate string
)

func main() {
	if err := cmd.Execute(version, builddate); err != nil {
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

	viper.SetDefault("backup_dir", defaultConfigDir)
	viper.SetDefault("map_file", "cfgrrmap.yaml")
	viper.SetDefault("ignore_file", ".cfgrrignore")

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
