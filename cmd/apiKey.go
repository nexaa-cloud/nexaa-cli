package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

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
		client := api.NewClient()
		apiKeys, err := client.ApiKeysList()

		if err != nil {
			log.Fatalf("Failed to list API keys: %v", err)
		}

		if len(apiKeys) == 0 {
			fmt.Println("No API keys found.")
			return
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)

		fmt.Fprintln(writer, "IDENTIFIER\t NAME\t AUTHORIZATION\t START DATE\t END DATE\t STATUS\t")

		for _, apiKey := range apiKeys {
			fmt.Fprintf(
				writer, "%s\t %s\t %s\t %s\t %s\t %s\t\n",
				apiKey.Identifier, apiKey.Name, strings.Join(apiKey.Authorizations, ", "), apiKey.StartDate, apiKey.EndDate, apiKey.Status,
			)
		}

		writer.Flush()
	},
}

var createApiKeyCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new API key",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		authorizations, _ := cmd.Flags().GetString("authorizations")
		//endDate, _ := cmd.Flags().GetString("endDate")

		var auths []api.ApiKeyAuthorizationType
		switch authorizations {
		case "read":
			auths = []api.ApiKeyAuthorizationType{api.ApiKeyAuthorizationTypeRead}
		case "read/write":
			auths = []api.ApiKeyAuthorizationType{api.ApiKeyAuthorizationTypeRead, api.ApiKeyAuthorizationTypeWrite}
		default:
			log.Fatalf("invalid authorization %q: must be \"read\" or \"read/write\"", authorizations)
		}

		input := api.ApiKeyInput{
			Name:           &name,
			Authorizations: auths,
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
	createApiKeyCmd.Flags().StringP("authorizations", "a", "read", "The authorizations of the API key (read or read/write). Defaults to read")
	createApiKeyCmd.Flags().String("end-date", "", "The API key is valid until this date. Defaults to 6 months")
	createApiKeyCmd.MarkFlagRequired("name")
	apiKeyCmd.AddCommand(createApiKeyCmd)

	apiKeyCmd.AddCommand(revokeApiKeyCmd)
}
