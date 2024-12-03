package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"gitlab.com/Tilaa/tilaa-cli/api"
)

var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Manage private registries",
}

var listRegistriesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all private registries",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")

		registries, err := api.ListRegistries(namespace)

		if err != nil {
			log.Fatalf("Failed to list registries: %v", err)
		}

		if len(registries) == 0 {
			fmt.Println("No registries found.")
			return
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)

		fmt.Fprintln(writer, "NAME\t STATE\t IMAGE\t")

		for _, registry := range registries {
			fmt.Fprintf(writer, "%s\t %s\t\n", registry.Id, registry.Name)
		}

		writer.Flush()
	},
}

var createRegistryCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new private registry",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetInt("namespace")
		name, _ := cmd.Flags().GetString("name")
		source, _ := cmd.Flags().GetString("source")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		input := api.RegistryInput{
			Namespace: namespace,
			Name:      name,
			Source:    source,
			Username:  username,
			Password:  password,
			Verify:    true,
		}

		_, err := api.CreateRegistry(input)

		if err != nil {
			log.Fatalf("Failed to create registry: %v", err)
		}
	},
}

func init() {
	listRegistriesCmd.Flags().StringP("namespace", "n", "", "Namespace")
	listRegistriesCmd.MarkFlagRequired("namespace")
	registryCmd.AddCommand(listRegistriesCmd)

	createRegistryCmd.Flags().IntP("namespace", "n", 0, "Namespace")
	createRegistryCmd.Flags().String("name", "", "A unique name for the private registry")
	createRegistryCmd.Flags().String("source", "", "Source URL for the private registry")
	createRegistryCmd.Flags().String("username", "", "Username for the private registry")
	createRegistryCmd.Flags().String("password", "", "Password for the private registry")

	createRegistryCmd.MarkFlagRequired("namespace")
	createRegistryCmd.MarkFlagRequired("name")
	createRegistryCmd.MarkFlagRequired("source")
	createRegistryCmd.MarkFlagRequired("username")
	createRegistryCmd.MarkFlagRequired("password")

	registryCmd.AddCommand(createRegistryCmd)
}
