package main

import (
	"os"

	"github.com/dpattee/speedbeat/cmd"

	_ "github.com/dpattee/speedbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
