package main

import (
	"github.com/jktr/strichliste-cli/cmd"
	"os"
)

func main() {
	err := cmd.NewCLI().RootCommand.Execute()
	if err != nil {
		os.Exit(1)
	}
}
