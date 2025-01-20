package dynamicimage

import (
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

type Option func(d *DynamicImage)

func WithProgress(enabled bool) Option {
	return func(d *DynamicImage) {
		if !enabled {
			return
		}
		d.bar = progressbar.NewOptions64(
			int64(d.config.Bounds.Dy()),
			progressbar.OptionSetDescription("Downloading"),
			progressbar.OptionSetWriter(os.Stderr),
			progressbar.OptionThrottle(65*time.Millisecond),
			progressbar.OptionOnCompletion(func() {
				_, _ = fmt.Fprintln(os.Stderr)
			}),
			progressbar.OptionFullWidth(),
			progressbar.OptionShowElapsedTimeOnFinish(),
		)
	}
}
