package main

import (
	"fmt"

	"github.com/osamaadam/gocfgr/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Print(err)
	}
}
