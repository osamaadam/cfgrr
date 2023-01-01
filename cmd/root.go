package cmd

import (
	"path/filepath"

	cf "github.com/osamaadam/gocfgr/configfile"
	"github.com/osamaadam/gocfgr/prompt"
	"github.com/osamaadam/gocfgr/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gocfgr",
	Short: `A one-hit solution to your configuration trouble`,
	RunE:  runRoot,
}

func runRoot(cmd *cobra.Command, args []string) error {
	root := args[0]

	configFiles, err := util.FindFiles(root, "testdata/.gocfgrignore", ".*")

	if err != nil {
		return err
	}

	selectedFiles, err := prompt.PromptForFileSelection(configFiles)

	if err != nil {
		return errors.WithStack(err)
	}

	copyDir, err := filepath.Abs("./testdata/lab")
	if err != nil {
		return errors.WithStack(err)
	}

	if err := cf.CopyFiles(copyDir, selectedFiles...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func Execute() error {
	return rootCmd.Execute()
}
