package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const Version = "0.1.3"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of Fuzzagotchi",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Fuzzagotchi v%s\n", Version)
		os.Exit(0)
	},
}
