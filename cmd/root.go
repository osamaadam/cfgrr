package cmd

import (
	"fmt"

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

	testFilePath := "./testdata/config.yaml"

	if err := cf.CreateYamlFile(testFilePath, selectedFiles...); err != nil {
		return errors.WithStack(err)
	}

	readFiles, err := cf.ReadYamlFile(testFilePath)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, v := range readFiles {
		fmt.Println(v)
	}

	return nil
}

func Execute() error {
	return rootCmd.Execute()
}
