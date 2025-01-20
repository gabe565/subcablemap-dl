package config

import (
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"log/slog"
	"math"
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
	FullImage   bool
	Crop        image.Rectangle
	Bounds      image.Rectangle
	Zoom        int
	Parallelism int
	Format      string
	Compression CompressionLevel
	NoProgress  bool
}

func (c *Config) TilesHorizontal() int {
	return int(math.Ceil(float64(c.Bounds.Dx()) / float64(c.TileSize)))
}

func (c *Config) TilesVertical() int {
	return int(math.Ceil(float64(c.Bounds.Dy()) / float64(c.TileSize)))
}

func (c *Config) TileCount() int {
	return c.TilesHorizontal() * c.TilesVertical()
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
		return image.Pt(Zoom6Max*c.TileSize, Zoom6Max*c.TileSize), nil
	case 5:
		return image.Pt(Zoom5Max*c.TileSize, Zoom5Max*c.TileSize), nil
	case 4:
		return image.Pt(Zoom4Max*c.TileSize, Zoom4Max*c.TileSize), nil
	case 3:
		return image.Pt(Zoom3Max*c.TileSize, Zoom3Max*c.TileSize), nil
	case 2:
		return image.Pt(Zoom2Max*c.TileSize, Zoom2Max*c.TileSize), nil
	}
	return image.Point{}, fmt.Errorf("%w: %d", ErrInvalidZoom, c.Zoom)
}

var (
	ErrBoundsTooSmall = errors.New("bounds too small")
	ErrBoundsTooLarge = errors.New("bounds too large")
)

func (c *Config) UpdateBounds() error {
	var err error
	c.Bounds, err = c.GetYearBounds()
	if err != nil {
		return err
	}

	c.Bounds = image.Rect(
		c.Bounds.Min.X+c.Crop.Min.X,
		c.Bounds.Min.Y+c.Crop.Min.Y,
		c.Bounds.Max.X-c.Crop.Max.X,
		c.Bounds.Max.Y-c.Crop.Max.Y,
	)

	maxPoint, err := c.MaxForZoom()
	if err != nil {
		return err
	}

	if c.Bounds.Min.X < 0 || c.Bounds.Min.Y < 0 {
		return fmt.Errorf("%w: %s must be greater than %s", ErrBoundsTooSmall, c.Bounds.Min, image.Pt(0, 0))
	}
	if c.Bounds.Max.X > maxPoint.X || c.Bounds.Max.Y > maxPoint.Y {
		return fmt.Errorf("%w: %s must be less than %s", ErrBoundsTooLarge, c.Bounds.Max, maxPoint)
	}
	return nil
}

func (c *Config) GetYearBounds() (image.Rectangle, error) {
	bounds, err := c.MaxForZoom()
	if err != nil {
		return image.Rectangle{}, err
	}

	frame := image.Rectangle{Max: bounds}
	if c.FullImage {
		return frame, nil
	}

	zoomOffset := int(math.Pow(2, float64(c.Zoom)-2))
	switch c.Year {
	case 2013:
		frame.Min.Y = 80 * zoomOffset
		frame.Max.Y -= 144 * zoomOffset
	default:
		frame.Min.Y = 144 * zoomOffset
		frame.Max.Y -= 144 * zoomOffset
	}
	return frame, nil
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
		u := c.BuildURL(c.Year, c.Zoom, image.Point{}, v)
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
