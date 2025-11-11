package cmd

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/apixify/lockify/internal/service"
	"github.com/apixify/lockify/internal/vault"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set [env]",
	Short: "add/update a new entry to the vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("seting a new secret to the vault")
		env := args[0]
		var key, value string

		fmt.Println(env, key, value)
		isSecret, _ := cmd.Flags().GetBool("secret")

		prompt := &survey.Input{Message: "Enter key:"}
		survey.AskOne(prompt, &key)

		if isSecret {
			prompt := &survey.Password{Message: "Enter secret:"}
			survey.AskOne(prompt, &value)
		} else {
			prompt = &survey.Input{Message: "Enter value:"}
			survey.AskOne(prompt, &value)
		}

		envPassphraseKey, _ := cmd.Flags().GetString("passphrase-env")
		passphrase := service.NewPassphraseService(env, envPassphraseKey)
		vault, err := vault.Open(env)
		if err != nil {
			return fmt.Errorf("failed to open vault for environment %s: %w", env, err)
		}

		if !vault.VerifyFingerPrint(passphrase.GetPassphrase()) {
			passphrase.ClearPassphrase()
			return errors.New("invalid credentials")
		}

		vault.SetEntry(key, value)
		vault.Save()

		return nil
	},
}

func init() {
	setCmd.Flags().String("passphrase-env", "LOCKIFY_PASSPHRASE", "Name of the environment variable that holds the passphrase")
	setCmd.Flags().BoolP("secret", "s", false, "States that value to set is a secret and should be hidden in the terminal")

	rootCmd.AddCommand(setCmd)
}
