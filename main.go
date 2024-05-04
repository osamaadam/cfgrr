package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/osamaadam/cfgrr/cmd"
)

var (
	version string
	tagdate string
	pkgPath string
)

func main() {
	info, ok := debug.ReadBuildInfo()
	pkgPath = info.Main.Path
	if ok {
		if version == "" {
			version = info.Main.Version
		}
		pkgPath = info.Main.Path
	}
	if err := cmd.Execute(version, tagdate, pkgPath); err != nil {
		fmt.Printf("Run `%s` to print usage.\n", os.Args[0]+" --help")
	}
}
