// Viper config manager.
package vconfig

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	BackupDir  string `mapstructure:"backup_dir"`
	MapFile    string `mapstructure:"map_file"`
	IgnoreFile string `mapstructure:"ignore_file"`
}

// Gets the config from the config file.
func GetConfig() *Config {
	var c Config

	if err := viper.Unmarshal(&c); err != nil {
		panic(err)
	}

	return &c
}

// Gets the full path of the map file.
func (c *Config) GetMapFilePath() string {
	return filepath.Join(c.BackupDir, c.MapFile)
}

// Gets the full path of the ignore file.
func (c *Config) GetIgnoreFilePath() string {
	return filepath.Join(c.BackupDir, c.IgnoreFile)
}

// Sets the backup directory.
// Does not save the config.
func (c *Config) SetBackupDir(path string) {
	path = filepath.Clean(path)
	viper.Set("backup_dir", path)
	c.BackupDir = path
}

// Sets the map file.
// Does not save the config.
func (c *Config) SetMapFile(name string) {
	viper.Set("map_file", name)
	c.MapFile = name
}

// Sets the ignore file.
// Does not save the config.
func (c *Config) SetIgnoreFile(name string) {
	viper.Set("ignore_file", name)
	c.IgnoreFile = name
}

// Sets a key and value to the config file.
func (c *Config) Set(key string, values ...string) error {
	viper.Set(key, values)

	switch key {
	case "backup_dir":
		c.SetBackupDir(values[0])
	case "map_file":
		c.SetMapFile(values[0])
	case "ignore_file":
		c.SetIgnoreFile(values[0])
	}

	if err := c.Save(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Sets all the config values for viper.
func (c *Config) setAll() {
	viper.Set("backup_dir", c.BackupDir)
	viper.Set("map_file", c.MapFile)
	viper.Set("ignore_file", c.IgnoreFile)
}

// Saves the current config to the config file.
func (c *Config) Save() error {
	c.setAll()
	if err := viper.WriteConfig(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Initializes the config.
// This should be called on the startup of the app.
func (c *Config) Init() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return errors.WithStack(err)
	}
	viper.AddConfigPath(homedir)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".cfgrr")

	userConfig, _ := os.UserConfigDir()
	if userConfig == "" {
		userConfig = filepath.Join(homedir, ".config")
	}

	defaultConfigDir := filepath.Join(userConfig, "cfgrr")

	c.SetBackupDir(defaultConfigDir)
	c.SetMapFile("cfgrrmap.yaml")
	c.SetIgnoreFile(".cfgrrignore")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; creating it.
			if err := viper.SafeWriteConfig(); err != nil {
				return errors.WithStack(err)
			}
		} else {
			// Config file was found but another error was produced
			return errors.WithStack(err)
		}
	}

	return nil
}
