package config

import (
	"strings"

	"gabe565.com/utils/cobrax"
	"gabe565.com/utils/must"
	"github.com/spf13/cobra"
)

const (
	FlagBaseURL     = "base-url"
	FlagInsecure    = "insecure"
	FlagYear        = "year"
	FlagFullImage   = "full-image"
	FlagCropLeft    = "crop-left"
	FlagCropRight   = "crop-right"
	FlagCropTop     = "crop-top"
	FlagCropBottom  = "crop-bottom"
	FlagZoom        = "zoom"
	FlagParallelism = "parallelism"
	FlagFormat      = "format"
	FlagCompression = "compression"
	FlagNoProgress  = "no-progress"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	must.Must(cobrax.RegisterCompletionFlag(cmd))
	fs := cmd.Flags()
	fs.Var(&c.BaseURL, FlagBaseURL, "Base tile download URL")
	fs.BoolVarP(&c.Insecure, FlagInsecure, "k", c.Insecure, "Skip HTTPS TLS verification")
	fs.IntVarP(&c.Year, FlagYear, "y", c.Year, "Year to download (default latest available)")
	fs.BoolVar(&c.FullImage, FlagFullImage, c.FullImage, "Download the entire square map instead of cropping")
	fs.IntVar(&c.Crop.Min.X, FlagCropLeft, c.Crop.Min.X, "Adjust the number of pixels to crop on the left side (can be positive or negative)")
	fs.IntVar(&c.Crop.Max.X, FlagCropRight, c.Crop.Max.X, "Adjust the number of pixels to crop on the right side (can be positive or negative)")
	fs.IntVar(&c.Crop.Min.Y, FlagCropTop, c.Crop.Min.Y, "Adjust the number of pixels to crop on the top side (can be positive or negative)")
	fs.IntVar(&c.Crop.Max.Y, FlagCropBottom, c.Crop.Max.Y, "Adjust the number of pixels to crop on the bottom side (can be positive or negative)")
	fs.IntVarP(&c.Zoom, FlagZoom, "z", c.Zoom, "Zoom level")
	fs.IntVarP(&c.Parallelism, FlagParallelism, "p", c.Parallelism, "Number of goroutines to use")
	fs.StringVarP(&c.Format, FlagFormat, "f", c.Format, "Tile format. Try png, png8, png24. (default detected)")
	fs.VarP(&c.Compression, FlagCompression, "c", "PNG compression level (one of "+strings.Join(CompressionLevelStrings(), ", ")+")")
	fs.BoolVar(&c.NoProgress, FlagNoProgress, c.NoProgress, "Do not show progress bar")
}
