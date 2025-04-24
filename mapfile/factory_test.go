package mapfile

import (
	"path/filepath"
	"testing"

	"github.com/osamaadam/cfgrr/vconfig"
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
		{"json format", "/some/path/config.json", ".json", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := vconfig.GetConfig()
			c.SetBackupDir(t.TempDir())
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
