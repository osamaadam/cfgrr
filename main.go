package main

import (
	"fmt"

	"github.com/osamaadam/cfgrr/cmd"
)

var (
	version string
	tagdate string
)

func main() {
	if err := cmd.Execute(version, tagdate); err != nil {
		fmt.Printf("%+v", err)
	}
}
