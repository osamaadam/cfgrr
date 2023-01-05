package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/osamaadam/cfgrr/core"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup the configuration file",
	Args:  cobra.NoArgs,
	RunE:  runSetup,
}

func runSetup(cmd *cobra.Command, args []string) error {
	config := &core.Config{}

	if err := viper.Unmarshal(config); err != nil {
		return err
	}

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

	viper.Set("backup_dir", config.BackupDir)
	viper.Set("map_file", config.MapFile)
	viper.Set("ignore_file", config.IgnoreFile)

	if err := viper.WriteConfig(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
