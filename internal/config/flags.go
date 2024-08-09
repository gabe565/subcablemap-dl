package config

import (
	"strings"

	"github.com/spf13/cobra"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().IntVarP(&c.Year, "year", "y", c.Year, "Year to download (default latest available)")
	cmd.Flags().BoolVarP(&c.NoCrop, "no-crop", "n", c.NoCrop, "Download the entire square map instead of cropping")
	cmd.Flags().IntVar(&c.Tiles.Min.X, "tile-min-x", c.Tiles.Min.X, "X tile min (default determined by year and zoom)")
	cmd.Flags().IntVar(&c.Tiles.Max.X, "tile-max-x", c.Tiles.Max.X, "X tile max (default determined by year and zoom)")
	cmd.Flags().IntVar(&c.Tiles.Min.Y, "tile-min-y", c.Tiles.Min.Y, "Y tile min (default determined by year and zoom)")
	cmd.Flags().IntVar(&c.Tiles.Max.Y, "tile-max-y", c.Tiles.Max.Y, "Y tile max (default determined by year and zoom)")
	cmd.Flags().IntVarP(&c.Zoom, "zoom", "z", c.Zoom, "Zoom level")
	cmd.Flags().IntVarP(&c.Parallelism, "parallelism", "p", c.Parallelism, "Number of goroutines to use")
	cmd.Flags().StringVarP(&c.Format, "format", "f", c.Format, "Tile format. Try png, png8, png24. (default detected)")
	cmd.Flags().VarP(&c.Compression, "compression", "c", "PNG compression level (one of "+strings.Join(CompressionLevelStrings(), ", ")+")")
	cmd.Flags().StringVar(&c.URLTemplate, "url-template", c.URLTemplate, "URL template. Variables are: year, zoom, x, y, format.")
}
