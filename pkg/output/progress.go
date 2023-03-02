package output

import (
	"time"

	"github.com/schollz/progressbar/v3"
)

type ProgressBar *progressbar.ProgressBar

func NewProgressBar(max int, desc string) ProgressBar {
	return progressbar.NewOptions(max,
		progressbar.OptionUseANSICodes(true),
		// progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionThrottle(1*time.Millisecond),
		progressbar.OptionSetWidth(10),
		progressbar.OptionShowDescriptionAtLineEnd(),
		progressbar.OptionSetDescription("\r"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[yellow]+[reset]",
			SaucerHead:    "[yellow]+[reset]",
			SaucerPadding: " ",
			BarStart:      "[cyan]|[reset]",
			BarEnd:        "[cyan]|[reset]",
		}),
		progressbar.OptionClearOnFinish())
}
