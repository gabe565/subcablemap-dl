package main

import (
	"log/slog"
	"net/http"
	"os"

	"gabe565.com/subcablemap-dl/cmd"
	"gabe565.com/subcablemap-dl/internal/config"
	"gabe565.com/utils/cobrax"
	"gabe565.com/utils/httpx"
)

var version = "beta"

func main() {
	config.InitLog()
	root := cmd.New(cobrax.WithVersion(version))
	http.DefaultTransport = httpx.NewUserAgentTransport(nil, cobrax.BuildUserAgent(root))
	if err := root.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
