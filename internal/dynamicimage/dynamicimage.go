package dynamicimage

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"sync"

	"gabe565.com/subcablemap-dl/internal/config"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/sync/errgroup"
)

func New(ctx context.Context, conf *config.Config, opts ...Option) (*DynamicImage, error) {
	d := &DynamicImage{
		ctx:    ctx,
		config: conf,
		tiles:  make([]image.Image, conf.Tiles.Dx()),
	}
	if err := d.downloadRow(0); err != nil {
		d.error = err
		return d, err
	}
	for _, opt := range opts {
		opt(d)
	}
	return d, nil
}

type DynamicImage struct {
	ctx    context.Context
	config *config.Config
	row    int
	tiles  []image.Image
	once   sync.Once
	error  error
	bar    *progressbar.ProgressBar
}

func (d *DynamicImage) ColorModel() color.Model {
	return color.NRGBAModel
}

func (d *DynamicImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, d.config.OutputWidth(), d.config.OutputHeight())
}

func (d *DynamicImage) Opaque() bool {
	return true
}

func (d *DynamicImage) At(x, y int) color.Color {
	if d.error != nil {
		return color.NRGBA{}
	}

	tileRow := y / d.config.TileSize
	if tileRow != d.row {
		if err := d.downloadRow(tileRow); err != nil {
			d.error = err
			return color.NRGBA{}
		}
	}

	if x == 0 && d.bar != nil {
		_ = d.bar.Add64(1)
	}

	img := d.tiles[x/d.config.TileSize]
	if img == nil {
		return color.NRGBA{}
	}

	partX, partY := x%d.config.TileSize, y%d.config.TileSize
	return img.At(partX, partY)
}

func (d *DynamicImage) Error() error {
	return d.error
}

func (d *DynamicImage) DownloadFull() (image.Image, error) {
	img := image.NewNRGBA(d.Bounds())
	draw.Draw(img, d.Bounds(), d, image.Point{}, draw.Src)
	return img, d.error
}

func (d *DynamicImage) downloadRow(y int) error {
	clear(d.tiles)

	group, ctx := errgroup.WithContext(d.ctx)
	group.SetLimit(d.config.Parallelism)
	var mu sync.Mutex

	for x := range d.config.Tiles.Dx() {
		group.Go(func() error {
			tileData, err := DownloadTile(ctx, d.config, image.Pt(x+d.config.Tiles.Min.X, y+d.config.Tiles.Min.Y))
			if err != nil {
				return err
			}

			d.once.Do(func() {
				d.config.TileSize = tileData.Bounds().Max.X
			})

			mu.Lock()
			defer mu.Unlock()
			d.tiles[x] = tileData
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return err
	}
	d.row = y
	return nil
}

var ErrUnexpectedResponse = errors.New("unexpected response")

func DownloadTile(ctx context.Context, conf *config.Config, point image.Point) (image.Image, error) {
	url := conf.BuildURL(conf.Year, conf.Zoom, point.X, point.Y, conf.Format)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w from %s: %s", ErrUnexpectedResponse, url, resp.Status)
	}

	tileData, err := png.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return tileData, nil
}
