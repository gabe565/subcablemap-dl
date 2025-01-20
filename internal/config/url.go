package config

import (
	"image"
	"net/url"
	"path"
	"strconv"
)

func (c *Config) BuildURL(year, zoom int, pt image.Point, format string) *url.URL {
	u := *c.BaseURL.URL
	u.Path = path.Join(
		u.Path,
		"maps",
		"submarine-cable-map-"+strconv.Itoa(year),
		strconv.Itoa(zoom),
		strconv.Itoa(pt.X),
		strconv.Itoa(pt.Y)+"."+format,
	)
	return &u
}
