package main

import (
	"fmt"
	"os"

	"github.com/kyma-incubator/milv/cli"
	milv "github.com/kyma-incubator/milv/pkg"
)

func main() {
	cliCommands := cli.ParseCommands()
	config, err := milv.NewConfig(cliCommands)
	if err != nil {
		panic(err)
	}

	files, _ := milv.NewFiles(cliCommands.Files, config)
	files.Run(cliCommands.Verbose)

	if files.Summary() {
		os.Exit(1)
	}

	fmt.Println("NO ISSUES :-)")
}
