package config

import (
	"image"
	"net/url"
	"path"
	"strconv"
)

func (c *Config) BuildURL(year, zoom int, pt image.Point, format string) *url.URL {
	u := *c.BaseURL.URL
	u.Path = path.Join(u.Path, "maps", "submarine-cable-map-"+strconv.Itoa(year))
	if c.Year == 2012 {
		u.Path = path.Join(u.Path, "1.0.0", "submarine-cable-map-"+strconv.Itoa(year))
	}
	u.Path = path.Join(u.Path, strconv.Itoa(zoom), strconv.Itoa(pt.X), strconv.Itoa(pt.Y))
	u.Path += "." + format
	return &u
}
