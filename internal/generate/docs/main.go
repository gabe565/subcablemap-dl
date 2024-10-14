package main

import (
	"os"

	"gabe565.com/subcablemap-dl/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	output := "./docs"

	if err := os.RemoveAll(output); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		panic(err)
	}

	root := cmd.New(cmd.WithVersion("beta"))
	if err := doc.GenMarkdownTree(root, output); err != nil {
		panic(err)
	}
}
