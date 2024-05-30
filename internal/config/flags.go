package config

import (
	"github.com/spf13/cobra"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&c.Year, "year", 0, "Year to download (default current year)")
	cmd.Flags().IntVar(&c.TileSize, "tile-size", 256, "Tile size")
	cmd.Flags().IntVar(&c.TileMaxX, "tile-max-x", 64, "X tile max")
	cmd.Flags().IntVar(&c.TileMaxY, "tile-max-y", 64, "Y tile max")
	cmd.Flags().IntVar(&c.Zoom, "zoom", 6, "Zoom level")
	cmd.Flags().IntVar(&c.Parallelism, "parallelism", 16, "Number of goroutines to use")
	cmd.Flags().StringVar(&c.Format, "format", "", "Tile format. Try png, png8, png24. (default detected)")
	cmd.Flags().StringVar(&c.URLTemplate, "url-template", "https://tiles.telegeography.com/maps/submarine-cable-map-%d/%d/%d/%d.%s", "URL template. Variables are: year, zoom, x, y, format.")
}
