package main

import (
	"os"

	"github.com/gabe565/submarine-cable-map-downloader/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		os.Exit(1)
	}
}
