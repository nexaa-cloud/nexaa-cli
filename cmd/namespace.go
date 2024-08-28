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

func init() {
	namespaceCmd.AddCommand(listNamespacesCmd)

	rootCmd.AddCommand(namespaceCmd)
}
