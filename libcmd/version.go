package libcmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of fuzzagotchi",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Fuzzagotchi version %s\n", rootCmd.Version)
	},
}
