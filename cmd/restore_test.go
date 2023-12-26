package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/osamaadam/cfgrr/core"
	"github.com/osamaadam/cfgrr/helpers"
	"github.com/osamaadam/cfgrr/vconfig"
)

func TestRestoreCmd(t *testing.T) {
	tests := []struct {
		name                  string
		args                  []string
		filesToBackup         []string
		expectedRestoredFiles []string
	}{
		{"no args", []string{}, _dummyTestFiles, nil},
		{"restore all files", []string{"-a"}, _dummyTestFiles, _dummyTestFiles},
		{"no files to restore", []string{"-a"}, nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgDir, _ := _restoreSetup(t, tt.filesToBackup...)

			args := append([]string{"restore"}, tt.args...)

			t.Logf("$ cfgrr %v", strings.Join(args, " "))
			tCmd := rootCmd
			tCmd.SetArgs(args)

			if err := tCmd.Execute(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for _, file := range tt.expectedRestoredFiles {
				filePath := filepath.Join(orgDir, file)
				if exists := helpers.CheckFileExists(filePath); !exists {
					t.Errorf("expected file %s to exist, but it doesn't", file)
				}
			}
		})
	}
}

func _restoreSetup(t *testing.T, fileNames ...string) (orgDir, backupDir string) {
	orgDir = t.TempDir()
	backupDir = t.TempDir()

	vConfig := vconfig.GetConfig()
	vConfig.SetBackupDir(backupDir)

	files := _createFilesToBackup(orgDir, fileNames...)

	if err := core.BackupFiles(files...); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := os.RemoveAll(orgDir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return orgDir, backupDir
}
