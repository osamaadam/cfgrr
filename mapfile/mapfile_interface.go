package mapfile

import (
	"fmt"
	cf "github.com/osamaadam/cfgrr/configfile"
)

type IMapFile interface {
	fmt.Stringer
	Parse() (map[string]*cf.ConfigFile, error)
	AddFiles(files ...*cf.ConfigFile) error
	RemoveFiles(files ...*cf.ConfigFile) error
	Tidy() error
}
