package main

import (
	"fmt"

	"github.com/osamaadam/cfgrr/cmd"
	"github.com/osamaadam/cfgrr/vconfig"
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

func init() {
	if err := vconfig.GetConfig().Init(); err != nil {
		panic(err)
	}
}
