package cmd

import (
	"strings"

	"github.com/osamaadam/cfgrr/vconfig"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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
	vals := args[1:]

	config := vconfig.GetConfig()

	if err := config.Set(key, vals...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
