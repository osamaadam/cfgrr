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
		// Dear future me:
		// This test is a little unorthodox.
		// While the out array seems to be empty,
		// The default behaviour of the function is to look for
		// the files matching the name in both the current working
		// directory and in the backup_dir from the viper config.
		// So if the the name is `.cfgrrignore`, it looks looks for it
		// in both `.../backup_dir/.cfgrrignore`, and `./.cfgrrignore`.
		{"no names", []string{}, []string{}},
		{"one name", []string{"yo"}, []string{}},
		{"multiple names", []string{"yo", "hi"}, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := vconfig.GetConfig()
			temp := t.TempDir()
			c.SetBackupDir(temp)
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
