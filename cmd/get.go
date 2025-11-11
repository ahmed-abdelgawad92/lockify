package cmd

import (
	"errors"
	"fmt"

	"github.com/apixify/lockify/internal/service"
	"github.com/apixify/lockify/internal/vault"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [env]",
	Short: "get a secret from the vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("getting a secret from the vault")
		env := args[0]
		key, _ := cmd.Flags().GetString("key")
		passphrase := service.NewPassphraseService(env)
		vault, err := vault.Open(env)
		if err != nil {
			return fmt.Errorf("failed to open vault for environment %s: %w", env, err)
		}
		if !vault.VerifyFingerPrint(passphrase.GetPassphrase()) {
			passphrase.ClearPassphrase()
			return errors.New("invalid credentials")
		}

		entry, err := vault.GetEntry(key)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		fmt.Println(entry.Value)

		return nil
	},
}

func init() {
	getCmd.Flags().StringP("key", "k", "", "The key to use for getting the secret")

	rootCmd.AddCommand(getCmd)
}
