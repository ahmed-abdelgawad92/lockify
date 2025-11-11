package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import a secret from the vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("importing a secret from the vault")
		return nil
	},
}

func init() {
	importCmd.Flags().String("env", defaultEnv, "The environment for which to import the secret")
	importCmd.Flags().StringP("key", "k", "", "The key to use for importing the secret")
	importCmd.Flags().StringP("value", "v", "", "The value to import the secret to")

	rootCmd.AddCommand(importCmd)
}
