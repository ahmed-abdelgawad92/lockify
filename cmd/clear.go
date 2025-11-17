package cmd

import (
	"github.com/apixify/lockify/internal/di"
	"github.com/spf13/cobra"
)

// lockify cache clear
var clearCmd = &cobra.Command{
	Use:   "cache clear",
	Short: "Clear cached passphrase.",
	RunE: func(cmd *cobra.Command, args []string) error {
		di.GetLogger().Progress("clearing cached passphrases")
		useCase := di.BuildClearCachedPassphrase()

		ctx := getContext()
		err := useCase.Execute(ctx)
		if err != nil {
			di.GetLogger().Error("failed to cleare cached passphrases")
			return err
		}

		di.GetLogger().Success("cleared cached passphrases")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
