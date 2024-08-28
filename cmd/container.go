package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"gitlab.com/Tilaa/tilaa-cli/api"

	"github.com/spf13/cobra"
)

var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "Manage containers",
}

var listContainersCmd = &cobra.Command{
	Use:   "list",
	Short: "List all containers",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")

		containers, err := api.ListContainers(namespace)

		if err != nil {
			log.Fatalf("Failed to list containers: %v", err)
		}

		if len(containers) == 0 {
			fmt.Println("No containers found.")
			return
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)

		fmt.Fprintln(writer, "NAME\t STATE\t IMAGE\t")

		for _, container := range containers {
			fmt.Fprintf(writer, "%s\t %s\t %s\t\n", container.Name, container.State, container.Image)
		}

		writer.Flush()
	},
}

func init() {
	listContainersCmd.Flags().StringP("namespace", "n", "", "Namespace")
	listContainersCmd.MarkFlagRequired("namespace")

	containerCmd.AddCommand(listContainersCmd)
}
