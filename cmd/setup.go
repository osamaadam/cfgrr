package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/osamaadam/cfgrr/vconfig"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Args:  cobra.NoArgs,
	RunE:  runSetup,
	Short: "Setup the configuration file",
	Long: `Setup the configuration file.
Prompts the user to set their application variables. It is a good idea to run this command after updating cfgrr.`,
}

func runSetup(cmd *cobra.Command, args []string) error {
	config := vconfig.GetConfig()

	questions := []*survey.Question{
		{
			Name: "BackupDir",
			Prompt: &survey.Input{
				Message: "Backup directory: ",
				Default: config.BackupDir,
			},
		},
		{
			Name: "MapFile",
			Prompt: &survey.Input{
				Message: "Map file: ",
				Default: config.MapFile,
			},
		},
		{
			Name: "IgnoreFile",
			Prompt: &survey.Input{
				Message: "Ignore file: ",
				Default: config.IgnoreFile,
			},
		},
	}

	if err := survey.Ask(questions, config); err != nil {
		return errors.WithStack(err)
	}

	if err := config.Save(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
