package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/nexaa-cloud/nexaa-cli/api"
	"github.com/spf13/cobra"
)

var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "Manage containers",
}

var getContainerCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a container",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")
		client := api.NewClient()

		container, err := client.ListContainerByName(namespace, name)
		if err != nil {
			log.Fatalf("Failed to get container : %v", err)
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)

		fmt.Fprintln(writer, "NAME\t IMAGE\t RESOURCES\t")

		fmt.Fprintf(writer, "%s\t %s\t %s\t\n", container.Name, container.Image, container.Resources)

		writer.Flush()
	},
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
		environmentVariables, _ := cmd.Flags().GetStringArray("env")
		secrets, _ := cmd.Flags().GetStringArray("secret")
		command, _ := cmd.Flags().GetStringArray("command")
		entrypoint, _ := cmd.Flags().GetStringArray("entrypoint")

		envs := append(envsToApi(environmentVariables, false, api.StatePresent), envsToApi(secrets, true, api.StatePresent)...)
		conResources := api.ContainerResources(resources)

		input := api.ContainerCreateInput{
			Name:                 name,
			Namespace:            namespace,
			Resources:            &conResources,
			Image:                &image,
			EnvironmentVariables: envs,
			Entrypoint:           entrypoint,
			Command:              command,
			Mounts:               []api.MountInput{},
			Ports:                []string{},
			Ingresses:            []api.IngressInput{},
			Type:                 api.ContainerTypeDefault,
		}

		client := api.NewClient()

		container, err := client.ContainerCreate(input)

		if err != nil {
			log.Fatalf("Failed to create container: %v", err)
			return
		}

		log.Println("Created container:", container.Name)
	},
}

var createStarterContainerCmd = &cobra.Command{
	Use:   "create-starter",
	Short: "Create a new starter container (CPU: 0.25, RAM: 0.5GB)",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")
		image, _ := cmd.Flags().GetString("image")
		environmentVariables, _ := cmd.Flags().GetStringArray("env")
		secrets, _ := cmd.Flags().GetStringArray("secret")
		command, _ := cmd.Flags().GetStringArray("command")
		entrypoint, _ := cmd.Flags().GetStringArray("entrypoint")

		envs := append(envsToApi(environmentVariables, false, api.StatePresent), envsToApi(secrets, true, api.StatePresent)...)
		resources := api.ContainerResourcesCpu250Ram500

		input := api.ContainerCreateInput{
			Name:                 name,
			Namespace:            namespace,
			Resources:            &resources,
			Image:                &image,
			EnvironmentVariables: envs,
			Entrypoint:           entrypoint,
			Command:              command,
			Mounts:               []api.MountInput{},
			Ports:                []string{},
			Ingresses:            []api.IngressInput{},
			Type:                 api.ContainerTypeStarter,
		}

		client := api.NewClient()

		container, err := client.ContainerCreate(input)

		if err != nil {
			log.Fatalf("Failed to create starter container: %v", err)
			return
		}

		log.Println("Created starter container:", container.Name)
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
		environmentVariables, _ := cmd.Flags().GetStringArray("env")
		secrets, _ := cmd.Flags().GetStringArray("secret")
		removedEnvironmentVariables, _ := cmd.Flags().GetStringArray("remove-env")
		registry, _ := cmd.Flags().GetString("registry")
		removeRegistry, _ := cmd.Flags().GetBool("remove-registry")
		command, _ := cmd.Flags().GetStringArray("command")
		removeCommand, _ := cmd.Flags().GetBool("remove-command")
		entrypoint, _ := cmd.Flags().GetStringArray("entrypoint")
		removeEntrypoint, _ := cmd.Flags().GetBool("remove-entrypoint")
		healthCheckPort, _ := cmd.Flags().GetInt("health-check-port")
		healthCheckPath, _ := cmd.Flags().GetString("health-check-path")
		removeHealthCheck, _ := cmd.Flags().GetBool("remove-health-check")

		client := api.NewClient()

		oldContainer, err := client.ListContainerByName(namespace, name)
		if err != nil {
			log.Fatalf("Container not found: %v", err)
		}

		input := buildContainerModifyInput(oldContainer, modifyContainerOptions{
			Namespace:                   namespace,
			Name:                        name,
			Image:                       image,
			Resources:                   resources,
			Registry:                    registry,
			RemoveRegistry:              removeRegistry,
			Command:                     command,
			RemoveCommand:               removeCommand,
			Entrypoint:                  entrypoint,
			RemoveEntrypoint:            removeEntrypoint,
			EnvironmentVariables:        environmentVariables,
			Secrets:                     secrets,
			RemovedEnvironmentVariables: removedEnvironmentVariables,
			HealthCheckPort:             healthCheckPort,
			HealthCheckPath:             healthCheckPath,
			RemoveHealthCheck:           removeHealthCheck,
		})

		container, err := client.ContainerModify(input)
		if err != nil {
			log.Fatalf("Failed to modify container: %v", err)
			return
		}

		log.Println("Modified container: ", container.Name)
	},
}

