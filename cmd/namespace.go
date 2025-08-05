package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/nexaa-cloud/nexaa-cli/api"

	"github.com/spf13/cobra"
)

// Define the namespace command
var namespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "Manage namespaces",
}

var listNamespacesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all namespaces",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient()
		namespaces, err := client.NamespacesList()

		if err != nil {
			log.Fatalf("Failed to list namespaces: %v", err)
		}

		if len(namespaces) == 0 {
			fmt.Println("No namespaces found.")
			return
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)

		fmt.Fprintln(writer, "NAME\t DESCRIPTION\t")

		for _, namespace := range namespaces {
			fmt.Fprintf(writer, "%s\t %s\t\n", namespace.Name, namespace.Description)
		}

		writer.Flush()
	},
}

var createNamespaceCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new namespace",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := api.NamespaceCreateInput{
			Name:        name,
			Description: &description,
		}

		client := api.NewClient()

		namespace, err := client.NamespaceCreate(input)
		if err != nil {
			log.Fatalf("Failed to create namespace: %v", err)
		}

		fmt.Println("Created namespace: ", namespace.Name)
	},
}

var deleteNamespaceCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a namespace",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")

		client := api.NewClient()

		_, err := client.NamespaceDelete(name)
		if err != nil {
			log.Fatalf("Failed to delete namespace: %v", err)
		}

		fmt.Println("Deleted namespace with name: ", name)
	},
}

func init() {
	namespaceCmd.AddCommand(listNamespacesCmd)

	createNamespaceCmd.Flags().StringP("name", "n", "", "Name")
	createNamespaceCmd.Flags().StringP("description", "d", "", "Description")
	createNamespaceCmd.MarkFlagRequired("name")
	namespaceCmd.AddCommand(createNamespaceCmd)

	deleteNamespaceCmd.Flags().StringP("name", "n", "", "Name")
	deleteNamespaceCmd.MarkFlagRequired("name")
	namespaceCmd.AddCommand(deleteNamespaceCmd)
}
