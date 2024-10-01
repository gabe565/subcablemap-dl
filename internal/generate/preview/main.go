package main

import (
	"context"
	"image"
	"image/jpeg"
	"log/slog"
	"os"

	"github.com/disintegration/gift"
	"github.com/gabe565/subcablemap-dl/internal/config"
	"github.com/gabe565/subcablemap-dl/internal/downloader"
)

func main() {
	if err := generate(context.Background()); err != nil {
		slog.Error("Exiting due to an error", "error", err.Error())
		os.Exit(1)
	}
}

func generate(ctx context.Context) error {
	years := []int{2013, 2014, 2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024}
	const TileWidth, TileHeight, Cols = 256, 183, 4
	previewImg := image.NewNRGBA(image.Rect(0, 0, TileWidth*Cols, TileHeight*len(years)/Cols))

	g := gift.New(
		gift.ResizeToFill(TileWidth, TileHeight, gift.LanczosResampling, gift.CenterAnchor),
		gift.UnsharpMask(1, 0.5, 0),
	)

	for i, year := range years {
		slog.Info("Fetching tile", "year", year)

		conf := config.New()
		conf.Year = year
		conf.Zoom = 2
		if err := conf.DetermineOffsetsByYear(); err != nil {
			return err
		}

		img, err := downloader.New(conf).Do(ctx)
		if err != nil {
			return err
		}

		pt := image.Point{
			X: (i % Cols) * TileWidth,
			Y: (i / Cols) * TileHeight,
		}
		slog.Info("Writing tile", "year", year, "origin", pt)
		g.DrawAt(previewImg, img, pt, gift.CopyOperator)
	}

	path := "preview.jpg"
	slog.Info("Writing preview image", "path", path)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	if err := jpeg.Encode(f, previewImg, &jpeg.Options{Quality: 100}); err != nil {
		return err
	}

	return f.Close()
}
