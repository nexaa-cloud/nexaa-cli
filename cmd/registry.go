package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"gitlab.com/tilaa/tilaa-cli/api"
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
		client := api.NewClient()

		registries, err := client.ListRegistries(namespace)

		if err != nil {
			log.Fatalf("Failed to list registries: %v", err)
		}

		if len(registries) == 0 {
			fmt.Println("No registries found.")
			return
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)

		fmt.Fprintln(writer, "NAME\t SOURCE\t USERNAME\t")

		for _, registry := range registries {
			fmt.Fprintf(writer, "%s\t %s\t %s\t\n", registry.Name, registry.Source, registry.Username)
		}

		writer.Flush()
	},
}

var createRegistryCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new private registry",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")
		source, _ := cmd.Flags().GetString("source")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		input := api.RegistryCreateInput{
			Namespace: namespace,
			Name:      name,
			Source:    source,
			Username:  username,
			Password:  password,
			Verify:    true,
		}

		client := api.NewClient()
		registry, err := client.RegistryCreate(input)

		if err != nil {
			log.Fatalf("Failed to create registry: %v", err)
		}

		fmt.Println("Created registry: ", registry.Name)
	},
}

var deleteRegistryCmd = &cobra.Command{
	Use: "delete",
	Short: "Delete a private registry",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")

		client := api.NewClient()

		result, err := client.RegistryDelete(namespace, name)
		if err != nil {
			log.Fatalf("Failed to delete registry: %q", err)
			return
		}

		if !result {
			log.Fatalf("Could not delete registry with name: %q", name)
			return
		}

		log.Println("deleted registry with name: ", name)
	},
}

func init() {
	listRegistriesCmd.Flags().StringP("namespace", "", "", "Namespace")
	listRegistriesCmd.MarkFlagRequired("namespace")
	registryCmd.AddCommand(listRegistriesCmd)

	createRegistryCmd.Flags().String("namespace", "", "Namespace")
	createRegistryCmd.Flags().String("name", "", "Name for the private registry")
	createRegistryCmd.Flags().String("source", "", "Source URL for the private registry")
	createRegistryCmd.Flags().String("username", "", "Username for the private registry")
	createRegistryCmd.Flags().String("password", "", "Password for the private registry")
	createRegistryCmd.MarkFlagRequired("namespace")
	createRegistryCmd.MarkFlagRequired("name")
	createRegistryCmd.MarkFlagRequired("source")
	createRegistryCmd.MarkFlagRequired("username")
	createRegistryCmd.MarkFlagRequired("password")
	registryCmd.AddCommand(createRegistryCmd)

	deleteRegistryCmd.Flags().String("namespace", "", "Namespace")
	deleteRegistryCmd.Flags().String("name", "", "Name of the private registry")
	deleteRegistryCmd.MarkFlagRequired("namespace")
	deleteRegistryCmd.MarkFlagRequired("name")
	registryCmd.AddCommand(deleteRegistryCmd)
}
