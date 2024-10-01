package main

import (
	"os"

	"github.com/gabe565/subcablemap-dl/cmd"
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

	root := cmd.New()
	if err := doc.GenMarkdownTree(root, output); err != nil {
		panic(err)
	}
}
