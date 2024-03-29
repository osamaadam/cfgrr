package main

import (
	"fmt"
	"os"

	"github.com/osamaadam/cfgrr/cmd"
)

var (
	version    string
	tagdate    string
	releaseUrl string
)

func main() {
	if err := cmd.Execute(version, tagdate, releaseUrl); err != nil {
		fmt.Printf("Run `%s` to print usage.\n", os.Args[0]+" --help")
	}
}
