package ignorefile

import (
	"reflect"
	"sort"
	"testing"

	"github.com/osamaadam/cfgrr/helpers"
	"github.com/osamaadam/cfgrr/vconfig"
)

func TestInitDefaultIgnoreFile(t *testing.T) {
	t.Run("creates the file properly", func(t *testing.T) {
		c := vconfig.GetConfig()
		c.SetBackupDir(t.TempDir())
		c.SetIgnoreFile(".cfgrrignore")

		ignFile, err := InitDefaultIgnoreFile()
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if !helpers.CheckFileExists(ignFile.Path()) {
			t.Fatalf("expected %s to exist", ignFile.Path())
		}

		lines, err := ignFile.ReadLines()
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if !reflect.DeepEqual(lines, defaultIgnores) {
			t.Fatalf("expected %v, got %v", defaultIgnores, lines)
		}
	})
}

// Implicitely tests Read, Write and Append.
func TestIgnoreFile(t *testing.T) {
	tests := []struct {
		name     string
		in       []string
		existing []string
		out      []string
		wantErr  bool
	}{
		{"no lines, no existing", []string{}, []string{}, []string{}, false},
		{"no lines, multiple existing", []string{}, []string{"hi", "yo", "bye"}, []string{"hi", "yo", "bye"}, false},
		{"multiple lines, no existing", []string{"hi", "yo", "bye"}, []string{}, []string{"hi", "yo", "bye"}, false},
		{"multiple lines, multiple existing", []string{"hi", "yo", "bye"}, []string{"hi2", "yo2", "bye2"}, []string{"hi2", "yo2", "bye2", "hi", "yo", "bye"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := vconfig.GetConfig()
			c.SetBackupDir(t.TempDir())
			c.SetIgnoreFile(".cfgrrignore")
			i := NewIgnoreFile(c.GetIgnoreFilePath())

			if err := i.WriteLines(tt.existing...); err != nil {
				t.Errorf("expected nil, got %v", err)
			}
			if err := i.WriteLines(tt.in...); err != nil {
				t.Errorf("expected nil, got %v", err)
			}

			lines, err := i.ReadLines()
			if err != nil && !tt.wantErr {
				t.Errorf("expected nil, got %v", err)
			}

			sort.Strings(lines)
			sort.Strings(tt.out)

			if !reflect.DeepEqual(lines, tt.out) {
				// Apparently, [] != [] is true in Go, it's like Javascript all over again.
				if len(lines) != 0 || len(tt.out) != 0 {
					t.Errorf("expected %v, got %v", tt.out, lines)
				}
			}
		})
	}
}
