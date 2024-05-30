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

func (d *Downloader) Do(ctx context.Context) (*image.RGBA, error) {
	img := image.NewRGBA(image.Rect(0, 0, d.config.TileSize*(d.config.TileMaxX+1), d.config.TileSize*(d.config.TileMaxY+1)))
	tileChan := make(chan image.Point)
	group, ctx := errgroup.WithContext(ctx)

	slog.Info("Spawning downloaders", "count", d.config.Parallelism)
	bar := progressbar.Default(int64(d.config.TileMaxX*d.config.TileMaxY), "Creating mosaic")
	for range d.config.Parallelism {
		group.Go(func() error {
			for tile := range tileChan {
				url := fmt.Sprintf(d.config.URLTemplate, d.config.Year, d.config.Zoom, tile.X, tile.Y, d.config.Format)
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
				if err != nil {
					return err
				}

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					return err
				}
				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("%w from %s: %s", ErrUnexpectedResponse, url, resp.Status)
				}

				tileData, err := png.Decode(resp.Body)
				if err != nil {
					return err
				}
				_ = resp.Body.Close()

				pt := image.Point{X: tile.X * d.config.TileSize, Y: tile.Y * d.config.TileSize}
				r := image.Rectangle{Min: pt, Max: pt.Add(image.Point{X: d.config.TileSize, Y: d.config.TileSize})}
				draw.Draw(img, r, tileData, image.Point{}, draw.Src)
				_ = bar.Add(1)
			}

			return nil
		})
	}

	group.Go(func() error {
		defer close(tileChan)
		for x := range d.config.TileMaxX {
			for y := range d.config.TileMaxY {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case tileChan <- image.Point{X: x, Y: y}:
				}
			}
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		_ = bar.Finish()
		return nil, err
	}
	_ = bar.Finish()

	return img, nil
}
