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

		fmt.Fprintln(writer, "ID\tNAME\t STATE\t IMAGE\t")

		for _, container := range containers {
			fmt.Fprintf(writer, "%s\t%s\t %s\t %s\t\n", container.Id, container.Name, container.State, container.Image)
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
		httpPort, _ := cmd.Flags().GetInt("httpPort")
		https, _ := cmd.Flags().GetString("https")
		httpsPort, _ := cmd.Flags().GetInt("httpsPort")
		ports, _ := cmd.Flags().GetStringArray("ports")
		registry, _ := cmd.Flags().GetInt("registry")
		env, _ := cmd.Flags().GetStringArray("env")
		secret, _ := cmd.Flags().GetStringArray("secret")

		input := api.ContainerInput{
			Name:      name,
			Namespace: namespace,
			Image:     image,
			Http:      http,
			HttpPort:  httpPort,
			Https:     https,
			HttpsPort: httpsPort,
			Ports:     ports,
			Registry:  registry,
			Env:       env,
			Secret:    secret,
		}

		_, err := api.CreateContainer(input)

		if err != nil {
			log.Fatalf("Failed to create container: %v", err)
		}
	},
}

var modifyContainerCmd = &cobra.Command{
	Use:   "modify",
	Short: "Modify an existing container",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetInt("id")
		image, _ := cmd.Flags().GetString("image")
		http, _ := cmd.Flags().GetString("http")
		httpPort, _ := cmd.Flags().GetInt("httpPort")
		https, _ := cmd.Flags().GetString("https")
		httpsPort, _ := cmd.Flags().GetInt("httpsPort")
		ports, _ := cmd.Flags().GetStringArray("port")
		registry, _ := cmd.Flags().GetInt("registry")

		input := api.ContainerInput{
			Id:        id,
			Image:     image,
			Http:      http,
			HttpPort:  httpPort,
			Https:     https,
			HttpsPort: httpsPort,
			Ports:     ports,
			Registry:  registry,
		}

		_, err := api.ModifyContainer(input)

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
	createContainerCmd.Flags().String("name", "", "Name for this container")
	createContainerCmd.Flags().String("image", "", "Container image")
	createContainerCmd.Flags().String("http", "", "HTTP ingress hostname")
	createContainerCmd.Flags().Int("httpPort", 0, "HTTP ingress port")
	createContainerCmd.Flags().String("https", "", "HTTPS ingress hostname")
	createContainerCmd.Flags().Int("httpsPort", 0, "HTTPS ingress port")
	createContainerCmd.Flags().Int("registry", 0, "What registry to use")
	createContainerCmd.Flags().StringArrayP("env", "e", []string{}, "Environment variables")
	createContainerCmd.Flags().StringArrayP("secret", "s", []string{}, "Secret environment variables")
	createContainerCmd.Flags().StringArrayP("port", "p", []string{}, "Port mappings")
	createContainerCmd.MarkFlagRequired("namespace")
	createContainerCmd.MarkFlagRequired("name")
	createContainerCmd.MarkFlagRequired("image")
	containerCmd.AddCommand(createContainerCmd)

	modifyContainerCmd.Flags().IntP("id", "i", 0, "Container ID to modify")
	modifyContainerCmd.Flags().String("name", "", "Name for this container")
	modifyContainerCmd.Flags().String("image", "", "Container image")
	modifyContainerCmd.Flags().String("http", "", "HTTP ingress hostname")
	modifyContainerCmd.Flags().Int("httpPort", 0, "HTTP ingress port")
	modifyContainerCmd.Flags().String("https", "", "HTTPS ingress hostname")
	modifyContainerCmd.Flags().Int("httpsPort", 0, "HTTPS ingress port")
	modifyContainerCmd.Flags().Int("registry", 0, "What registry to use")
	modifyContainerCmd.Flags().StringArrayP("port", "p", []string{}, "Port mappings")
	modifyContainerCmd.MarkFlagRequired("id")
	containerCmd.AddCommand(modifyContainerCmd)
}
