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
	Browsable  bool   `mapstructure:"browsable"`
}

var v *viper.Viper
var vc Config

// Returns a pointer to the viper instance.
func GetViper() *viper.Viper {
	if v == nil {
		vc.init()
	}
	return v
}

// Gets the config from the config file.
func GetConfig() *Config {
	if v == nil {
		// Viper is not initialized.
		if err := vc.init(); err != nil {
			panic(err)
		}
		// Refreshing the values to read from Viper.
		if err := vc.refresh(); err != nil {
			panic(err)
		}
	}

	return &vc
}

// Refreshes the config struct.
func (c *Config) refresh() error {
	return v.Unmarshal(c)
}

// Sets the main config file to read from.
func (c *Config) SetConfigFile(file string) error {
	v.SetConfigFile(file)

	if err := c.refresh(); err != nil {
		return errors.WithStack(err)
	}

	return nil
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
	v.Set("backup_dir", path)
	c.BackupDir = path
}

// Sets the map file.
// Does not save the config.
func (c *Config) SetMapFile(name string) {
	v.Set("map_file", name)
	c.MapFile = name
}

// Sets the ignore file.
// Does not save the config.
func (c *Config) SetIgnoreFile(name string) {
	v.Set("ignore_file", name)
	c.IgnoreFile = name
}

func (c *Config) SetBrowsable(browsable bool) {
	viper.Set("browsable", browsable)
	c.Browsable = browsable
}

// Sets a key and value to the config file.
func (c *Config) Set(key string, values ...string) error {
	v.Set(key, values)

	switch key {
	case "backup_dir":
		c.SetBackupDir(values[0])
	case "map_file":
		c.SetMapFile(values[0])
	case "ignore_file":
		c.SetIgnoreFile(values[0])
	case "browsable":
		browsable := values[0] == "true"
		c.SetBrowsable(browsable)
	}

	if err := c.Save(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Sets all the config values for v.
func (c *Config) setAll() {
	v.Set("backup_dir", c.BackupDir)
	v.Set("map_file", c.MapFile)
	v.Set("ignore_file", c.IgnoreFile)
}

// Saves the current config to the config file.
func (c *Config) Save() error {
	c.setAll()
	if err := v.WriteConfig(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Initializes the config.
// This should be called on the startup of the app.
func (c *Config) init() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return errors.WithStack(err)
	}
	v = viper.New()
	v.AddConfigPath(homedir)
	v.SetConfigType("yaml")
	v.SetConfigName(".cfgrr")

	userConfig, _ := os.UserConfigDir()
	if userConfig == "" {
		userConfig = filepath.Join(homedir, ".config")
	}

	defaultConfigDir := filepath.Join(userConfig, "cfgrr")

	v.SetDefault("backup_dir", defaultConfigDir)
	v.SetDefault("map_file", "cfgrrmap.yaml")
	v.SetDefault("ignore_file", ".cfgrrignore")
	if err := v.ReadInConfig(); err != nil {
		if err := c.refresh(); err != nil {
			return errors.WithStack(err)
		}
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; creating it.
			c.SetBackupDir(defaultConfigDir)
			c.SetMapFile("cfgrrmap.yaml")
			c.SetIgnoreFile(".cfgrrignore")
			if err := v.SafeWriteConfig(); err != nil {
				return errors.WithStack(err)
			}
		} else {
			// Config file was found but another error was produced
			return errors.WithStack(err)
		}
	}

	return nil
}
