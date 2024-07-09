## submarine-cable-map-downloader



```
submarine-cable-map-downloader [path] [flags]
```

### Options

```
      --format string         Tile format. Try png, png8, png24. (default detected)
  -h, --help                  help for submarine-cable-map-downloader
      --no-crop               Download the entire square map instead of cropping
      --parallelism int       Number of goroutines to use (default 16)
      --tile-max-x int        X tile max (default determined by year) (default 63)
      --tile-max-y int        Y tile max (default determined by year) (default 63)
      --tile-min-x int        X tile min (default determined by year)
      --tile-min-y int        Y tile min (default determined by year)
      --url-template string   URL template. Variables are: year, zoom, x, y, format. (default "https://tiles.telegeography.com/maps/submarine-cable-map-%d/%d/%d/%d.%s")
      --year int              Year to download (default latest available)
      --zoom int              Zoom level (default 6)
```

