package cmd

import (
	"context"
	"image/png"
	"io"
	"log/slog"
	"os"
	"strconv"

	"github.com/gabe565/subcablemap-dl/internal/config"
	"github.com/gabe565/subcablemap-dl/internal/downloader"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "subcablemap-dl [path]",
		Short:   "Download full-resolution versions of Telegeography Submarine Cable Maps",
		RunE:    run,
		Version: buildVersion(),

		SilenceErrors:     true,
		DisableAutoGenTag: true,
	}
	cmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "commit %s" .Version}}
`)

	conf := config.New()
	conf.RegisterFlags(cmd)
	cmd.SetContext(config.NewContext(context.Background(), conf))
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	config.InitLog()

	conf, ok := config.FromContext(cmd.Context())
	if !ok {
		panic("Config not added to context")
	}

	if conf.Completion != "" {
		return completion(cmd)
	}

	if err := conf.DetermineOffsetsByYear(); err != nil {
		return err
	}

	img, err := downloader.New(conf).Do(cmd.Context())
	if err != nil {
		return err
	}

	path := "submarine-cable-map-" + strconv.Itoa(conf.Year) + ".png"
	if len(args) > 0 {
		path = args[0]
	}

	slog.Info("Creating file", "path", path, "dimensions", img.Bounds().Max)
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	bar := progressbar.DefaultBytes(-1, "Writing to file")
	encoder := png.Encoder{CompressionLevel: conf.Compression.ToPNG()}
	if err := encoder.Encode(io.MultiWriter(out, bar), img); err != nil {
		return err
	}
	_ = bar.Exit()

	if err := out.Close(); err != nil {
		return err
	}

	slog.Info("Done", "path", path)
	return nil
}
