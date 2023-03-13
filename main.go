package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/osamaadam/cfgrr/cmd"
	"github.com/osamaadam/cfgrr/helpers"
)

var (
	version     string
	tagdate     string
	ErrorLogger *log.Logger
)

func main() {
	if err := cmd.Execute(version, tagdate); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		fmt.Printf("Run `%s` to print usage.\n", os.Args[0]+" --help")
		ErrorLogger.Fatalf("\nargs: %s\n%+v", strings.Join(os.Args, " "), err)
	}
}

func init() {
	ErrorLogger = log.New(os.Stderr, "ERROR at ", log.Ldate|log.Ltime|log.Lshortfile)
	homedir, _ := os.UserHomeDir()
	shareDir := filepath.Join(homedir, ".local", "share", "cfgrr")
	if err := helpers.EnsureDirExists(shareDir); err != nil {
		return
	}
	logfilePath := filepath.Join(shareDir, "cfgrr.log")
	file, err := os.OpenFile(logfilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	ErrorLogger.SetOutput(file)
}
