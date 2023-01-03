package configfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitFile(t *testing.T) {
	t.Setenv("HOME", "/home/user")
	homedir, _ := os.UserHomeDir()

	tests := []struct {
		name string
		in   string
		out  *ConfigFile
		err  bool
	}{
		{"empty path", "", nil, true},
		{"valid path", filepath.Join(homedir, "path/to/file"), &ConfigFile{Path: "path/to/file"}, false},
		{"clean path", filepath.Join(homedir, "path/../path/.config"), &ConfigFile{Path: "path/.config"}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file, err := InitFile(test.in)
			if err != nil && !test.err {
				t.Errorf("got error: %v", err)
			}

			t.Log(file)
			if file != nil && file.Path != test.out.Path {
				t.Errorf("got path: %v, want: %v", file.Path, test.out.Path)
			}
		})
	}
}
