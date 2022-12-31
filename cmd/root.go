package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/osamaadam/gocfgr/configfile"
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

	configFiles, err := util.FindFiles(root, "testdata/.gocfgrignore", ".*")

	filesMap := make(map[string]*configfile.ConfigFile)

	for _, file := range configFiles {
		filesMap[file.String()] = file
	}

	if err != nil {
		return err
	}

	filteredFiles := []string{}

	prompt := &survey.MultiSelect{
		Message:  "Which files would you like to track?",
		Options:  configfile.ArrToString(configFiles),
		PageSize: 15,
	}

	if err := survey.AskOne(prompt, &filteredFiles, survey.WithKeepFilter(true)); err != nil {
		return err
	}

	selectedFiles := make([]*configfile.ConfigFile, 0, len(filteredFiles))

	for _, file := range filteredFiles {
		selectedFiles = append(selectedFiles, filesMap[file])
	}

	fmt.Println(selectedFiles)

	return nil
}

func Execute() error {
	return rootCmd.Execute()
}
