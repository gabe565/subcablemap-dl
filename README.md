# Submarine Cable Map Downloader
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/gabe565/subcablemap-dl)](https://github.com/gabe565/subcablemap-dl/releases)
[![Build](https://github.com/gabe565/subcablemap-dl/actions/workflows/build.yaml/badge.svg)](https://github.com/gabe565/subcablemap-dl/actions/workflows/build.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gabe565/subcablemap-dl)](https://goreportcard.com/report/github.com/gabe565/subcablemap-dl)

Downloads and combines all tiles for the [Telegeography Submarine Cable Map](https://submarine-cable-map-2024.telegeography.com/) to create a full-resolution image. All maps since 2013 are supported.

![preview](https://github.com/gabe565/subcablemap-dl/assets/7717888/db101cfe-db1a-4c85-a91f-2e2a74d55041)

## Installation

### APT (Ubuntu, Debian)

<details>
  <summary>Click to expand</summary>

1. If you don't have it already, install the `ca-certificates` package
   ```shell
   sudo apt install ca-certificates
   ```

2. Add gabe565 apt repository
   ```
   echo 'deb [trusted=yes] https://apt.gabe565.com /' | sudo tee /etc/apt/sources.list.d/gabe565.list
   ```

3. Update apt repositories
   ```shell
   sudo apt update
   ```

4. Install subcablemap-dl
   ```shell
   sudo apt install subcablemap-dl
   ```
</details>

### RPM (CentOS, RHEL)

<details>
  <summary>Click to expand</summary>

1. If you don't have it already, install the `ca-certificates` package
   ```shell
   sudo dnf install ca-certificates
   ```

2. Add gabe565 rpm repository to `/etc/yum.repos.d/gabe565.repo`
   ```ini
   [gabe565]
   name=gabe565
   baseurl=https://rpm.gabe565.com
   enabled=1
   gpgcheck=0
   ```

3. Install subcablemap-dl
   ```shell
   sudo dnf install subcablemap-dl
   ```
</details>

### AUR (Arch Linux)

<details>
  <summary>Click to expand</summary>

Install [subcablemap-dl-bin](https://aur.archlinux.org/packages/subcablemap-dl-bin) with your [AUR helper](https://wiki.archlinux.org/index.php/AUR_helpers) of choice.
</details>

### Homebrew (macOS, Linux)

<details>
  <summary>Click to expand</summary>

Install subcablemap-dl from [gabe565/homebrew-tap](https://github.com/gabe565/homebrew-tap):
```shell
brew install gabe565/tap/subcablemap-dl
```
</details>

### Docker

<details>
  <summary>Click to expand</summary>

A Docker image is available at [`ghcr.io/gabe565/subcablemap-dl`](https://ghcr.io/gabe565/subcablemap-dl)

#### Usage
```shell
docker run --rm -it -v "$PWD:/data" ghcr.io/gabe565/subcablemap-dl --year=2024
```
</details>

### Manual Installation

<details>
  <summary>Click to expand</summary>

A binary is built for each release. You can either download one of these pre-built release assets, or you can perform a local Go build.

#### Released Binary
1. Download and run the [latest release](https://github.com/gabe565/subcablemap-dl/releases/latest) for your system and architecture.
2. Extract the binary and place it in the desired directory.

#### Local Go Build
```shell
go install github.com/gabe565/subcablemap-dl@latest
```

</details>

## Usage

To download the 2024 map, run
```shell
subcablemap-dl --year=2024
```
When done, the map will be available at `submarine-cable-map-2024.png`.

For full command-line reference, see [docs](./docs/subcablemap-dl.md).
