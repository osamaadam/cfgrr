package ignorefile

import (
	"os"
	"path/filepath"
	"strings"

	cf "github.com/osamaadam/cfgrr/configfile"
	"github.com/pkg/errors"
)

func InitIgnoreFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := cf.EnsureDirExists(filepath.Dir(path)); err != nil {
			return errors.WithStack(err)
		}
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := file.Write([]byte(strings.Join(defaultIgnores, "\n"))); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
