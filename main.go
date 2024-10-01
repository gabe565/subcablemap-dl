package main

import (
	"log/slog"
	"os"

	"github.com/gabe565/subcablemap-dl/cmd"
)

var version = "beta"

func main() {
	root := cmd.New(cmd.WithVersion(version))
	if err := root.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
