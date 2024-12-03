package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"gitlab.com/Tilaa/tilaa-cli/api"

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
		namespaces, err := api.ListNamespaces()

		if err != nil {
			log.Fatalf("Failed to list namespaces: %v", err)
		}

		if len(namespaces) == 0 {
			fmt.Println("No namespaces found.")
			return
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)

		fmt.Fprintln(writer, "ID\t NAME\t")

		for _, namespace := range namespaces {
			fmt.Fprintf(writer, "%s\t %s\t\n", namespace.Id, namespace.Name)
		}

		writer.Flush()
	},
}

var createNamespaceCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new namespace",
	Run: func(cmd *cobra.Command, args []string) {
		customerId, _ := cmd.Flags().GetInt("customer_id")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		err := api.CreateNamespace(customerId, name, description)
		if err != nil {
			log.Fatalf("Failed to list namespaces: %v", err)
		}

		fmt.Println("Namespace created successfully.")
	},
}

var deleteNamespaceCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a namespace",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetInt("id")

		err := api.DeleteNamespace(id)
		if err != nil {
			log.Fatalf("Failed to list namespaces: %v", err)
		}

		fmt.Println("Namespace removed successfully.")
	},
}

func init() {
	namespaceCmd.AddCommand(listNamespacesCmd)

	createNamespaceCmd.Flags().StringP("name", "n", "", "Name")
	createNamespaceCmd.Flags().IntP("customer_id", "c", 0, "Your customer id")
	createNamespaceCmd.Flags().StringP("description", "d", "", "Description")
	createNamespaceCmd.MarkFlagRequired("name")
	createNamespaceCmd.MarkFlagRequired("customer_id")
	namespaceCmd.AddCommand(createNamespaceCmd)

	deleteNamespaceCmd.Flags().IntP("id", "i", 0, "Namespace id")
	deleteNamespaceCmd.MarkFlagRequired("id")
	namespaceCmd.AddCommand(deleteNamespaceCmd)
}
