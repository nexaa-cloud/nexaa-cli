package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/nexaa-cloud/nexaa-cli/api"
	"github.com/spf13/cobra"
)

var containerJobCmd = &cobra.Command{
	Use:   "container_job",
	Short: "Manage container jobs",
}

var createContainerJobCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new container job",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")
		image, _ := cmd.Flags().GetString("image")
		resources, _ := cmd.Flags().GetString("resources")
		schedule, _ := cmd.Flags().GetString("schedule")
		enabled, _ := cmd.Flags().GetBool("enable")
		command, _ := cmd.Flags().GetStringArray("command")
		entrypoint, _ := cmd.Flags().GetStringArray("entrypoint")
		environmentVariables, _ := cmd.Flags().GetStringArray("env")
		secrets, _ := cmd.Flags().GetStringArray("secret")

		envs := append(envsToApi(environmentVariables, false, api.StatePresent), envsToApi(secrets, true, api.StatePresent)...)

		input := api.ContainerJobCreateInput{
			Name:                 name,
			Namespace:            namespace,
			Resources:            api.ContainerResources(resources),
			Image:                image,
			Entrypoint:           entrypoint,
			Command:              command,
			Enabled:              enabled,
			Schedule:             schedule,
			EnvironmentVariables: envs,
			Mounts:               []api.MountInput{},
		}

		client := api.NewClient()

		containerJob, err := client.ContainerJobCreate(input)

		if err != nil {
			log.Fatalf("Failed to create container job: %v", err)
			return
		}

		log.Println("Created container job: ", containerJob.Name)
	},
}

var modifyContainerJobCmd = &cobra.Command{
	Use:   "modify",
	Short: "Modify a container job",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")
		image, _ := cmd.Flags().GetString("image")
		resources, _ := cmd.Flags().GetString("resources")
		schedule, _ := cmd.Flags().GetString("schedule")
		command, _ := cmd.Flags().GetStringArray("command")
		entrypoint, _ := cmd.Flags().GetStringArray("entrypoint")
		enabled, _ := cmd.Flags().GetBool("enable")
		environmentVariables, _ := cmd.Flags().GetStringArray("env")
		secrets, _ := cmd.Flags().GetStringArray("secret")
		removedEnvironmentVariables, _ := cmd.Flags().GetStringArray("remove-env")
		envs := append(
			envsToApi(environmentVariables, false, api.StatePresent),
			envsToApi(secrets, true, api.StatePresent)...,
		)
		envs = append(
			envs,
			envsToApi(removedEnvironmentVariables, false, api.StateAbsent)...,
		)

		input := api.ContainerJobModifyInput{
			Name:                 name,
			Namespace:            namespace,
			Enabled:              &enabled,
			EnvironmentVariables: envs,
		}

		if image != "" {
			input.Image = &image
		}

		if resources != "" {
			resources := api.ContainerResources(resources)
			input.Resources = &resources
		}

		if schedule != "" {
			input.Schedule = &schedule
		}

		if len(command) > 0 {
			input.Command = command
		}

		if len(entrypoint) > 0 {
			input.Entrypoint = entrypoint
		}

		client := api.NewClient()

		containerJob, err := client.ContainerJobModify(input)

		if err != nil {
			log.Fatalf("Failed to modify container job: %v", err)
			return
		}

		log.Println("Modified container job: ", containerJob.Name)
	},
}

var getContainerJobCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a container job",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")
		client := api.NewClient()

		container, err := client.ContainerJobByName(namespace, name)
		if err != nil {
			log.Fatalf("Failed to list container jobs: %v", err)
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)

		fmt.Fprintln(writer, "NAME\t STATE\t IMAGE\t ENTRYPOINT\t COMMAND\t ENABLED\t")

		fmt.Fprintf(writer, "%s\t %s\t %s\t %s\t %s\t %s\t\n", container.Name, container.State, container.Image, commandApiToString(container.Entrypoint), commandApiToString(container.Command), enabledApiToString(container.Enabled))

		writer.Flush()
	},
}

var listContainerJobsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all container jobs",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		client := api.NewClient()

		containerJobs, err := client.ContainerJobList(namespace)

		if err != nil {
			log.Fatalf("Failed to list container jobs: %v", err)
		}

		if len(containerJobs) == 0 {
			fmt.Println("No containerjobs found.")
			return
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)

		fmt.Fprintln(writer, "NAME\t STATE\t IMAGE\t ENTRYPOINT\t COMMAND\t ENABLED\t")

		for _, container := range containerJobs {
			// Skip containers with empty names to avoid having a FALSE enabled empty row
			if container.Name == "" {
				continue
			}

			fmt.Fprintf(writer, "%s\t %s\t %s\t %s\t %s\t %s\t\n", container.Name, container.State, container.Image, commandApiToString(container.Entrypoint), commandApiToString(container.Command), enabledApiToString(container.Enabled))
		}

		writer.Flush()
	},
}

var deleteContainerJobCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a container job",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")

		client := api.NewClient()

		result, err := client.ContainerJobDelete(namespace, name)
		if err != nil {
			log.Fatalf("Failed to delete container job: %q", err)
			return
		}

		if !result {
			log.Fatalf("Could not delete container job with name: %q", name)
		}

		log.Println("deleted container job with name: ", name)
	},
}

func init() {
	createContainerJobCmd.Flags().String("namespace", "", "Namespace")
	createContainerJobCmd.Flags().String("name", "", "Name for this container job")
	createContainerJobCmd.Flags().String("image", "", "Container job image")
	createContainerJobCmd.Flags().String("resources", "", "Container job resources")
	createContainerJobCmd.Flags().String("schedule", "", "Container job schedule")
	createContainerJobCmd.Flags().StringArray("env", []string{}, "Container job environment variables")
	createContainerJobCmd.Flags().StringArray("secret", []string{}, "Container job secrets")
	createContainerJobCmd.Flags().Bool("enable", true, "enable container job")
	createContainerJobCmd.MarkFlagRequired("namespace")
	createContainerJobCmd.MarkFlagRequired("name")
	createContainerJobCmd.MarkFlagRequired("image")
	createContainerJobCmd.MarkFlagRequired("resources")
	createContainerJobCmd.MarkFlagRequired("schedule")
	containerJobCmd.AddCommand(createContainerJobCmd)

	modifyContainerJobCmd.Flags().String("namespace", "", "Namespace")
	modifyContainerJobCmd.Flags().String("name", "", "Name for this container job")
	modifyContainerJobCmd.Flags().String("image", "", "Container job image")
	modifyContainerJobCmd.Flags().String("resources", "", "Container job resources")
	modifyContainerJobCmd.Flags().String("schedule", "", "Container job schedule")
	modifyContainerJobCmd.Flags().Bool("enable", true, "enable container job")
	modifyContainerJobCmd.Flags().StringArray("env", []string{}, "Container job environment variables")
	modifyContainerJobCmd.Flags().StringArray("remove-env", []string{}, "Container job remove environment variables")
	modifyContainerJobCmd.Flags().StringArray("secret", []string{}, "Container job secrets")
	modifyContainerJobCmd.MarkFlagRequired("namespace")
	modifyContainerJobCmd.MarkFlagRequired("name")
	containerJobCmd.AddCommand(modifyContainerJobCmd)

	listContainerJobsCmd.Flags().String("namespace", "", "Namespace")
	listContainerJobsCmd.MarkFlagRequired("namespace")
	containerJobCmd.AddCommand(listContainerJobsCmd)

	deleteContainerJobCmd.Flags().String("namespace", "", "Namespace")
	deleteContainerJobCmd.Flags().String("name", "", "Name of the containerjob")
	deleteContainerJobCmd.MarkFlagRequired("namespace")
	deleteContainerJobCmd.MarkFlagRequired("name")
	containerJobCmd.AddCommand(deleteContainerJobCmd)

	getContainerJobCmd.Flags().String("namespace", "", "Namespace")
	getContainerJobCmd.Flags().String("name", "", "Name of the container job")
	getContainerJobCmd.MarkFlagRequired("namespace")
	getContainerJobCmd.MarkFlagRequired("Name")
	containerJobCmd.AddCommand(getContainerJobCmd)
}
