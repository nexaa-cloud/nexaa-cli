package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/nexaa-cloud/nexaa-cli/api"
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
		client := api.NewClient()

		containers, err := client.ListContainers(namespace)

		if err != nil {
			log.Fatalf("Failed to list containers: %v", err)
		}

		if len(containers) == 0 {
			fmt.Println("No containers found.")
			return
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)

		fmt.Fprintln(writer, "NAME\t IMAGE\t RESOURCES\t")

		for _, container := range containers {
			fmt.Fprintf(writer, "%s\t %s\t %s\t\n", container.Name, container.Image, container.Resources)
		}

		writer.Flush()
	},
}

var createContainerCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new container",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")
		image, _ := cmd.Flags().GetString("image")
		resources, _ := cmd.Flags().GetString("resources")

		input := api.ContainerCreateInput{
			Name:                 name,
			Namespace:            namespace,
			Resources:            api.ContainerResources(resources),
			Image:                image,
			EnvironmentVariables: []api.EnvironmentVariableInput{},
			Mounts:               []api.MountInput{},
			Ports:                []string{},
			Ingresses:            []api.IngressInput{},
			HealthCheck:          &api.HealthCheckInput{},
		}

		client := api.NewClient()

		container, err := client.ContainerCreate(input)

		if err != nil {
			log.Fatalf("Failed to create container: %v", err)
			return
		}

		log.Println("Created container: ", container.Name)
	},
}

var modifyContainerCmd = &cobra.Command{
	Use:   "modify",
	Short: "Modify a container",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")
		image, _ := cmd.Flags().GetString("image")
		resources, _ := cmd.Flags().GetString("resources")

		input := api.ContainerModifyInput{
			Name:      name,
			Namespace: namespace,
		}

		if image != "" {
			input.Image = &image
		}

		if resources != "" {
			resource := api.ContainerResources(resources)
			input.Resources = &resource
		}

		client := api.NewClient()

		container, err := client.ContainerModify(input)
		if err != nil {
			log.Fatalf("Failed to modify container: %v", err)
			return
		}

		log.Println("Modified container: ", container.Name)
	},
}

var deleteContainerCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a container",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")

		client := api.NewClient()

		result, err := client.ContainerDelete(namespace, name)
		if err != nil {
			log.Fatalf("Failed to delete container: %v", err)
			return
		}

		if !result {
			log.Fatalf("Could not delete container with name: %s", name)
			return
		}
		log.Println("Deleted container with name: ", name)
	},
}

func init() {
	createContainerCmd.Flags().String("namespace", "", "Namespace")
	createContainerCmd.Flags().String("name", "", "Name for the container")
	createContainerCmd.Flags().String("image", "", "Container image")
	createContainerCmd.Flags().String("resources", "", "Container resources")
	createContainerCmd.MarkFlagRequired("namespace")
	createContainerCmd.MarkFlagRequired("name")
	createContainerCmd.MarkFlagRequired("image")
	createContainerCmd.MarkFlagRequired("resources")
	containerCmd.AddCommand(createContainerCmd)

	modifyContainerCmd.Flags().String("namespace", "", "Namespace")
	modifyContainerCmd.Flags().String("name", "", "Name for the container")
	modifyContainerCmd.Flags().String("image", "", "Container image")
	modifyContainerCmd.Flags().String("resources", "", "Container resources")
	modifyContainerCmd.MarkFlagRequired("namespace")
	modifyContainerCmd.MarkFlagRequired("name")
	containerCmd.AddCommand(modifyContainerCmd)

	listContainersCmd.Flags().String("namespace", "", "Namespace")
	listContainersCmd.MarkFlagRequired("namespace")
	containerCmd.AddCommand(listContainersCmd)

	deleteContainerCmd.Flags().String("namespace", "", "Namespace")
	deleteContainerCmd.Flags().String("name", "", "Name of this container")
	deleteContainerCmd.MarkFlagRequired("namespace")
	deleteContainerCmd.MarkFlagRequired("name")
	containerCmd.AddCommand(deleteContainerCmd)
}
