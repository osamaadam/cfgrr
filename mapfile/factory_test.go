package mapfile

import (
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestNewMapFile(t *testing.T) {
	tests := []struct {
		name      string
		in        string
		outExt    string
		expectNil bool
	}{
		{"no path", "", ".yaml", false},
		{"normal path", "/some/path/config", ".yaml", false},
		{"already yaml", "/some/path/config.yaml", ".yaml", false},
		{"already yml", "/some/path/config.yml", ".yml", false},
		{"UNIMPLEMENTED: json", "/some/path/config.json", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("backup_dir", t.TempDir())
			mf := NewMapFile(tt.in)
			if mf == nil && !tt.expectNil {
				t.Errorf("expected non-nil mapfile, got nil")
			} else if mf == nil && tt.expectNil {
				return
			}
			ext := filepath.Ext(mf.Path())

			if ext != tt.outExt {
				t.Errorf("expected %s, got %s", tt.outExt, ext)
			}
		})
	}
}
