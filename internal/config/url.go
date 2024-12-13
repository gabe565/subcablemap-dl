package config

import (
	"net/url"
	"path"
	"strconv"
)

func (c *Config) BuildURL(year, zoom, x, y int, format string) *url.URL {
	u := *c.BaseURL.URL
	u.Path = path.Join(
		u.Path,
		"maps",
		"submarine-cable-map-"+strconv.Itoa(year),
		strconv.Itoa(zoom),
		strconv.Itoa(x),
		strconv.Itoa(y)+"."+format,
	)
	return &u
}
