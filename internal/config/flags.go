package config

import (
	"strings"

	"gabe565.com/utils/cobrax"
	"gabe565.com/utils/must"
	"github.com/spf13/cobra"
)

const (
	FlagBaseURL     = "base-url"
	FlagInsecure    = "insecure"
	FlagYear        = "year"
	FlagNoCrop      = "no-crop"
	FlagTileMinX    = "tile-min-x"
	FlagTileMaxX    = "tile-max-x"
	FlagTileMinY    = "tile-min-y"
	FlagTileMaxY    = "tile-max-y"
	FlagZoom        = "zoom"
	FlagParallelism = "parallelism"
	FlagFormat      = "format"
	FlagCompression = "compression"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	must.Must(cobrax.RegisterCompletionFlag(cmd))
	cmd.Flags().Var(&c.BaseURL, FlagBaseURL, "Base tile download URL")
	cmd.Flags().BoolVarP(&c.Insecure, FlagInsecure, "k", c.Insecure, "Skip HTTPS TLS verification")
	cmd.Flags().IntVarP(&c.Year, FlagYear, "y", c.Year, "Year to download (default latest available)")
	cmd.Flags().BoolVarP(&c.NoCrop, FlagNoCrop, "n", c.NoCrop, "Download the entire square map instead of cropping")
	cmd.Flags().IntVar(&c.Tiles.Min.X, FlagTileMinX, c.Tiles.Min.X, "X tile min (default determined by year and zoom)")
	cmd.Flags().IntVar(&c.Tiles.Max.X, FlagTileMaxX, c.Tiles.Max.X, "X tile max (default determined by year and zoom)")
	cmd.Flags().IntVar(&c.Tiles.Min.Y, FlagTileMinY, c.Tiles.Min.Y, "Y tile min (default determined by year and zoom)")
	cmd.Flags().IntVar(&c.Tiles.Max.Y, FlagTileMaxY, c.Tiles.Max.Y, "Y tile max (default determined by year and zoom)")
	cmd.Flags().IntVarP(&c.Zoom, FlagZoom, "z", c.Zoom, "Zoom level")
	cmd.Flags().IntVarP(&c.Parallelism, FlagParallelism, "p", c.Parallelism, "Number of goroutines to use")
	cmd.Flags().StringVarP(&c.Format, FlagFormat, "f", c.Format, "Tile format. Try png, png8, png24. (default detected)")
	cmd.Flags().VarP(&c.Compression, FlagCompression, "c", "PNG compression level (one of "+strings.Join(CompressionLevelStrings(), ", ")+")")
}
