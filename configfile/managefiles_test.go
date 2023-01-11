package configfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckFileExists(t *testing.T) {
	t.Run("returns false for non-existent file", func(t *testing.T) {
		tempDir := t.TempDir()
		file := filepath.Join(tempDir, "testfile")
		if exists := CheckFileExists(file); exists {
			t.Errorf("file exists")
		}
	})

	t.Run("works on directories too", func(t *testing.T) {
		tempDir := t.TempDir()
		if exists := CheckFileExists(tempDir); !exists {
			t.Errorf("directory does not exist")
		}
	})

	t.Run("returns true for existing file", func(t *testing.T) {
		tempDir := t.TempDir()
		file := filepath.Join(tempDir, "testfile")
		f, err := os.Create(file)
		if err != nil {
			t.Errorf("got error: %v", err)
		}
		f.Close()

		if exists := CheckFileExists(file); !exists {
			t.Errorf("file does not exist")
		}
	})
}

func TestEnsureDirExists(t *testing.T) {
	t.Run("does nothing if dir exists", func(t *testing.T) {
		tempDir := t.TempDir()
		if err := EnsureDirExists(tempDir); err != nil {
			t.Errorf("got error: %v", err)
		}

		if exists := CheckFileExists(tempDir); !exists {
			t.Errorf("directory not created")
		}
	})

	t.Run("actually creates dir", func(t *testing.T) {
		tempDir := t.TempDir()
		createdDir := filepath.Join(tempDir, "testdir")
		if err := EnsureDirExists(createdDir); err != nil {
			t.Errorf("got error: %v", err)
		}

		if exists := CheckFileExists(createdDir); !exists {
			t.Errorf("directory not created")
		}
	})
}

func TestCopyFile(t *testing.T) {
	t.Run("returns error if source file does not exist", func(t *testing.T) {
		tempDir := t.TempDir()
		src := filepath.Join(tempDir, "src")
		dest := filepath.Join(tempDir, "dest")
		srcFile, _ := NewConfigFile(src)
		if err := copyFile(dest, srcFile); err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}

func TestRestoreSymLink(t *testing.T) {
	t.Run("creates symlink to the backedup file", func(t *testing.T) {
		backupDir := t.TempDir()
		homedir := t.TempDir()

		expectedSrcFile := filepath.Join(homedir, ".src")

		backupFile, _ := NewConfigFile(expectedSrcFile)

		backupFilePath := filepath.Join(backupDir, backupFile.HashShort())

		if _, err := os.Create(backupFilePath); err != nil {
			t.Errorf("got error: %v", err)
		}

		if err := restoreSymLink(backupDir, backupFile); err != nil {
			t.Errorf("got error: %v", err)
		}

		if exists := CheckFileExists(expectedSrcFile); !exists {
			t.Errorf("file not restored")
		}

		link, err := filepath.EvalSymlinks(expectedSrcFile)
		if err != nil {
			t.Errorf("got error: %v", err)
		}

		if link != backupFilePath {
			t.Errorf("got link: %v, want: %v", link, backupFilePath)
		}
	})
}

func TestCopyAndReplaceFile(t *testing.T) {
	t.Run("fails if the src file doesn't exist", func(t *testing.T) {
		cfgFile, err := NewConfigFile("/tmp/doesnotexist")
		if err != nil {
			t.Errorf("got error: %v", err)
		}

		if err := copyAndReplaceFile(t.TempDir(), cfgFile); err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("copies file, and replaces it with a symlink", func(t *testing.T) {
		backupDir := t.TempDir()
		homedir := t.TempDir()

		srcFile := filepath.Join(homedir, ".src")

		cfgFile, err := NewConfigFile(srcFile)
		if err != nil {
			t.Errorf("got error: %v", err)
		}

		backupFile := filepath.Join(backupDir, cfgFile.HashShort())

		// Create the source file
		f, err := os.Create(srcFile)
		if err != nil {
			t.Errorf("got error: %v", err)
		}
		f.Close()

		if err := copyAndReplaceFile(backupDir, cfgFile); err != nil {
			t.Errorf("got error: %v", err)
		}

		link, err := filepath.EvalSymlinks(srcFile)
		if err != nil {
			t.Errorf("got error: %v", err)
		}

		if link != backupFile {
			t.Errorf("got link: %v, want: %v", link, backupFile)
		}
	})
}
