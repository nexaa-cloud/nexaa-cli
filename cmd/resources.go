package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/nexaa-cloud/nexaa-cli/api"
	"github.com/spf13/cobra"
)

var resourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "List Resources",
}

var listResourcesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all resources",
	Run: func(cmd *cobra.Command, args []string) {
		resources := api.AllContainerResources

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)

		fmt.Fprintln(writer, "RESOURCE\t")
		for _, resource := range resources {
			fmt.Fprintf(writer, "%s\t\n", resource)
		}
		writer.Flush()
	},
}

func init() {
	resourcesCmd.AddCommand(listResourcesCmd)
}
