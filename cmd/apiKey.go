package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var apiKeyCmd = &cobra.Command{
	Use: "apikey",
	Short: "Manage API keys",
}

var listApiKeysCmd = &cobra.Command{
	Use: "list",
	Short: "List all API Keys",
	Run: func(ccmd *cobra.Command, args []string) {
		fmt.Errorf("Not implemented yet")
	},
}

var createApiKeyCmd = &cobra.Command{
	Use: "create",
	Short: "Create a new API key",
		Run: func(ccmd *cobra.Command, args []string) {
		fmt.Errorf("Not implemented yet")
	},
}

var revokeApiKeyCmd = &cobra.Command{
	Use: "revoke",
	Short: "Revoke a API key",
		Run: func(ccmd *cobra.Command, args []string) {
		fmt.Errorf("Not implemented yet")
	},
}

func init() {
	apiKeyCmd.AddCommand(listApiKeysCmd)

	apiKeyCmd.AddCommand(createApiKeyCmd)

	apiKeyCmd.AddCommand(revokeApiKeyCmd)
}