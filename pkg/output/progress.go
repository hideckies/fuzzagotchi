package output

import (
	"fmt"

	"github.com/schollz/progressbar/v3"
)

type ProgressBar *progressbar.ProgressBar

func NewProgressBar(max int, desc string) ProgressBar {
	return progressbar.NewOptions(max,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan]%s[reset]", desc)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[yellow]=[reset]",
			SaucerHead:    "[yellow]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionClearOnFinish())
}
