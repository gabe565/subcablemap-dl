package config

import (
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"gabe565.com/utils/pflagx"
)

func New() *Config {
	return &Config{
		BaseURL:     pflagx.URL{URL: &url.URL{Scheme: "https", Host: "tiles.telegeography.com"}},
		Client:      &http.Client{},
		TileSize:    256,
		Zoom:        6,
		Parallelism: 16,
	}
}

type Config struct {
	Completion  string
	BaseURL     pflagx.URL
	Insecure    bool
	Client      *http.Client
	Year        int
	TileSize    int
	NoCrop      bool
	Tiles       image.Rectangle
	Zoom        int
	Parallelism int
	Format      string
	Compression CompressionLevel
}

func (c *Config) OutputWidth() int {
	return (c.Tiles.Max.X - c.Tiles.Min.X) * c.TileSize
}

func (c *Config) OutputHeight() int {
	return (c.Tiles.Max.Y - c.Tiles.Min.Y) * c.TileSize
}

func (c *Config) TileCount() int {
	return c.Tiles.Dx() * c.Tiles.Dy()
}

var ErrInvalidZoom = errors.New("invalid zoom")

const (
	Zoom6Max = 64
	Zoom5Max = 32
	Zoom4Max = 16
	Zoom3Max = 8
	Zoom2Max = 4
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
				newTiles.Max.Y = 56
			case 2020:
				newTiles.Min.Y = 7
				newTiles.Max.Y = 55
			default:
				newTiles.Min.Y = 8
				newTiles.Max.Y = 56
			}
		case 5:
			newTiles.Max.Y = 28
			switch c.Year {
			case 2013:
				newTiles.Min.Y = 2
			case 2020:
				newTiles.Min.Y = 3
			default:
				newTiles.Min.Y = 4
			}
		case 4:
			newTiles.Max.Y = 14
			switch c.Year {
			case 2013:
				newTiles.Min.Y = 1
			default:
				newTiles.Min.Y = 2
			}
		case 3:
			newTiles.Max.Y = 7
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

var ErrMissingYear = errors.New("map is not available for year")

func (c *Config) CheckYear(ctx context.Context) error {
	var latest bool
	if c.Year == 0 {
		latest = true
		c.Year = time.Now().Year()
	}

	u := *c.BaseURL.URL
	u.Path = path.Join(u.Path, "maps/submarine-cable-map-"+strconv.Itoa(c.Year))

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(req)
	if resp != nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}
	if err != nil || resp.StatusCode == http.StatusNotFound {
		if latest {
			slog.Warn("Map for " + strconv.Itoa(c.Year) + " is not yet available, downloading " + strconv.Itoa(c.Year-1) + ".")
			c.Year--
			return c.CheckYear(ctx)
		}
		return fmt.Errorf("%w: %d", ErrMissingYear, c.Year)
	}

	return nil
}

var (
	ErrNoFormat           = errors.New("could not discover file format")
	ErrUnexpectedResponse = errors.New("unexpected response")
)

func (c *Config) FindFormat(ctx context.Context) error {
	if c.Format != "" {
		return nil
	}

	var errs []error
	for _, v := range []string{"png", "png8", "png24"} {
		u := c.BuildURL(c.Year, c.Zoom, 0, 0, v)
		req, err := http.NewRequestWithContext(ctx, http.MethodHead, u.String(), nil)
		if err != nil {
			return err
		}

		resp, err := c.Client.Do(req)
		if resp != nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}
		switch {
		case err != nil:
			errs = append(errs, err)
		case resp.StatusCode != http.StatusOK:
			errs = append(errs, fmt.Errorf("%w from %q: %s", ErrUnexpectedResponse, u.String(), resp.Status))
		default:
			c.Format = v
			return nil
		}
	}

	errs = append([]error{ErrNoFormat}, errs...)
	return errors.Join(errs...)
}
