package cmd

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setCmd = &cobra.Command{
	Use:     "set [key] [value]",
	Aliases: []string{"s"},
	Short:   "Set the value of a configuration variable",
	Args:    cobra.ExactArgs(2),
	RunE:    runSet,
	Example: strings.Join([]string{
		`cfgrr set backup_dir /path/to/backup/dir`,
		`cfgrr s map_file cfgrrmap.yaml`,
		`cfgrr s ignore_file .cfgrrignore`,
	}, "\n"),
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
