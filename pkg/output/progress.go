package output

import (
	"fmt"
	"os"

	"github.com/schollz/progressbar/v3"
)

type ProgressBar *progressbar.ProgressBar

func NewProgressBar(max int, desc string) ProgressBar {
	return progressbar.NewOptions(max,
		progressbar.OptionUseANSICodes(true),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetWidth(20),
		// progressbar.OptionShowDescriptionAtLineEnd(),
		progressbar.OptionSetDescription(fmt.Sprintf("[yellow]%s[reset]", desc)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[yellow]+[reset]",
			SaucerHead:    "[yellow]+[reset]",
			SaucerPadding: " ",
			BarStart:      "[cyan]|[reset]",
			BarEnd:        "[cyan]|[reset]",
		}),
		progressbar.OptionClearOnFinish())
}
