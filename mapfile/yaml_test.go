package mapfile

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	cf "github.com/osamaadam/cfgrr/configfile"
	"github.com/osamaadam/cfgrr/helpers"
	"github.com/spf13/viper"
)

func TestNewYamlFile(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{"normal path", "/some/path/config", "/some/path/config.yaml"},
		{"yaml extension", "/some/path/config.yaml", "/some/path/config.yaml"},
		{"yml extension", "/some/path/config.yml", "/some/path/config.yml"},
		{"weird extension", "/some/path/config.weird", "/some/path/config.weird.yaml"},
		{"hidden file", "/some/path/.config", "/some/path/.config.yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yf := NewYamlMapFile(tt.in)
			if yf.Path() != tt.out {
				t.Errorf("Expected %s, got %s", tt.out, yf.Path())
			}
		})
	}
}

func TestYamlMapFile_Parse(t *testing.T) {
	tests := []struct {
		name    string
		in      int
		wantErr bool
	}{
		{"no files", 0, false},
		{"one file", 1, false},
		{"69 files fuckit", 69, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			temp := t.TempDir()
			in, expected := _createConfigFiles(temp, tt.in)
			// Creation
			path := filepath.Join(temp, "test.yaml")
			yf := NewYamlMapFile(path)
			if err := yf.AddFiles(in...); err != nil && !tt.wantErr {
				t.Errorf("expected no error, got %s", err)
			}
			// Comparison
			parsed, err := yf.Parse()
			if err != nil && !tt.wantErr {
				t.Errorf("expected no error, got %s", err)
			}
			if !reflect.DeepEqual(parsed, expected) {
				t.Errorf("Expected %s, got %s", expected, parsed)
			}
		})
	}
}

func TestYamlMapFile_AddFiles(t *testing.T) {
	tests := []struct {
		name          string
		in            int
		existingFiles int
		wantErr       bool
	}{
		{"no file", 0, 0, false},
		{"multiple files", 7, 0, false},
		{"no files, existing files", 0, 7, false},
		{"multiple files, existing files", 7, 7, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp := t.TempDir()
			yf := _createYamlFileWithExistingFiles(temp, tt.existingFiles)

			in, expected := _createConfigFiles(temp, tt.in)
			preExisting, _ := yf.Parse()

			// Testing
			if err := yf.AddFiles(in...); err != nil && !tt.wantErr {
				t.Errorf("expected no error, got %s", err)
			}

			for k, v := range preExisting {
				expected[k] = v
			}

			if files, err := yf.Parse(); err != nil && !tt.wantErr {
				t.Errorf("expected no error, got %s", err)
			} else {
				if !reflect.DeepEqual(files, expected) {
					t.Errorf("Expected %s, got %s", expected, files)
				} else {
					t.Logf("Expected has %d keys, got has %d keys", len(expected), len(files))
				}
			}
		})
	}
}

func TestYamlMapFile_RemoveFiles(t *testing.T) {
	tests := []struct {
		name     string
		remove   int
		existing int
		out      int
		wantErr  bool
	}{
		{"no files", 0, 0, 0, false},
		{"no files, existing files", 0, 7, 7, false},
		{"multiple files, existing files", 7, 7, 0, false},
		{"multiple files but not all, existing files", 3, 7, 4, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp := t.TempDir()
			yf := _createYamlFileWithExistingFiles(temp, tt.existing)
			existing, _ := yf.Parse()
			toRemove := helpers.GetMapValues(existing)[:tt.remove]
			expected := existing

			for _, v := range toRemove {
				delete(expected, v.HashShort())
			}

			if err := yf.RemoveFiles(toRemove...); err != nil && !tt.wantErr {
				t.Errorf("expected no error, got %s", err)
			}

			parsed, _ := yf.Parse()

			if len(parsed) != tt.out {
				t.Errorf("Expected %d, got %d", tt.out, len(parsed))
			}

			if !reflect.DeepEqual(parsed, expected) {
				t.Errorf("Expected %s, got %s", expected, parsed)
			}
		})

	}
}

func TestYamlMapFile_Tidy(t *testing.T) {
	tests := []struct {
		name     string
		remove   int
		existing int
		out      int
		wantErr  bool
	}{
		{"no files", 0, 0, 0, false},
		{"no files, existing files", 0, 7, 7, false},
		{"remove multiple files, existing files", 7, 7, 0, false},
		{"remove multiple files but not all, existing files", 3, 7, 4, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp := t.TempDir()
			yf := _createYamlFileWithExistingFiles(temp, tt.existing)
			existing, _ := yf.Parse()
			expected := existing
			viper.Set("backup_dir", temp)
			i := 0
			for k, v := range expected {
				if i < tt.remove {
					delete(expected, k)
				} else {
					v.Backup()
				}
				i++
			}

			if err := yf.Tidy(); err != nil && !tt.wantErr {
				t.Errorf("expected no error, got %s", err)
			}

			parsed, _ := yf.Parse()
			if len(parsed) != tt.out {
				t.Errorf("Expected %d, got %d", tt.out, len(parsed))
			}

			if !reflect.DeepEqual(parsed, expected) {
				t.Errorf("Expected %s, got %s", expected, parsed)
			}

		})
	}
}

func _createConfigFiles(dir string, num int) (in []*cf.ConfigFile, expected map[string]*cf.ConfigFile) {
	expected = make(map[string]*cf.ConfigFile)
	in = make([]*cf.ConfigFile, num)
	for i := 0; i < num; i++ {
		tempFile, _ := ioutil.TempFile(dir, "")
		file, _ := cf.NewConfigFile(tempFile.Name())
		tempFile.Close()
		expected[file.HashShort()] = file
		in[i] = file
	}

	return
}

func _createYamlFileWithExistingFiles(dir string, num int) (yf *YamlMapFile) {
	in, _ := _createConfigFiles(dir, num)
	yf = NewYamlMapFile(filepath.Join(dir, "test.yaml"))

	yf.AddFiles(in...)

	return
}
