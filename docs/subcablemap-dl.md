## subcablemap-dl

Download full-resolution versions of Telegeography Submarine Cable Maps

```
subcablemap-dl [path] [flags]
```

### Options

```
      --base-url string      Base tile download URL (default "https://tiles.telegeography.com")
  -c, --compression string   PNG compression level (one of default, none, fast, best) (default "default")
      --crop-bottom int      Adjust the number of pixels to crop on the bottom side (can be positive or negative)
      --crop-left int        Adjust the number of pixels to crop on the left side (can be positive or negative)
      --crop-right int       Adjust the number of pixels to crop on the right side (can be positive or negative)
      --crop-top int         Adjust the number of pixels to crop on the top side (can be positive or negative)
  -f, --format string        Tile format. Try png, png8, png24. (default detected)
      --full-image           Download the entire square map instead of cropping
  -h, --help                 help for subcablemap-dl
  -k, --insecure             Skip HTTPS TLS verification
      --no-progress          Do not show progress bar
  -p, --parallelism int      Number of goroutines to use (default 16)
  -v, --version              version for subcablemap-dl
  -y, --year int             Year to download (default latest available)
  -z, --zoom int             Zoom level (default 6)
```

