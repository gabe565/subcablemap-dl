package cmd

import (
	"context"
	"image/png"
	"io"
	"log/slog"
	"os"
	"strconv"

	"gabe565.com/subcablemap-dl/internal/config"
	"gabe565.com/subcablemap-dl/internal/downloader"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func New(options ...Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subcablemap-dl [path]",
		Short: "Download full-resolution versions of Telegeography Submarine Cable Maps",
		Args:  cobra.MaximumNArgs(1),
		RunE:  run,

		ValidArgsFunction: validArgs,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
	}

	conf := config.New()
	conf.RegisterFlags(cmd)
	conf.RegisterCompletions(cmd)
	cmd.SetContext(config.NewContext(context.Background(), conf))

	for _, option := range options {
		option(cmd)
	}

	return cmd
}

func validArgs(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return []string{"png"}, cobra.ShellCompDirectiveFilterFileExt
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	conf, err := config.Load(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	if conf.Completion != "" {
		return completion(cmd)
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
