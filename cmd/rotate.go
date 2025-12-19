package cmd

import (
	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/di"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
	"github.com/spf13/cobra"
)

type RotateCommand struct {
	useCase app.RotatePassphraseUc
	prompt  service.PromptService
	logger  domain.Logger
}

func NewRotateCommand(useCase app.RotatePassphraseUc, prompt service.PromptService, logger domain.Logger) *cobra.Command {
	cmd := &RotateCommand{useCase, prompt, logger}

	// lockify rotate-key --env [env]
	cobraCmd := &cobra.Command{
		Use:   "rotate-key",
		Short: "Rotate the passphrase for a vault",
		Long: `Rotate the passphrase for a vault.

This command allows you to change the passphrase for a vault by re-encrypting all entries
with a new passphrase. You will be prompted for the current passphrase and a new passphrase.`,
		Example: `  lockify rotate-key --env prod
  lockify rotate-key --env staging`,
		RunE: cmd.runE,
	}

	cobraCmd.Flags().StringP("env", "e", "", "Environment Name")
	cobraCmd.MarkFlagRequired("env")

	return cobraCmd
}

func (c *RotateCommand) runE(cmd *cobra.Command, args []string) error {
	env, err := requireEnvFlag(cmd)
	if err != nil {
		return err
	}

	passphrase := c.prompt.GetPassphraseInput("Enter current passphrase:")
	newPassphrase := c.prompt.GetPassphraseInput("Enter new passphrase:")

	c.logger.Progress("Rotating passphrase for %s...\n", env)
	ctx := getContext()
	err = c.useCase.Execute(ctx, env, passphrase, newPassphrase)
	if err != nil {
		return err
	}

	clearCacheUseCase := di.BuildClearEnvCachedPassphrase()
	clearCacheUseCase.Execute(ctx, env)

	c.logger.Success("Passphrase rotated successfully")

	return nil
}

func init() {
	rotateCmd := NewRotateCommand(di.BuildRotatePassphrase(), di.BuildPromptService(), di.GetLogger())
	rootCmd.AddCommand(rotateCmd)
}
