package cmd

import (
	"context"
	"image"
	"image/png"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gabe565.com/subcablemap-dl/internal/config"
	"gabe565.com/subcablemap-dl/internal/dynamicimage"
	"gabe565.com/utils/bytefmt"
	"gabe565.com/utils/cobrax"
	"github.com/spf13/cobra"
)

func New(options ...cobrax.Option) *cobra.Command {
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

	slog.Info("Submarine Cable Map Downloader",
		"version", cobrax.GetVersion(cmd),
		"commit", cobrax.GetCommit(cmd),
	)

	conf, err := config.Load(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	path := "submarine-cable-map-" + strconv.Itoa(conf.Year) + ".png"
	tmp := "." + path
	if len(args) > 0 {
		path = args[0]
		base, file := filepath.Split(path)
		tmp = filepath.Join(base, "."+file)
	}

	slog.Info("Starting download",
		"year", conf.Year,
		"tiles", conf.TileCount(),
		"bounds", conf.Bounds,
	)

	img, err := dynamicimage.New(cmd.Context(), conf, dynamicimage.WithProgress(!conf.NoProgress))
	if err != nil {
		return err
	}

	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
		_ = os.Remove(tmp)
	}()

	start := time.Now()
	encoder := png.Encoder{CompressionLevel: conf.Compression.ToPNG()}
	if err := encoder.Encode(out, img); err != nil {
		return err
	}
	if img.Error() != nil {
		return img.Error()
	}

	if err := out.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmp, path); err != nil {
		return err
	}

	log := slog.With(
		"path", path,
		"dimensions", image.Pt(conf.Bounds.Dx(), conf.Bounds.Dy()),
	)
	if stat, err := os.Stat(tmp); err == nil {
		log = log.With("size", bytefmt.Encode(stat.Size()))
	}

	log.Info("Done", "took", time.Since(start).Truncate(100*time.Millisecond))
	return nil
}
