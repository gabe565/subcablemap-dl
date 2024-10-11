package config

import "image/png"

//go:generate go run github.com/dmarkham/enumer -type CompressionLevel -trimprefix Compression -transform lower -text -output compression_string.go

type CompressionLevel int

const (
	CompressionDefault CompressionLevel = iota
	CompressionNone
	CompressionFast
	CompressionBest
)

func (c *CompressionLevel) Set(s string) error {
	return c.UnmarshalText([]byte(s))
}

func (c *CompressionLevel) Type() string {
	return "string"
}

func (c *CompressionLevel) ToPNG() png.CompressionLevel {
	switch *c {
	case CompressionDefault:
		return png.DefaultCompression
	case CompressionNone:
		return png.NoCompression
	case CompressionFast:
		return png.BestSpeed
	case CompressionBest:
		return png.BestCompression
	default:
		return png.NoCompression
	}
}
