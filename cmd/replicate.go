package cmd

import (
	"fmt"
	"os"
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
	Aliases: []string{"rep", "browse"},
	Args:    cobra.MaximumNArgs(1),
	Example: strings.Join([]string{
		`cfgrr replicate`,
		`cfgrr replicate --all`,
		`cfgrr replicate -a`,
		`cfgrr replicate -a --clean`,
		`cfgrr replicate ~/browsable/`,
		`cfgrr replicate ~/browsable/ -a`,
	}, "\n"),
	RunE:  runReplicate,
	Short: "Creates a replica of the configuration files to root_dir. If the file is already browsable, updates the browsable replica",
	Long: `Creates a replica of the configuration files to root_dir. If the file is already browsable, updates the browsable replica.
This should be run if the user intends to put their configuration on display on any platform. By default cfgrr saves the backed up files as hashes.
This is to avoid GNU stow's method of replicating the entire file's path structure, and instead relies on a map file to keep track which file should be restored where.
If the user intends to keep the files private, it wouldn't make sense for them to replicate them. However, they may find it convenient for readability and syntax highlighting.`,
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

	if clean {
		if err := os.RemoveAll(baseDir); err != nil {
			return errors.WithStack(err)
		}
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
	replicateCmd.Flags().BoolVar(&clean, "clean", false, "remove all files in the replica directory before replicating")
}
