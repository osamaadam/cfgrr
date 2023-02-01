package ignorefile

import (
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/osamaadam/cfgrr/vconfig"
)

func TestNewIgnoresContainer(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		out  []string
	}{
		// A default name is always provided from the viper config, keep that in mind.
		// So even if you pass no names, the default name will be added.
		{"no names", []string{}, []string{}},
		{"one name", []string{"yo"}, []string{}},
		{"multiple names", []string{"yo", "hi"}, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := vconfig.GetConfig()
			temp := t.TempDir()
			c.SetBackupDir(temp)
			// Adding the default file name to the arr of names.
			tt.in = append(tt.in, c.IgnoreFile)
			// Adding the backup_dir/name path for each name.
			for _, name := range tt.in {
				path := filepath.Join(c.BackupDir, name)
				curDirPath := filepath.Join("./", name)
				tt.out = append(tt.out, path, curDirPath)
			}

			ic := NewIgnoresContainer(tt.in...)
			paths := ic.Paths()
			sort.Strings(paths)
			sort.Strings(tt.out)
			if !reflect.DeepEqual(paths, tt.out) {
				t.Errorf("Expected %v, got %v", tt.out, paths)
			}
		})
	}
}
