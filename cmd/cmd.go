package cmd

import (
	"context"
	"image/png"
	"io"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/gabe565/submarine-cable-map-downloader/internal/config"
	"github.com/gabe565/submarine-cable-map-downloader/internal/downloader"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "submarine-cable-map-downloader",
		RunE: run,

		DisableAutoGenTag: true,
	}

	conf := &config.Config{}
	conf.RegisterFlags(cmd)
	cmd.SetContext(config.NewContext(context.Background(), conf))
	return cmd
}

func run(cmd *cobra.Command, _ []string) error {
	cmd.SilenceUsage = true

	conf, ok := config.FromContext(cmd.Context())
	if !ok {
		panic("Config not added to context")
	}

	if conf.Year == 0 {
		conf.Year = time.Now().Year()
	}

	dl := downloader.New(conf)

	if err := dl.CheckYear(cmd.Context()); err != nil {
		return err
	}

	if err := dl.FindFormat(); err != nil {
		return err
	}

	img, err := dl.Do(cmd.Context())
	if err != nil {
		return err
	}

	path := "submarine-cable-map-" + strconv.Itoa(conf.Year) + ".png"
	slog.Info("Creating file", "path", path)
	out, err := os.Create(path)
	if err != nil {
		return err
	}

	bar := progressbar.DefaultBytes(-1, "Writing to file")
	if err := png.Encode(io.MultiWriter(out, bar), img); err != nil {
		return err
	}
	_ = bar.Finish()
	_, _ = io.WriteString(os.Stderr, "\n")

	slog.Info("Done", "path", path)
	return nil
}
