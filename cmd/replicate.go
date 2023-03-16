package cmd

import (
	"fmt"
	"strings"

	"github.com/osamaadam/cfgrr/core"
	"github.com/osamaadam/cfgrr/helpers"
	"github.com/osamaadam/cfgrr/mapfile"
	"github.com/osamaadam/cfgrr/prompt"
	"github.com/osamaadam/cfgrr/vconfig"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var replicateCmd = &cobra.Command{
	Use:     "replicate [root_dir]",
	Short:   "Creates a replica of the configuration files to the root dir",
	Aliases: []string{"rep", "browse"},
	Args:    cobra.MaximumNArgs(1),
	Example: strings.Join([]string{
		`cfgrr replicate ~/browsable/`,
		`cfgrr replicate`,
	}, "\n"),
	RunE: runReplicate,
}

func runReplicate(cmd *cobra.Command, args []string) error {
	baseDir := ""
	if len(args) > 0 {
		baseDir = args[0]
	}
	config := vconfig.GetConfig()

	mapFile := mapfile.NewMapFile(config.GetMapFilePath())

	m, err := mapFile.Parse()
	if err != nil {
		return errors.WithStack(err)
	}

	files := helpers.GetMapValues(m)

	if baseDir == "" {
		baseDir = "home"
	}

	if !all {
		files, err = prompt.PromptForFileSelection(files, "Select the files to replicate: ")
		if err != nil {
			return errors.WithStack(err)
		}
	}

	if len(files) == 0 {
		fmt.Println("No files selected, terminating...")
		return nil
	}

	if err := core.MakeFilesBrowsable(baseDir, files...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func init() {
	replicateCmd.Flags().BoolVarP(&all, "all", "a", false, "replicate all files in the backup directory (skip prompt)")
}
