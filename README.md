# Submarine Cable Map Downloader

Downloads and combines all tiles for the [Telegeography Submarine Cable Map](https://submarine-cable-map-2024.telegeography.com/) to create a full-resolution image. All maps since 2013 are supported.

![preview](https://github.com/gabe565/subcablemap-dl/assets/7717888/db101cfe-db1a-4c85-a91f-2e2a74d55041)

## Installation

```shell
go install github.com/gabe565/subcablemap-dl@latest
```

## Usage

To download the 2024 map, run
```shell
subcablemap-dl --year=2024
```
When done, the map will be available at `submarine-cable-map-2024.png`.

For full command-line reference, see [docs](./docs/subcablemap-dl.md).
