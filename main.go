package main

import (
	"os"

	"github.com/SietseVisser/springbeat/cmd"

	_ "github.com/SietseVisser/springbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
