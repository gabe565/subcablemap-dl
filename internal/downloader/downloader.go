package downloader

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gabe565/submarine-cable-map-downloader/internal/config"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/sync/errgroup"
)

func New(c *config.Config) *Downloader {
	return &Downloader{
		config: c,
	}
}

type Downloader struct {
	config *config.Config
}

var ErrUnexpectedResponse = errors.New("unexpected response")

const URLTemplate = "https://tiles.telegeography.com/maps/submarine-cable-map-%d/%d/%d/%d.%s"

func (d *Downloader) Do(ctx context.Context) (image.Image, error) {
	if err := d.CheckYear(ctx); err != nil {
		return nil, err
	}

	if err := d.FindFormat(ctx); err != nil {
		return nil, err
	}

	slog.Info("Starting download",
		"year", d.config.Year,
		"tiles", d.config.TileCount(),
		"tile_offsets", d.config.Tiles,
		"workers", d.config.Parallelism,
	)

	var img *image.NRGBA
	group, ctx := errgroup.WithContext(ctx)
	group.SetLimit(d.config.Parallelism)
	bar := progressbar.Default(int64(d.config.TileCount()), "Creating mosaic")
	var once sync.Once

	for x := d.config.Tiles.Min.X; x <= d.config.Tiles.Max.X; x++ {
		for y := d.config.Tiles.Min.Y; y <= d.config.Tiles.Max.Y; y++ {
			group.Go(func() error {
				tile := image.Pt(x, y)

				url := fmt.Sprintf(URLTemplate, d.config.Year, d.config.Zoom, tile.X, tile.Y, d.config.Format)
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
				if err != nil {
					return err
				}

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					return err
				}
				defer func() {
					_ = resp.Body.Close()
				}()

				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("%w from %s: %s", ErrUnexpectedResponse, url, resp.Status)
				}

				tileData, err := png.Decode(resp.Body)
				if err != nil {
					return err
				}

				once.Do(func() {
					d.config.TileSize = tileData.Bounds().Max.X
					img = image.NewNRGBA(image.Rect(0, 0, d.config.OutputWidth(), d.config.OutputHeight()))
				})

				draw.Draw(img, d.config.TileRect(tile), tileData, image.Point{}, draw.Src)
				_ = bar.Add(1)
				return nil
			})
		}
	}

	err := group.Wait()
	return img, err
}
