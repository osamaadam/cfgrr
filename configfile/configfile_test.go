package configfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/osamaadam/cfgrr/helpers"
	"github.com/osamaadam/cfgrr/vconfig"
)

func TestNewConfigFile(t *testing.T) {
	t.Setenv("HOME", "/home/user")
	homedir, _ := os.UserHomeDir()

	tests := []struct {
		name string
		in   string
		out  *ConfigFile
		err  bool
	}{
		{"empty path", "", nil, true},
		{"valid path", filepath.Join(homedir, "path/to/file"), &ConfigFile{Path: "path/to/file"}, false},
		{"clean path", filepath.Join(homedir, "path/../path/.config"), &ConfigFile{Path: "path/.config"}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file, err := NewConfigFile(test.in)
			if err != nil && !test.err {
				t.Errorf("got error: %v", err)
			}

			t.Log(file)
			if file != nil && file.Path != test.out.Path {
				t.Errorf("got path: %v, want: %v", file.Path, test.out.Path)
			}
		})
	}
}

func TestConfigFile_Backup(t *testing.T) {
	t.Run("actually backs up", func(t *testing.T) {
		files := _setupBackupEnv(t.TempDir(), t.TempDir(), 1)
		if err := files[0].Backup(); err != nil {
			t.Errorf("expected no error, got %s", err)
		}
		if !helpers.CheckFileExists(files[0].BackupPath()) {
			t.Errorf("expected backup file to exist at %s, but it doesn't", files[0].BackupPath())
		}
	})
}

func TestConfigFile_Restore(t *testing.T) {
	t.Run("actually restores", func(t *testing.T) {
		files := _setupRestoreEnv(t.TempDir(), t.TempDir(), 1)
		file := files[0]
		if err := file.Restore(); err != nil {
			t.Errorf("expected no error, got %s", err)
		}
		if !helpers.CheckFileExists(file.PathAbs()) {
			t.Errorf("expected restored file to exist at %s, but it doesn't", file.PathAbs())
		}
		if ok, _ := helpers.CheckIfSymlink(file.PathAbs()); !ok {
			t.Errorf("expected restored file to be a symlink, but it isn't")
		}
	})
}

func TestConfigFile_DeleteBackup(t *testing.T) {
	tests := []struct {
		name    string
		restore bool
		wantErr bool
	}{
		{"backup, no restore", false, false},
		{"backup, restore", true, false},
	}

	// symlink to file doesn't exist beforehand.
	for _, tt := range tests {
		t.Run("CLEAN: "+tt.name, func(t *testing.T) {
			files := _setupRestoreEnv(t.TempDir(), t.TempDir(), 1)
			file := files[0]
			if err := file.DeleteBackup(tt.restore); (err != nil) != tt.wantErr {
				t.Errorf("ConfigFile.DeleteBackup() error = %v, wantErr %v", err, tt.wantErr)
			}
			if helpers.CheckFileExists(file.BackupPath()) {
				t.Errorf("expected backup file to not exist at %s, but it does", file.BackupPath())
			}
			if tt.restore {
				if !helpers.CheckFileExists(file.PathAbs()) {
					t.Errorf("expected restored file to exist at %s, but it doesn't", file.PathAbs())
				}
				if ok, _ := helpers.CheckIfSymlink(file.PathAbs()); ok {
					t.Errorf("expected restored file to not be a symlink, but it is")
				}
			}
		})
	}

	// Symlink to file exists beforehand.
	for _, tt := range tests {
		t.Run("DIRTY: "+tt.name, func(t *testing.T) {
			files := _setupBackupEnv(t.TempDir(), t.TempDir(), 1)
			file := files[0]
			file.Backup()
			if err := file.DeleteBackup(tt.restore); (err != nil) != tt.wantErr {
				t.Errorf("ConfigFile.DeleteBackup() error = %v, wantErr %v", err, tt.wantErr)
			}
			if helpers.CheckFileExists(file.BackupPath()) {
				t.Errorf("expected backup file to not exist at %s, but it does", file.BackupPath())
			}
			if tt.restore {
				if !helpers.CheckFileExists(file.PathAbs()) {
					t.Errorf("expected restored file to exist at %s, but it doesn't", file.PathAbs())
				}
				if ok, _ := helpers.CheckIfSymlink(file.PathAbs()); ok {
					t.Errorf("expected restored file to not be a symlink, but it is")
				}
			}
		})
	}
}

func _setupBackupEnv(backupDir, dir string, num int) []*ConfigFile {
	c := vconfig.GetConfig()
	c.SetBackupDir(backupDir)
	cfs := make([]*ConfigFile, num)
	for i := 0; i < num; i++ {
		f, _ := os.CreateTemp(dir, "")
		f.Close()
		cfs[i], _ = NewConfigFile(f.Name())
	}

	return cfs
}

func _setupRestoreEnv(backupDir, dir string, num int) []*ConfigFile {
	files := _setupBackupEnv(backupDir, dir, num)
	for _, f := range files {
		f.Backup()
		os.Remove(f.PathAbs())
	}
	return files
}
