package cmd

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var resetCmd = &cobra.Command{
	Use:     "reset",
	Aliases: []string{"rs", "rst"},
	RunE:    reset,
	Args:    cobra.NoArgs,
	Short:   "Reset the configuration files to their original state",
	Long: `Reset the configuration files to their original state
This will remove the configuration file, and running the program will create a new one with the default values`,
}

func reset(cmd *cobra.Command, args []string) error {
	if err := viper.SafeWriteConfig(); err != nil {
		configPath := viper.ConfigFileUsed()
		if err := os.Remove(configPath); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
