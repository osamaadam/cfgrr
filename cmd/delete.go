package cmd

import (
	"strings"

	cf "github.com/osamaadam/cfgrr/configfile"
	"github.com/osamaadam/cfgrr/core"
	"github.com/osamaadam/cfgrr/helpers"
	"github.com/osamaadam/cfgrr/mapfile"
	"github.com/osamaadam/cfgrr/prompt"
	"github.com/osamaadam/cfgrr/vconfig"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete [...paths]",
	Aliases: []string{"d", "del"},
	RunE:    deleteRun,
	Example: strings.Join([]string{
		"cfgrr delete",
		"cfgrr delete ~/.vimrc",
		"cfgrr delete ~/.vimrc ~/.zshrc",
		"cfgrr delete -r ~/.vimrc",
	}, "\n"),
	Short: "Delete the configuration files from the backup directory",
	Long:  `Delete the configuration files from the backup directory, also replacing symlinks with the original target if used with the --replace flag.`,
}

func deleteRun(cmd *cobra.Command, args []string) (err error) {
	files := make([]*cf.ConfigFile, 0)

	for _, path := range args {
		file, err := cf.NewConfigFile(path)
		if err != nil {
			return err
		}
		files = append(files, file)
	}

	config := vconfig.GetConfig()

	if len(files) == 0 {
		// User didn't specify any files, so we'll prompt them to select some
		// from the map file.
		mapFile := mapfile.NewMapFile(config.GetMapFilePath())
		m, err := mapFile.Parse()
		if err != nil {
			return errors.WithStack(err)
		}

		files = helpers.GetMapValues(m)

		files, err = prompt.PromptForFileSelection(files, "Select the files to delete: ")
		if err != nil {
			return errors.WithStack(err)
		}
	}

	if len(files) == 0 {
		return nil
	}

	if err := core.DeleteFiles(replace, files...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func init() {
	deleteCmd.Flags().BoolVarP(&replace, "replace", "r", false, "replace the symlinks with the original target")
}
