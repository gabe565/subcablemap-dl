package config

import (
	"github.com/spf13/cobra"
)

const (
	DefaultFetchMin = 0
	DefaultFetchMax = 63
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&c.Year, "year", 0, "Year to download (default current year)")
	cmd.Flags().IntVar(&c.TileSize, "tile-size", 256, "Tile size")
	cmd.Flags().IntVar(&c.Tiles.Min.X, "tile-min-x", DefaultFetchMin, "X tile min (default determined by year)")
	cmd.Flags().IntVar(&c.Tiles.Max.X, "tile-max-x", DefaultFetchMax, "X tile max (default determined by year)")
	cmd.Flags().IntVar(&c.Tiles.Min.Y, "tile-min-y", DefaultFetchMin, "Y tile min (default determined by year)")
	cmd.Flags().IntVar(&c.Tiles.Max.Y, "tile-max-y", DefaultFetchMax, "Y tile max (default determined by year)")
	cmd.Flags().IntVar(&c.Zoom, "zoom", 6, "Zoom level")
	cmd.Flags().IntVar(&c.Parallelism, "parallelism", 16, "Number of goroutines to use")
	cmd.Flags().StringVar(&c.Format, "format", "", "Tile format. Try png, png8, png24. (default detected)")
	cmd.Flags().StringVar(&c.URLTemplate, "url-template", "https://tiles.telegeography.com/maps/submarine-cable-map-%d/%d/%d/%d.%s", "URL template. Variables are: year, zoom, x, y, format.")
}
