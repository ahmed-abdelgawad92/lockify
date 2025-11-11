package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export a secret from the vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("exporting a secret from the vault")
		return nil
	},
}

func init() {
	exportCmd.Flags().String("env", defaultEnv, "The environment for which to export the secret")
	exportCmd.Flags().StringP("key", "k", "", "The key to use for exporting the secret")
	exportCmd.Flags().StringP("value", "v", "", "The value to export the secret to")

	rootCmd.AddCommand(exportCmd)
}
