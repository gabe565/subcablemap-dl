## submarine-cable-map-downloader



```
submarine-cable-map-downloader [flags]
```

### Options

```
      --format string         Tile format. Try png, png8, png24. (default detected)
  -h, --help                  help for submarine-cable-map-downloader
      --parallelism int       Number of goroutines to use (default 16)
      --tile-max-x int        X tile max (default 64)
      --tile-max-y int        Y tile max (default 64)
      --tile-size int         Tile size (default 256)
      --url-template string   URL template. Variables are: year, zoom, x, y, format. (default "https://tiles.telegeography.com/maps/submarine-cable-map-%d/%d/%d/%d.%s")
      --year int              Year to download (default current year)
      --zoom int              Zoom level (default 6)
```

