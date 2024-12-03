package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"gitlab.com/Tilaa/tilaa-cli/api"
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

var createContainerCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new container",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetInt("namespace")
		name, _ := cmd.Flags().GetString("name")
		image, _ := cmd.Flags().GetString("image")
		http, _ := cmd.Flags().GetString("http")
		https, _ := cmd.Flags().GetString("https")
		ports, _ := cmd.Flags().GetStringArray("ports")

		input := api.ContainerInput{
			Name:      name,
			Namespace: namespace,
			Image:     image,
			Http:      http,
			Https:     https,
			Ports:     ports,
		}

		_, err := api.CreateContainer(input)

		if err != nil {
			log.Fatalf("Failed to create container: %v", err)
		}
	},
}

func init() {
	listContainersCmd.Flags().StringP("namespace", "n", "", "Namespace")
	listContainersCmd.MarkFlagRequired("namespace")
	containerCmd.AddCommand(listContainersCmd)

	createContainerCmd.Flags().IntP("namespace", "n", 0, "Namespace")
	createContainerCmd.Flags().String("name", "", "Namespace")
	createContainerCmd.Flags().String("image", "", "Container image")
	createContainerCmd.Flags().String("https", "", "HTTPS ingress hostname")
	createContainerCmd.Flags().String("http", "", "HTTP ingress hostname")
	createContainerCmd.Flags().StringArrayP("port", "p", []string{}, "Port mappings")
	createContainerCmd.MarkFlagRequired("namespace")
	createContainerCmd.MarkFlagRequired("name")
	createContainerCmd.MarkFlagRequired("image")
	containerCmd.AddCommand(createContainerCmd)
}
