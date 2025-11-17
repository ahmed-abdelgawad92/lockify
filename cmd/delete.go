package cmd

import (
	"fmt"

	"github.com/apixify/lockify/internal/di"
	"github.com/spf13/cobra"
)

// lockify del --env [env] --key [key]
var delCmd = &cobra.Command{
	Use:   "del",
	Short: "delete an entry from the vault of an env",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("⏳ removing key...")
		env, err := requireEnvFlag(cmd)
		if err != nil {
			return err
		}

		key, err := cmd.Flags().GetString("key")
		if err != nil {
			return fmt.Errorf("failed to retrieve key flag")
		}
		if key == "" {
			return fmt.Errorf("key is required")
		}

		ctx := getContext()
		useCase := di.BuildDeleteEntry()
		err = useCase.Execute(ctx, env, key)
		if err != nil {
			return err
		}

		di.GetLogger().Success("key %s is removed successfully.\n", key)

		return nil
	},
}

func init() {
	delCmd.Flags().StringP("env", "e", "", "Environment Name")
	delCmd.Flags().StringP("key", "k", "", "key to delete from the vault")
	delCmd.MarkFlagRequired("env")
	delCmd.MarkFlagRequired("key")

	rootCmd.AddCommand(delCmd)
}
