package config

import "image"

func New() *Config {
	return &Config{
		Zoom: 6,
	}
}

type Config struct {
	Year        int
	TileSize    int
	NoCrop      bool
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

func (c *Config) DetermineOffsetsByYear() {
	if c.NoCrop {
		return
	}

	newTiles := image.Rectangle{Min: c.Tiles.Min, Max: image.Point{X: DefaultFetchMax, Y: DefaultFetchMax}}
	switch c.Year {
	case 2013:
		newTiles.Min.Y = 5
		newTiles.Max.Y = 55
	case 2020:
		newTiles.Min.Y = 7
		newTiles.Max.Y = 54
	case 2014, 2015, 2016, 2017, 2018, 2019, 2021, 2022, 2023, 2024:
		newTiles.Min.Y = 8
		newTiles.Max.Y = 55
	default:
		return
	}

	if c.Tiles.Min.X == DefaultFetchMin {
		c.Tiles.Min.X = newTiles.Min.X
	}
	if c.Tiles.Min.Y == DefaultFetchMin {
		c.Tiles.Min.Y = newTiles.Min.Y
	}
	if c.Tiles.Max.X == DefaultFetchMax {
		c.Tiles.Max.X = newTiles.Max.X
	}
	if c.Tiles.Max.Y == DefaultFetchMax {
		c.Tiles.Max.Y = newTiles.Max.Y
	}
}
