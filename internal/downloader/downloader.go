package downloader

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/draw"
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

func (d *Downloader) Do(ctx context.Context) (image.Image, error) {
	if err := d.CheckYear(ctx); err != nil {
		return nil, err
	}

	if err := d.FindFormat(ctx); err != nil {
		return nil, err
	}

	var img *image.NRGBA
	tileChan := make(chan image.Point)
	group, ctx := errgroup.WithContext(ctx)

	bar := progressbar.Default(int64(d.config.TileCount()), "Creating mosaic")
	var mu sync.RWMutex
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

				tileData, _, err := image.Decode(resp.Body)
				if err != nil {
					return err
				}
				_ = resp.Body.Close()

				mu.Lock()
				if img == nil {
					d.config.TileSize = tileData.Bounds().Max.X
					img = image.NewNRGBA(image.Rect(0, 0, d.config.OutputWidth(), d.config.OutputHeight()))
				}
				mu.Unlock()

				draw.Draw(img, d.config.TileRect(tile), tileData, image.Point{}, draw.Src)
				_ = bar.Add(1)
			}

			return nil
		})
	}

	group.Go(func() error {
		defer close(tileChan)
		for x := d.config.Tiles.Min.X; x <= d.config.Tiles.Max.X; x++ {
			for y := d.config.Tiles.Min.Y; y <= d.config.Tiles.Max.Y; y++ {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case tileChan <- image.Pt(x, y):
				}
			}
		}
		return nil
	})

	return img, group.Wait()
}
