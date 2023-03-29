package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/osamaadam/cfgrr/vconfig"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var unsetCmd = &cobra.Command{
	Use:     "unset [key]",
	Aliases: []string{"u"},
	Args:    cobra.MinimumNArgs(1),
	RunE:    runUnset,
	Example: strings.Join([]string{
		`cfgrr unset backup_dir`,
		`cfgrr u map_file`,
	}, "\n"),
	Short: "Unset the value of a configuration variable",
	Long: `Unset the value of a configuration variable. cfgrr would use a default value if a variable is unset.
Good idea to run this if the user thinks they've messed some application variable. A better idea would be to run 'setup' ('cfgrr setup --help' for more info).`,
}

func runUnset(cmd *cobra.Command, args []string) error {
	if err := unset(args...); err != nil {
		errors.WithStack(err)
	}

	return nil
}

// Unsets the given variables from the config file.
// Copied from https://github.com/spf13/viper/issues/632#issuecomment-869668629
func unset(vars ...string) error {
	vc := vconfig.GetViper()
	cfg := vc.AllSettings()
	vals := cfg

	for _, v := range vars {
		parts := strings.Split(v, ".")
		for i, k := range parts {
			v, ok := vals[k]
			if !ok {
				// Doesn't exist no action needed
				break
			}

			switch len(parts) {
			case i + 1:
				// Last part so delete.
				delete(vals, k)
			default:
				m, ok := v.(map[string]interface{})
				if !ok {
					return fmt.Errorf("unsupported type: %T for %q", v, strings.Join(parts[0:i], "."))
				}
				vals = m
			}
		}
	}

	b, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	if err = vc.ReadConfig(bytes.NewReader(b)); err != nil {
		return err
	}

	return vc.WriteConfig()
}
