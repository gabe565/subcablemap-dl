package config

type Config struct {
	Year        int
	TileSize    int
	TileMaxX    int
	TileMaxY    int
	Zoom        int
	Parallelism int
	URLTemplate string
	Format      string
}
