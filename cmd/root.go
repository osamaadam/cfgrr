package cmd

import (
	cf "github.com/osamaadam/gocfgr/configfile"
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

	if err := cf.RestoreConfig(root); err != nil {
		return errors.WithStack(err)
	}

	// configFiles, err := util.FindFiles(root, "testdata/.gocfgrignore", ".*")

	// if err != nil {
	// 	return err
	// }

	// selectedFiles, err := prompt.PromptForFileSelection(configFiles)

	// if err != nil {
	// 	return errors.WithStack(err)
	// }

	// copyDir, err := filepath.Abs("./testdata/lab")
	// if err != nil {
	// 	return errors.WithStack(err)
	// }

	// if err := cf.CopyAndReplaceFiles(copyDir, selectedFiles...); err != nil {
	// 	return errors.WithStack(err)
	// }

	return nil
}

func Execute() error {
	return rootCmd.Execute()
}
