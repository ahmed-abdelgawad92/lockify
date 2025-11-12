package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/apixify/lockify/internal/service"
	"github.com/apixify/lockify/internal/vault"
	"github.com/spf13/cobra"
)

// lockify rotate-key --env [env]
var rotateCmd = &cobra.Command{
	Use:   "rotate-key",
	Short: "Decrypt all variables and export them in a specific format.",
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := cmd.Flags().GetString("env")
		if err != nil {
			return fmt.Errorf("failed to retrieve env flag")
		}
		if env == "" {
			return fmt.Errorf("env is required")
		}
		fmt.Fprintf(os.Stderr, "⏳ Rotating passphrase for %s...\n", env)

		passphraseService := service.NewPassphraseService(env)
		passphraseService.ClearPassphrase()

		vault, err := vault.Open(env)
		if err != nil {
			return fmt.Errorf("failed to open vault for environment %s: %w", env, err)
		}

		passphrase := passphraseService.GetPassphrase()
		passphraseService.ClearPassphrase()
		if !vault.VerifyFingerPrint(passphrase) {
			return fmt.Errorf("invalid credentials")
		}

		oldCryptoService, err := service.NewCryptoService(vault.Meta.Salt, passphrase)
		if err != nil {
			return fmt.Errorf("failed to initialize crypto service: %w", err)
		}
		vault.Meta.Salt, err = service.GenerateSalt(16)
		if err != nil {
			return fmt.Errorf("failed to generate salt")
		}

		var newPassphrase string
		prompt := &survey.Password{Message: "Enter new passphrase:"}
		survey.AskOne(prompt, &newPassphrase)
		newCryptoService, err := service.NewCryptoService(vault.Meta.Salt, newPassphrase)
		if err != nil {
			return fmt.Errorf("failed to initialize crypto service: %w", err)
		}

		vault.Meta.FingerPrint, err = vault.GenerateFingerprint(newPassphrase)
		if err != nil {
			return fmt.Errorf("failed to generate fingerprint")
		}

		for key := range vault.Entries {
			entry := vault.Entries[key]
			entry.Value, _ = oldCryptoService.DecryptValue(entry.Value)
			entry.Value, _ = newCryptoService.EncryptValue([]byte(entry.Value))

			vault.Entries[key] = entry
		}

		err = vault.Save()
		if err != nil {
			return fmt.Errorf("failed to save vault")
		}

		return nil
	},
}

func init() {
	rotateCmd.Flags().StringP("env", "e", "", "Environment Name")

	rootCmd.AddCommand(rotateCmd)
}
