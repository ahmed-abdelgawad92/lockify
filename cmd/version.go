package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current version of Lockify",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Lockify CLI v%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
