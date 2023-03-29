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
	Args:    cobra.ExactArgs(2),
	RunE:    runSet,
	Example: strings.Join([]string{
		`cfgrr set backup_dir /path/to/backup/dir`,
		`cfgrr s map_file cfgrrmap.yaml`,
		`cfgrr s ignore_file .cfgrrignore`,
	}, "\n"),
	Short: "Set the value of a configuration variable",
	Long: `Set the value of a configuration variable.
This is the prefered method for the user to set the application variables instead of manually editing the '.cfgrr.yaml' file at their home directory.`,
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
