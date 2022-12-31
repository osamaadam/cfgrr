package configfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal("couldn't get the user's homedir")
	}
	tests := []struct {
		filePath string
		want     *ConfigFile
		wantErr  bool
	}{
		{
			filepath.Join(homedir, "test/probably/an/actual/.path"),
			&ConfigFile{
				Path: "test/probably/an/actual/.path",
				Name: ".path",
			},
			false,
		},
		{
			"/wrong/abs/.path",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		file, err := InitFile(tt.filePath)
		if err != nil && !tt.wantErr {
			t.Errorf("got an error %q", err)
		}

		if tt.want != nil {
			if file.Name != tt.want.Name {
				t.Errorf("got %s, was expecting %s", file.Name, tt.want.Name)
			}

			if file.Path != tt.want.Path {
				t.Errorf("got %s, was expecting %s", file.Path, tt.want.Path)
			}
		}
	}
}
