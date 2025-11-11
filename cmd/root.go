package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var defaultEnv string = "local"

var rootCmd = &cobra.Command{
	Use:   "lockify",
	Short: "Lockify securely manages your .env files and secrets",
	Long:  `Lockify is a lightweight CLI tool for securely managing environment variables and .env files locally.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Lockify! Use --help to see available commands.")
	},
}

func Execute() error {
	return rootCmd.Execute()
}
