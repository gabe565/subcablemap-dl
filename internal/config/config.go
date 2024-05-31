package config

import "image"

type Config struct {
	Year        int
	TileSize    int
	Tiles       image.Rectangle
	Zoom        int
	Parallelism int
	URLTemplate string
	Format      string
}

func (c *Config) OutputWidth() int {
	return (c.Tiles.Max.X - c.Tiles.Min.X + 1) * c.TileSize
}

func (c *Config) OutputHeight() int {
	return (c.Tiles.Max.Y - c.Tiles.Min.Y + 1) * c.TileSize
}

func (c *Config) TileCount() int {
	diff := c.Tiles.Max.Sub(c.Tiles.Min).Add(image.Point{X: 1, Y: 1})
	return diff.X * diff.Y
}
