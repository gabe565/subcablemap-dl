package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func InitLog() {
	slog.SetDefault(slog.New(log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
	})))
}
