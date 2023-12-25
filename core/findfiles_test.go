package core

import (
	"os"
	"path/filepath"
	"testing"

	cf "github.com/osamaadam/cfgrr/configfile"
	"github.com/osamaadam/cfgrr/helpers"
	"github.com/osamaadam/cfgrr/ignorefile"
)

func TestCheckGlobsMatch(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		patterns []string
		out      bool
	}{
		{"empty string", "", []string{"*"}, false},
		{"empty string, no match", "", []string{"/nomatch", "./*", "~/*"}, false},
		{"empty patterns", "test", []string{}, false},
		{"no match", "/path/nomatch/file", []string{"/path/match/**/file"}, false},
		{"match", "/path/match/file", []string{"/path/match/**/file"}, true},
		{"multiple matches", "/path/match/file", []string{"/path/match/**/file", "/path/match/file"}, true},
		{"multiple matches, no match", "/path/nomatch/file", []string{"/path/match/**/file", "/path/match/file"}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			matches := CheckIfGlobsMatch(test.in, test.patterns...)

			if matches != test.out {
				t.Errorf("expected %v, got %v", test.out, matches)
			}
		})
	}
}

func TestFindFiles(t *testing.T) {
	backupDir := t.TempDir()
	_createFiles(
		backupDir,
		".vimrc",
		".bashrc",
		".zshrc",
		".config/nvim/init.vim",
		".config/nvim/coc-settings.json",
	)

	tests := []struct {
		name          string
		filesToIgnore []string
		patterns      []string
		expectedFiles []string
		expectErr     bool
	}{
		{"no patterns", []string{}, []string{}, []string{}, true},
		{"accept all, ignore none", []string{}, []string{"**/*"}, []string{
			".vimrc",
			".bashrc",
			".zshrc",
			".config/nvim/init.vim",
			".config/nvim/coc-settings.json",
		}, false},
		{"ignore specific patterns", []string{"**/*rc"}, []string{"**/*"}, []string{
			".config/nvim/init.vim",
			".config/nvim/coc-settings.json",
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			igContainer := _createIgnoreFile(t.TempDir(), tt.filesToIgnore...)

			files, err := FindFiles(backupDir, igContainer, tt.patterns...)
			if tt.expectErr && err != nil {
				return
			} else if tt.expectErr && err == nil {
				t.Fatalf("expected error, got nil")
			} else if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(files) != len(tt.expectedFiles) {
				t.Fatalf("expected %d files, got %d", len(tt.expectedFiles), len(files))
			}

			for _, file := range files {
				relativePath, err := filepath.Rel(backupDir, file.PathAbs())
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !helpers.Contains(tt.expectedFiles, relativePath) {
					t.Errorf("expected to find %s but didn't", relativePath)
				}
			}
		})
	}
}

func _createFiles(dir string, names ...string) []*cf.ConfigFile {
	numOfFiles := len(names)
	files := make([]*cf.ConfigFile, numOfFiles)

	for i, name := range names {
		filePath := filepath.Join(dir, name)
		fileDir := filepath.Dir(filePath)
		if err := os.MkdirAll(fileDir, 0755); err != nil {
			panic(err)
		}
		f, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}
		f.Close()
		files[i], _ = cf.NewConfigFile(f.Name())
	}

	return files
}

func _createIgnoreFile(backupDir string, patterns ...string) ignorefile.IIgnoresContainer {
	ignFile, err := ignorefile.InitDefaultIgnoreFile()
	if err != nil {
		panic(err)
	}

	if err := ignFile.WriteLines(patterns...); err != nil {
		panic(err)
	}

	ignContainer := ignorefile.NewIgnoresContainer(filepath.Base(ignFile.Path()))

	return ignContainer
}
