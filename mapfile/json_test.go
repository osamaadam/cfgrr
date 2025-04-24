package mapfile

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/osamaadam/cfgrr/helpers"
	"github.com/osamaadam/cfgrr/vconfig"
)

func TestNewJsonFile(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{"normal path", "/some/path/config", "/some/path/config.json"},
		{"json extension", "/some/path/config.json", "/some/path/config.json"},
		{"weird extension", "/some/path/config.weird", "/some/path/config.weird.json"},
		{"hidden file", "/some/path/.config", "/some/path/.config.json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jf := NewJsonMapFile(tt.in)
			if jf.Path() != tt.out {
				t.Errorf("Expected %s, got %s", tt.out, jf.Path())
			}
		})
	}
}

func TestJsonMapFile_Parse(t *testing.T) {
	tests := []struct {
		name    string
		in      int
		wantErr bool
	}{
		{"no files", 0, false},
		{"one file", 1, false},
		{"multiple files", 69, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			temp := t.TempDir()
			in, expected := _createConfigFiles(temp, tt.in)
			// Creation
			path := filepath.Join(temp, "test.json")
			jf := NewJsonMapFile(path)
			if err := jf.AddFiles(in...); err != nil && !tt.wantErr {
				t.Errorf("expected no error, got %s", err)
			}
			// Comparison
			parsed, err := jf.Parse()
			if err != nil && !tt.wantErr {
				t.Errorf("expected no error, got %s", err)
			}
			if !reflect.DeepEqual(parsed, expected) {
				t.Errorf("Expected %s, got %s", expected, parsed)
			}
		})
	}
}

func TestJsonMapFile_AddFiles(t *testing.T) {
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
			jf := _createJsonFileWithExistingFiles(temp, tt.existingFiles)

			in, expected := _createConfigFiles(temp, tt.in)
			preExisting, _ := jf.Parse()

			// Testing
			if err := jf.AddFiles(in...); err != nil && !tt.wantErr {
				t.Errorf("expected no error, got %s", err)
			}

			for k, v := range preExisting {
				expected[k] = v
			}

			if files, err := jf.Parse(); err != nil && !tt.wantErr {
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

func TestJsonMapFile_RemoveFiles(t *testing.T) {
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
			jf := _createJsonFileWithExistingFiles(temp, tt.existing)
			existing, _ := jf.Parse()
			toRemove := helpers.GetMapValues(existing)[:tt.remove]
			expected := existing

			for _, v := range toRemove {
				delete(expected, v.HashShort())
			}

			if err := jf.RemoveFiles(toRemove...); err != nil && !tt.wantErr {
				t.Errorf("expected no error, got %s", err)
			}

			parsed, _ := jf.Parse()

			if len(parsed) != tt.out {
				t.Errorf("Expected %d, got %d", tt.out, len(parsed))
			}

			if !reflect.DeepEqual(parsed, expected) {
				t.Errorf("Expected %s, got %s", expected, parsed)
			}
		})

	}
}

func TestJsonMapFile_Tidy(t *testing.T) {
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
			jf := _createJsonFileWithExistingFiles(temp, tt.existing)
			existing, _ := jf.Parse()
			expected := existing
			c := vconfig.GetConfig()
			c.SetBackupDir(temp)
			i := 0
			for k, v := range expected {
				if i < tt.remove {
					delete(expected, k)
				} else {
					v.Backup()
				}
				i++
			}

			if err := jf.Tidy(); err != nil && !tt.wantErr {
				t.Errorf("expected no error, got %s", err)
			}

			parsed, _ := jf.Parse()
			if len(parsed) != tt.out {
				t.Errorf("Expected %d, got %d", tt.out, len(parsed))
			}

			if !reflect.DeepEqual(parsed, expected) {
				t.Errorf("Expected %s, got %s", expected, parsed)
			}

		})
	}
}

func _createJsonFileWithExistingFiles(dir string, num int) (jf *JsonMapFile) {
	in, _ := _createConfigFiles(dir, num)
	jf = NewJsonMapFile(filepath.Join(dir, "test.json"))

	jf.AddFiles(in...)

	return
}
