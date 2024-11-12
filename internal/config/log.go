package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func InitLog() {
	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		TimeFormat: time.Stamp,
	})))
}
