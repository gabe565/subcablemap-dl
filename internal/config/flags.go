package config

import (
	"github.com/spf13/cobra"
)

const (
	DefaultFetchMin = 0
	DefaultFetchMax = 63
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&c.Year, "year", c.Year, "Year to download (default latest available)")
	cmd.Flags().IntVar(&c.TileSize, "tile-size", c.TileSize, "Tile size")
	cmd.Flags().BoolVar(&c.NoCrop, "no-crop", c.NoCrop, "Download the entire square map instead of cropping")
	cmd.Flags().IntVar(&c.Tiles.Min.X, "tile-min-x", c.Tiles.Min.X, "X tile min (default determined by year)")
	cmd.Flags().IntVar(&c.Tiles.Max.X, "tile-max-x", c.Tiles.Max.X, "X tile max (default determined by year)")
	cmd.Flags().IntVar(&c.Tiles.Min.Y, "tile-min-y", c.Tiles.Min.Y, "Y tile min (default determined by year)")
	cmd.Flags().IntVar(&c.Tiles.Max.Y, "tile-max-y", c.Tiles.Max.Y, "Y tile max (default determined by year)")
	cmd.Flags().IntVar(&c.Parallelism, "parallelism", c.Parallelism, "Number of goroutines to use")
	cmd.Flags().StringVar(&c.Format, "format", c.Format, "Tile format. Try png, png8, png24. (default detected)")
	cmd.Flags().StringVar(&c.URLTemplate, "url-template", c.URLTemplate, "URL template. Variables are: year, zoom, x, y, format.")
}
