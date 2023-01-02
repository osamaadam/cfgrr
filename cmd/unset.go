package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var unsetCmd = &cobra.Command{
	Use:     "unset",
	Aliases: []string{"u"},
	Short:   "Unset the value of a configuration variable",
	Args:    cobra.MinimumNArgs(1),
	RunE:    runUnset,
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
	cfg := viper.AllSettings()
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

	if err = viper.ReadConfig(bytes.NewReader(b)); err != nil {
		return err
	}

	return viper.WriteConfig()
}
