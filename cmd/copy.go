package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "copy a secret from the vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("copying a secret from the vault")
		return nil
	},
}

func init() {
	copyCmd.Flags().String("env", defaultEnv, "The environment for which to copy the secret")
	copyCmd.Flags().StringP("key", "k", "", "The key to use for copying the secret")
	copyCmd.Flags().StringP("value", "v", "", "The value to copy the secret to")

	rootCmd.AddCommand(copyCmd)
}
