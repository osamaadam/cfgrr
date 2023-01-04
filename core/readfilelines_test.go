package core

import (
	"os"
	"strings"
	"testing"
)

func TestReadFileLines(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		out  []string
	}{
		{"empty", []string{}, []string{}},
		{"one line", []string{"one"}, []string{"one"}},
		{"two lines", []string{"one", "two"}, []string{"one", "two"}},
		{"multiple lines", []string{"one", "two", "three"}, []string{"one", "two", "three"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempFile, err := os.CreateTemp(t.TempDir(), "test-*")
			if err != nil {
				t.Errorf("error creating temp file: %v", err)
			}

			if _, err := tempFile.WriteString(strings.Join(test.in, "\n")); err != nil {
				t.Errorf("error writing to temp file: %v", err)
			}

			filePath := tempFile.Name()
			tempFile.Close()

			lines, err := ReadFileLines(filePath)
			if err != nil {
				t.Errorf("error reading file: %v", err)
			}

			for i, line := range lines {
				if line != test.out[i] {
					t.Errorf("expected %s, got %s", test.out[i], line)
				}
			}
		})
	}
}