type modifyContainerOptions struct {
	Namespace                   string
	Name                        string
	Image                       string
	Resources                   string
	Registry                    string
	RemoveRegistry              bool
	Command                     []string
	RemoveCommand               bool
	Entrypoint                  []string
	RemoveEntrypoint            bool
	EnvironmentVariables        []string
	Secrets                     []string
	RemovedEnvironmentVariables []string
	HealthCheckPort             int
	HealthCheckPath             string
	RemoveHealthCheck           bool
}

func buildContainerModifyInput(oldContainer api.ContainerResult, opts modifyContainerOptions) api.ContainerModifyInput {
	envs := append(
		envsToApi(opts.EnvironmentVariables, false, api.StatePresent),
		envsToApi(opts.Secrets, true, api.StatePresent)...,
	)
	envs = append(
		envs,
		envsToApi(opts.RemovedEnvironmentVariables, false, api.StateAbsent)...,
	)

	input := api.ContainerModifyInput{
		Name:                 opts.Name,
		Namespace:            opts.Namespace,
		EnvironmentVariables: envs,
	}

	if !opts.RemoveRegistry {
		if opts.Registry != "" {
			input.Registry = &opts.Registry
		} else if oldContainer.PrivateRegistry != nil {
			input.Registry = &oldContainer.PrivateRegistry.Name
		}
	}

	if opts.Image != "" {
		input.Image = &opts.Image
	}

	if opts.Resources != "" {
		resource := api.ContainerResources(opts.Resources)
		input.Resources = &resource
	} else {
		resources := oldContainer.Resources
		input.Resources = &resources
	}

	if !opts.RemoveCommand {
		if len(opts.Command) > 0 {
			input.Command = opts.Command
		} else if len(oldContainer.Command) > 0 {
			input.Command = oldContainer.Command
		}
	}

	if !opts.RemoveEntrypoint {
		if len(opts.Entrypoint) > 0 {
			input.Entrypoint = opts.Entrypoint
		} else if len(oldContainer.Entrypoint) > 0 {
			input.Entrypoint = oldContainer.Entrypoint
		}
	}

	if !opts.RemoveHealthCheck {
		if opts.HealthCheckPort > 0 || opts.HealthCheckPath != "" {
			port := opts.HealthCheckPort
			path := opts.HealthCheckPath
			if oldContainer.HealthCheck != nil {
				if port == 0 {
					port = oldContainer.HealthCheck.Port
				}
				if path == "" {
					path = oldContainer.HealthCheck.Path
				}
			}
			input.HealthCheck = &api.HealthCheckInput{
				Port: port,
				Path: path,
			}
		} else if oldContainer.HealthCheck != nil {
			input.HealthCheck = &api.HealthCheckInput{
				Port: oldContainer.HealthCheck.Port,
				Path: oldContainer.HealthCheck.Path,
			}
		}
	}

	return input
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
	createContainerCmd.Flags().StringP("namespace", "n", "", "Namespace")
	createContainerCmd.Flags().String("name", "", "Name for the container")
	createContainerCmd.Flags().String("image", "", "Container image")
	createContainerCmd.Flags().String("resources", "", "Container resources")
	createContainerCmd.Flags().StringArray("env", []string{}, "Container environment variables")
	createContainerCmd.Flags().StringArray("secret", []string{}, "Container secrets")
	createContainerCmd.MarkFlagRequired("namespace")
	createContainerCmd.MarkFlagRequired("name")
	createContainerCmd.MarkFlagRequired("image")
	createContainerCmd.Flags().StringArray("entrypoint", []string{}, "Entrypoint for the container")
	createContainerCmd.Flags().StringArray("command", []string{}, "Command to run in the container")
	containerCmd.AddCommand(createContainerCmd)

	createStarterContainerCmd.Flags().StringP("namespace", "n", "", "Namespace")
	createStarterContainerCmd.Flags().String("name", "", "Name for the container")
	createStarterContainerCmd.Flags().String("image", "", "Container image")
	createStarterContainerCmd.Flags().StringArray("env", []string{}, "Container environment variables")
	createStarterContainerCmd.Flags().StringArray("secret", []string{}, "Container secrets")
	createStarterContainerCmd.Flags().StringArray("entrypoint", []string{}, "Entrypoint for the container")
	createStarterContainerCmd.Flags().StringArray("command", []string{}, "Command to run in the container")
	createStarterContainerCmd.MarkFlagRequired("namespace")
	createStarterContainerCmd.MarkFlagRequired("name")
	createStarterContainerCmd.MarkFlagRequired("image")
	containerCmd.AddCommand(createStarterContainerCmd)

	modifyContainerCmd.Flags().StringP("namespace", "n", "", "Namespace")
	modifyContainerCmd.Flags().String("name", "", "Name for the container")
	modifyContainerCmd.Flags().String("image", "", "Container image")
	modifyContainerCmd.Flags().String("resources", "", "Container resources")
	modifyContainerCmd.Flags().StringArray("command", []string{}, "Command to run in the container")
	modifyContainerCmd.Flags().Bool("remove-command", false, "Remove the command from the container")
	modifyContainerCmd.Flags().StringArray("entrypoint", []string{}, "Entrypoint for the container")
	modifyContainerCmd.Flags().Bool("remove-entrypoint", false, "Remove the entrypoint from the container")
	modifyContainerCmd.Flags().StringArray("env", []string{}, "Container environment variables")
	modifyContainerCmd.Flags().StringArray("secret", []string{}, "Container secrets")
	modifyContainerCmd.Flags().StringArray("remove-env", []string{}, "Container remove environment variables")
	modifyContainerCmd.Flags().String("registry", "", "Registry name for container image")
	modifyContainerCmd.Flags().Bool("remove-registry", false, "Change the registry to public")
	modifyContainerCmd.Flags().Int("health-check-port", 0, "Port for container health check")
	modifyContainerCmd.Flags().String("health-check-path", "", "Path for container health check")
	modifyContainerCmd.Flags().Bool("remove-health-check", false, "Remove health check from the container")
	modifyContainerCmd.MarkFlagRequired("namespace")
	modifyContainerCmd.MarkFlagRequired("name")
	containerCmd.AddCommand(modifyContainerCmd)

	listContainersCmd.Flags().StringP("namespace", "n", "", "Namespace")
	listContainersCmd.MarkFlagRequired("namespace")
	containerCmd.AddCommand(listContainersCmd)

	deleteContainerCmd.Flags().StringP("namespace", "n", "", "Namespace")
	deleteContainerCmd.Flags().String("name", "", "Name of this container")
	deleteContainerCmd.MarkFlagRequired("namespace")
	deleteContainerCmd.MarkFlagRequired("name")
	containerCmd.AddCommand(deleteContainerCmd)

	getContainerCmd.Flags().StringP("namespace", "n", "", "Namespace")
	getContainerCmd.Flags().String("name", "", "Name of the container")
	getContainerCmd.MarkFlagRequired("namespace")
	getContainerCmd.MarkFlagRequired("Name")
	containerCmd.AddCommand(getContainerCmd)

	containerCmd.AddCommand(containerEnableExternalConnectionCmd)
}
