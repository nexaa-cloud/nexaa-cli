package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/tilaa/tilaa-cli/api"
	"log"
	"os"
	"text/tabwriter"
)

var containerjobCmd = &cobra.Command{
	Use:   "containerjob",
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

		input := api.ContainerJobCreateInput{
			Name:                 name,
			Namespace:            namespace,
			Resources:            api.ContainerResources(resources),
			Image:                image,
			Entrypoint:           []string{},
			Command:              []string{},
			Enabled:              enabled,
			Schedule:             schedule,
			EnvironmentVariables: nil,
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
		enabled, _ := cmd.Flags().GetBool("enable")

		input := api.ContainerJobModifyInput{
			Name:      name,
			Namespace: namespace,
			Enabled:   &enabled,
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

		client := api.NewClient()

		containerJob, err := client.ContainerJobModify(input)

		if err != nil {
			log.Fatalf("Failed to modify container job: %v", err)
			return
		}

		log.Println("Modified container job: ", containerJob.Name)
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
			log.Fatalf("Failed to list containerjobs: %v", err)
		}

		if len(containerJobs) == 0 {
			fmt.Println("No containerjobs found.")
			return
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)

		fmt.Fprintln(writer, "NAME\t STATE\t IMAGE\t ENABLED\t")

		for _, container := range containerJobs {
			enabled := "False"
			if container.Enabled {
				enabled = "True"
			}

			fmt.Fprintf(writer, "%s\t %s\t %s\t %s\t\n", container.Name, container.State, container.Image, enabled)
		}

		writer.Flush()
	},
}

var deleteContainerJobCmd = &cobra.Command{
	Use: "delete",
	Short: "Delete a container job",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")

		client := api.NewClient()

		result, err := client.ContainerJobDelete(namespace, name)
		if err != nil {
			log.Fatalf("Failed to delete containerJob: %q", err)
			return
		}

		if !result {
			log.Fatalf("Could not delete containerJob with name: %q", name)
		}

		log.Println("deleted containerjob with name: ", name)
	},
}

func init() {
	createContainerJobCmd.Flags().String("namespace", "", "Namespace")
	createContainerJobCmd.Flags().String("name", "", "Name for this container")
	createContainerJobCmd.Flags().String("image", "", "Container image")
	createContainerJobCmd.Flags().String("resources", "", "Container resources")
	createContainerJobCmd.Flags().String("schedule", "", "Container schedule")
	createContainerJobCmd.Flags().Bool("enable", true, "enable container job")
	createContainerJobCmd.MarkFlagRequired("namespace")
	createContainerJobCmd.MarkFlagRequired("name")
	createContainerJobCmd.MarkFlagRequired("image")
	createContainerJobCmd.MarkFlagRequired("resources")
	createContainerJobCmd.MarkFlagRequired("schedule")
	containerjobCmd.AddCommand(createContainerJobCmd)

	modifyContainerJobCmd.Flags().String("namespace", "", "Namespace")
	modifyContainerJobCmd.Flags().String("name", "", "Name for this container")
	modifyContainerJobCmd.Flags().String("image", "", "Container image")
	modifyContainerJobCmd.Flags().String("resources", "", "Container resources")
	modifyContainerJobCmd.Flags().String("schedule", "", "Container schedule")
	modifyContainerJobCmd.Flags().Bool("enable", true, "enable container job")
	createContainerJobCmd.MarkFlagRequired("namespace")
	modifyContainerJobCmd.MarkFlagRequired("name")
	containerjobCmd.AddCommand(modifyContainerJobCmd)

	listContainerJobsCmd.Flags().String("namespace", "", "Namespace")
	listContainerJobsCmd.MarkFlagRequired("namespace")
	containerjobCmd.AddCommand(listContainerJobsCmd)

	deleteContainerJobCmd.Flags().String("namespace", "", "Namespace")
	deleteContainerJobCmd.Flags().String("name", "", "Name of the containerjob")
	deleteContainerJobCmd.MarkFlagRequired("namespace")
	deleteContainerJobCmd.MarkFlagRequired("name")
	containerjobCmd.AddCommand(deleteContainerJobCmd)
}
