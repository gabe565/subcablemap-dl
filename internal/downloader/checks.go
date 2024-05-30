package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

var ErrMissingYear = errors.New("could not find year")

func (d *Downloader) CheckYear(ctx context.Context) error {
	url := "https://submarine-cable-map-" + strconv.Itoa(d.config.Year) + ".telegeography.com"

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %d", ErrMissingYear, d.config.Year)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	return nil
}

var ErrNoFormat = errors.New("could not discover file format")

func (d *Downloader) FindFormat() error {
	if d.config.Format == "" {
		for _, v := range []string{"png", "png8", "png24"} {
			//nolint:noctx
			resp, err := http.Head(fmt.Sprintf(d.config.URLTemplate, d.config.Year, d.config.Zoom, 0, 0, v))
			if err != nil || resp.StatusCode != http.StatusOK {
				continue
			}
			_ = resp.Body.Close()

			d.config.Format = v
			slog.Info("Discovered format", "format", d.config.Format)
		}
		if d.config.Format == "" {
			return ErrNoFormat
		}
	}
	return nil
}
