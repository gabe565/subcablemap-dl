package main

import (
	"log/slog"
	"os"

	"github.com/gabe565/submarine-cable-map-downloader/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
