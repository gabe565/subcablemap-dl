package config

import (
	"errors"
	"fmt"
	"image"
)

func New() *Config {
	return &Config{
		TileSize:    256,
		Zoom:        6,
		Parallelism: 16,
		URLTemplate: "https://tiles.telegeography.com/maps/submarine-cable-map-%d/%d/%d/%d.%s",
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

var ErrInvalidZoom = errors.New("invalid zoom")

const (
	Zoom6Max = 63
	Zoom5Max = 31
	Zoom4Max = 15
	Zoom3Max = 7
	Zoom2Max = 3
)

func (c *Config) MaxForZoom() (image.Point, error) {
	switch c.Zoom {
	case 6:
		return image.Pt(Zoom6Max, Zoom6Max), nil
	case 5:
		return image.Pt(Zoom5Max, Zoom5Max), nil
	case 4:
		return image.Pt(Zoom4Max, Zoom4Max), nil
	case 3:
		return image.Pt(Zoom3Max, Zoom3Max), nil
	case 2:
		return image.Pt(Zoom2Max, Zoom2Max), nil
	}
	return image.Point{}, fmt.Errorf("%w: %d", ErrInvalidZoom, c.Zoom)
}

func (c *Config) TileRect(tile image.Point) image.Rectangle {
	x := (tile.X - c.Tiles.Min.X) * c.TileSize
	y := (tile.Y - c.Tiles.Min.Y) * c.TileSize
	return image.Rect(x, y, x+c.TileSize, y+c.TileSize)
}

var (
	ErrMaxXTooLarge = errors.New("tile max x exceeds zoom level")
	ErrMaxYTooLarge = errors.New("tile max y exceeds zoom level")
)

func (c *Config) DetermineOffsetsByYear() error {
	maxPoint, err := c.MaxForZoom()
	if err != nil {
		return err
	}
	newTiles := image.Rectangle{Min: c.Tiles.Min, Max: maxPoint}

	if !c.NoCrop {
		switch c.Zoom {
		case 6:
			switch c.Year {
			case 2013:
				newTiles.Min.Y = 5
				newTiles.Max.Y = 55
			case 2020:
				newTiles.Min.Y = 7
				newTiles.Max.Y = 54
			default:
				newTiles.Min.Y = 8
				newTiles.Max.Y = 55
			}
		case 5:
			newTiles.Max.Y = 27
			switch c.Year {
			case 2013:
				newTiles.Min.Y = 2
			case 2020:
				newTiles.Min.Y = 3
			default:
				newTiles.Min.Y = 4
			}
		case 4:
			newTiles.Max.Y = 13
			switch c.Year {
			case 2013:
				newTiles.Min.Y = 1
			default:
				newTiles.Min.Y = 2
			}
		case 3:
			newTiles.Max.Y = 6
			switch c.Year {
			case 2013:
			default:
				newTiles.Min.Y = 1
			}
		}
	}

	if c.Tiles.Min.X == 0 {
		c.Tiles.Min.X = newTiles.Min.X
	}
	if c.Tiles.Min.Y == 0 {
		c.Tiles.Min.Y = newTiles.Min.Y
	}
	if c.Tiles.Max.X == 0 {
		c.Tiles.Max.X = newTiles.Max.X
	}
	if c.Tiles.Max.Y == 0 {
		c.Tiles.Max.Y = newTiles.Max.Y
	}

	if c.Tiles.Max.X > maxPoint.X {
		return ErrMaxXTooLarge
	}
	if c.Tiles.Max.Y > maxPoint.Y {
		return ErrMaxYTooLarge
	}
	return nil
}
