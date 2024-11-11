package config

import "fmt"

const urlTemplate = "https://tiles.telegeography.com/maps/submarine-cable-map-%d/%d/%d/%d.%s"

func BuildURL(year, zoom, x, y int, format string) string {
	return fmt.Sprintf(urlTemplate, year, zoom, x, y, format)
}
