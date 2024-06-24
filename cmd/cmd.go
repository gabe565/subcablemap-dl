package cmd

import (
	"context"
	"image/png"
	"io"
	"log/slog"
	"os"
	"strconv"

	"github.com/gabe565/submarine-cable-map-downloader/internal/config"
	"github.com/gabe565/submarine-cable-map-downloader/internal/downloader"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "submarine-cable-map-downloader [path]",
		RunE: run,

		DisableAutoGenTag: true,
	}

	conf := config.New()
	conf.RegisterFlags(cmd)
	cmd.SetContext(config.NewContext(context.Background(), conf))
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	conf, ok := config.FromContext(cmd.Context())
	if !ok {
		panic("Config not added to context")
	}

	dl := downloader.New(conf)

	if err := dl.CheckYear(cmd.Context()); err != nil {
		return err
	}
	conf.DetermineOffsetsByYear()

	if err := dl.FindFormat(cmd.Context()); err != nil {
		return err
	}

	slog.Info("Starting download",
		"year", conf.Year,
		"tiles", conf.TileCount(),
		"tile_offsets", conf.Tiles,
		"workers", conf.Parallelism,
	)

	img, err := dl.Do(cmd.Context())
	if err != nil {
		return err
	}

	path := "submarine-cable-map-" + strconv.Itoa(conf.Year) + ".png"
	if len(args) > 0 {
		path = args[0]
	}

	slog.Info("Creating file", "path", path)
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	bar := progressbar.DefaultBytes(-1, "Writing to file")
	if err := png.Encode(io.MultiWriter(out, bar), img); err != nil {
		return err
	}
	_ = bar.Exit()

	slog.Info("Done", "path", path)
	return out.Close()
}
