package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"s"},
	Short:   "Set the value of a configuration variable",
	Args:    cobra.ExactArgs(2),
	RunE:    runSet,
}

func runSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	val := args[1]

	viper.Set(key, val)

	if err := viper.WriteConfig(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
