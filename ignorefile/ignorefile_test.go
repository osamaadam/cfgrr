package ignorefile

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/osamaadam/cfgrr/helpers"
)

func TestInitIgnoreFile(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		if err := InitIgnoreFile(""); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("non-empty path", func(t *testing.T) {
		tempDir := t.TempDir()
		ignoreFilePath := filepath.Join(tempDir, ".cfgrrignore")

		if err := InitIgnoreFile(ignoreFilePath); err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		lines, err := helpers.ReadFileLines(ignoreFilePath)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		for i, line := range lines {
			if line != defaultIgnores[i] {
				t.Errorf("expected %v, got %v", defaultIgnores[i], line)
			}
		}
	})

	t.Run("non-empty path, file exists", func(t *testing.T) {
		tempDir := t.TempDir()
		ignoreFilePath := filepath.Join(tempDir, ".cfgrrignore")

		ignoredFiles := []string{"test", "hi", "yo"}

		if err := ioutil.WriteFile(ignoreFilePath, []byte(strings.Join(ignoredFiles, "\n")), 0644); err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if err := InitIgnoreFile(ignoreFilePath); err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		lines, err := helpers.ReadFileLines(ignoreFilePath)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		for i, line := range lines {
			if line != ignoredFiles[i] {
				t.Errorf("expected %v, got %v", defaultIgnores[i], line)
			}
		}
	})
}
