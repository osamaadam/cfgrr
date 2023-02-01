package core

import (
	"os"
	"testing"

	cf "github.com/osamaadam/cfgrr/configfile"
	"github.com/osamaadam/cfgrr/helpers"
	"github.com/osamaadam/cfgrr/vconfig"
)

func TestBackupFiles(t *testing.T) {
	tests := []struct {
		name        string
		in          int
		expectedErr bool
	}{
		{"no files", 0, false},
		{"multiple files", 7, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files := _setupBackupEnv(t.TempDir(), t.TempDir(), tt.in)
			if err := BackupFiles(files...); err != nil && !tt.expectedErr {
				t.Errorf("Expected no error, got %s", err)
			}
			for _, f := range files {
				if !helpers.CheckFileExists(f.BackupPath()) {
					t.Errorf("Expected backup file to exist at %s, but it doesn't", f.BackupPath())
				}
			}
		})
	}
}

func TestRestoreFiles(t *testing.T) {
	tests := []struct {
		name        string
		in          int
		expectedErr bool
	}{
		{"no files", 0, false},
		{"multiple files", 7, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files := _setupRestoreEnv(t.TempDir(), t.TempDir(), tt.in)
			if err := RestoreFiles(files...); err != nil && !tt.expectedErr {
				t.Errorf("Expected no error, got %s", err)
			}

			for _, f := range files {
				if !helpers.CheckFileExists(f.PathAbs()) {
					t.Errorf("Expected file to exist at %s, but it doesn't", f.Path)
				}
				if ok, _ := helpers.CheckIfSymlink(f.PathAbs()); !ok {
					t.Errorf("Expected file to be a symlink, but it's not")
				}
			}
		})
	}
}

func _setupBackupEnv(backupDir, dir string, num int) []*cf.ConfigFile {
	c := vconfig.GetConfig()
	c.SetBackupDir(backupDir)
	cfs := make([]*cf.ConfigFile, num)
	for i := 0; i < num; i++ {
		f, _ := os.CreateTemp(dir, "")
		f.Close()
		cfs[i], _ = cf.NewConfigFile(f.Name())
	}

	return cfs
}

func _setupRestoreEnv(backupDir, dir string, num int) []*cf.ConfigFile {
	files := _setupBackupEnv(backupDir, dir, num)
	BackupFiles(files...)
	for _, f := range files {
		os.Remove(f.PathAbs())
	}
	return files
}
