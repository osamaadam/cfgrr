package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	cf "github.com/osamaadam/cfgrr/configfile"
	"github.com/osamaadam/cfgrr/helpers"
	"github.com/osamaadam/cfgrr/mapfile"
	"github.com/osamaadam/cfgrr/vconfig"
)

func TestBackupCmd(t *testing.T) {
	dummyFiles := []string{
		".vimrc",
		".bashrc",
		".zshrc",
		".config/nvim/init.vim",
		".config/nvim/coc-settings.json",
	}

	backupDir := t.TempDir()

	tests := []struct {
		name          string
		args          []string
		expectError   bool
		expectedFiles []string
	}{
		{"no args", []string{}, true, nil},
		{"backup default files (dotfiles)", []string{backupDir, "-a"}, false, []string{".vimrc", ".bashrc", ".zshrc"}},
		{"backup only one file", []string{filepath.Join(backupDir, ".vimrc"), "-a"}, false, []string{".vimrc"}},
		{"backup nothing if no files match pattern", []string{backupDir, "-p", "**/*.yaml", "-a"}, false, []string{}},
		{"backup only files matching pattern", []string{backupDir, "-p", "**/*.json", "-a"}, false,
			[]string{".config/nvim/coc-settings.json"}},
		{"backup only files matching patterns", []string{backupDir, "-p", "**/*.json", "-p", "**/*.vim", "-a"}, false,
			[]string{".config/nvim/coc-settings.json", ".config/nvim/init.vim"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tCmd := rootCmd
			// Delete the backup dir if exists
			if err := os.RemoveAll(backupDir); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			// Creating files to backup
			_createFilesToBackup(
				backupDir,
				dummyFiles...,
			)
			// initialize a temp backup dir which is deleted after the test
			backupToDir := t.TempDir()
			c := vconfig.GetConfig()
			c.SetBackupDir(backupToDir)

			// cfgrr backup [args...]
			args := append([]string{"backup"}, tt.args...)
			t.Logf("cmd: cfgrr %v", strings.Join(args, " "))
			// Silence the output
			tCmd.SetArgs(args)

			err := tCmd.Execute()
			if tt.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			} else if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check if the files were backed up
			mf := mapfile.NewMapFile(c.GetMapFilePath())
			files, err := mf.Parse()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			backupFiles := []string{}
			for _, file := range files {
				backupFiles = append(backupFiles, file.PathAbs())
			}

			for _, expectedFile := range tt.expectedFiles {
				absPath := filepath.Join(backupDir, expectedFile)
				if !helpers.Contains(backupFiles, absPath) {
					t.Errorf("expected file %s to be backed up, but it wasn't", absPath)
				}
			}

			t.Logf("backed up files: %v", backupFiles)

			if len(tt.expectedFiles) != len(backupFiles) {
				t.Errorf("expected %d file(s) to be backed up, got %d", len(tt.expectedFiles), len(backupFiles))
			}
		})
	}
}

func _createFilesToBackup(dir string, names ...string) []*cf.ConfigFile {
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
