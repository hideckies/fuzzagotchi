package output

import (
	"fmt"

	"github.com/fatih/color"
)

const (
	TMPL_BAR_SINGLE_S = "--------------------------------------"
	TMPL_BAR_SINGLE_M = "----------------------------------------------------------------"
	TMPL_BAR_SINGLE_L = "-------------------------------------------------------------------------------------------------"
	TMPL_BAR_DOUBLE_S = "======================================"
	TMPL_BAR_DOUBLE_M = "================================================================"
	TMPL_BAR_DOUBLE_L = "================================================================================================="
)

func Head(title string) {
	color.Yellow(TMPL_BAR_DOUBLE_M)
	fmt.Printf("%s %s\n", color.YellowString("+"), title)
	color.Yellow(TMPL_BAR_DOUBLE_M)
}
