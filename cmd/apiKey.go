package cmd

import (
	"fmt"
	"log"

	"github.com/nexaa-cloud/nexaa-cli/api"
	"github.com/spf13/cobra"
)

var apiKeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "Manage API keys",
}

var listApiKeysCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API Keys",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Errorf("Not implemented yet")
	},
}

var createApiKeyCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new API key",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		authorizations, _ := cmd.Flags().GetString("authorizations")
		//endDate, _ := cmd.Flags().GetString("endDate")

		input := api.ApiKeyInput{
			Name:           &name,
			Authorizations: []string{authorizations},
		}

		client := api.NewClient()
		result, err := client.ApiKeyCreate(input)
		if err != nil {
			log.Println("Make sure you are logged in with your normal credentials for manage API keys.")
			log.Fatalf("Failed to create API key: %v", err)
			return
		}
		log.Println("Created API key: ", result)
	},
}

var revokeApiKeyCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke a API key",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Errorf("Not implemented yet")
	},
}

func init() {
	apiKeyCmd.AddCommand(listApiKeysCmd)

	createApiKeyCmd.Flags().StringP("name", "n", "", "Name of the API key")
	createApiKeyCmd.Flags().StringP("authorizations", "a", "read", "The authorities of the API key (read or read/write). Defaults to read")
	createApiKeyCmd.Flags().String("end-date", "", "The API key is valid until this date. Defaults to 6 months")
	createApiKeyCmd.MarkFlagRequired("name")
	apiKeyCmd.AddCommand(createApiKeyCmd)

	apiKeyCmd.AddCommand(revokeApiKeyCmd)
}
