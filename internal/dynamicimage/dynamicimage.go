package dynamicimage

import (
	"context"
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
		tiles:  make([]image.Image, conf.TilesHorizontal()),
		row:    -1,
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
	error  error
	bar    *progressbar.ProgressBar
}

func (d *DynamicImage) ColorModel() color.Model {
	return color.NRGBAModel
}

func (d *DynamicImage) Bounds() image.Rectangle {
	return image.Rectangle{Max: image.Pt(d.config.Bounds.Dx(), d.config.Bounds.Dy())}
}

func (d *DynamicImage) Opaque() bool {
	return true
}

func (d *DynamicImage) At(x, y int) color.Color {
	if d.error != nil {
		return color.NRGBA{}
	}

	x += d.config.Bounds.Min.X
	y += d.config.Bounds.Min.Y

	tileRow := y / d.config.TileSize
	if tileRow != d.row {
		if err := d.downloadRow(tileRow); err != nil {
			d.error = err
			return color.NRGBA{}
		}
	}

	if x == d.config.Bounds.Min.X && d.bar != nil {
		_ = d.bar.Add64(1)
	}

	img := d.tiles[x/d.config.TileSize]
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

	for x := range d.config.TilesHorizontal() {
		group.Go(func() error {
			tileData, err := DownloadTile(ctx, d.config, image.Pt(x, y))
			if err != nil {
				return err
			}

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

func DownloadTile(ctx context.Context, conf *config.Config, pt image.Point) (image.Image, error) {
	url := conf.BuildURL(conf.Year, conf.Zoom, pt, conf.Format)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := conf.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w from %q: %s", config.ErrUnexpectedResponse, url, resp.Status)
	}

	tileData, err := png.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return tileData, nil
}
