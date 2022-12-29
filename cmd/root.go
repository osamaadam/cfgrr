package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/osamaadam/gocfgr/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gocfgr",
	Short: `A one-hit solution to your configuration trouble`,
	RunE:  runRoot,
}

func runRoot(cmd *cobra.Command, args []string) error {
	root := args[0]

	uniqueFiles, err := util.FindFiles(root, "testdata/.gocfgrignore", ".*")

	if err != nil {
		return err
	}

	filteredFiles := []string{}

	prompt := &survey.MultiSelect{
		Message:  "Which files would you like to track?",
		Options:  uniqueFiles,
		Default:  uniqueFiles,
		PageSize: 15,
	}

	if err := survey.AskOne(prompt, &filteredFiles, survey.WithKeepFilter(true)); err != nil {
		return err
	}

	return nil
}

func Execute() error {
	return rootCmd.Execute()
}
