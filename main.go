package main

import (
	"os"

//	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cmd"

	"github.com/SietseVisser/springbeat/beater"
)

var RootCmd = cmd.GenRootCmd("springbeat", "", beater.New)

func main() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
