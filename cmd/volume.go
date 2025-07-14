package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"gitlab.com/tilaa/tilaa-cli/api"
)

var volumeCmd = &cobra.Command{
	Use: "volume",
	Short: "Manage persistent volumes",
}

var listVolumesCmd = &cobra.Command{
	Use: "list",
	Short: "List all persistent volumes",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		client := api.NewClient()

		volumes, err := client.ListVolumes(namespace)
		if err != nil {
			log.Fatalf("Failed to list volumes: %v", err)
			return
		}

		if len(volumes) == 0 {
			fmt.Println("No volumes found.")
			return
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)

		fmt.Fprintln(writer, "NAME\t SIZE\t USAGE\t")

		for _, volume := range volumes {
			fmt.Fprintf(writer, "%s\t %f\t %f	\t", volume.Name, volume.Size, volume.Usage)
		}

		writer.Flush()
	},
}

var createVolumeCmd = &cobra.Command{
	Use: "create",
	Short: "Create a new persistent volume",
	Run: func (cmd *cobra.Command, args []string)  {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")
		size, _ := cmd.Flags().GetInt("size")
		
		input := api.VolumeCreateInput{
			Namespace: namespace,
			Name: name,
			Size: size,
		}

		client := api.NewClient()
		volume, err := client.VolumeCreate(input)
		if err != nil {
			log.Fatalf("Failed to create volume: %v", err)
			return
		}

		fmt.Println("created volume: ", volume.Name)
	},
}

var increaseVolumeCmd = &cobra.Command{
	Use: "increase",
	Short: "Increase the size of the volume",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")
		size, _ := cmd.Flags().GetInt("size")

		input := api.VolumeModifyInput{
			Namespace: namespace,
			Name: name,
			Size: size,
		}

		client := api.NewClient()

		volume, err := client.VolumeIncrease(input)
		if err != nil {
			log.Fatalf("Failed to increase volume: %s", err)
			return
		}

		log.Println("increased volume: ", volume.Name)
	},
}

var deleteVolumeCmd = &cobra.Command{
	Use: "delete",
	Short: "Delete a persistent volume",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")

		client := api.NewClient()

		result, err := client.VolumeDelete(namespace, name)
		if err != nil {
			log.Fatalf("Failed to delete volume: %s", err)
			return
		}

		if !result {
			log.Fatalf("Could not delete volume with name: %s", name)
			return
		}
		log.Println("deleted volume with name: ", name)
	},
}


func init() {
	listVolumesCmd.Flags().String("namespace", "", "Namespace")
	listVolumesCmd.MarkFlagRequired("namespace")
	volumeCmd.AddCommand(listVolumesCmd)

	createVolumeCmd.Flags().String("namespace", "", "Namespace")
	createVolumeCmd.Flags().String("name", "", "Name for the volume")
	createVolumeCmd.Flags().Int("size", 0, "Size of the volume")
	createVolumeCmd.MarkFlagRequired("namespace")
	createVolumeCmd.MarkFlagRequired("name")
	createVolumeCmd.MarkFlagRequired("size")
	volumeCmd.AddCommand(createVolumeCmd)

	increaseVolumeCmd.Flags().String("namespace", "", "Namespace")
	increaseVolumeCmd.Flags().String("name", "", "Name of the volume")
	increaseVolumeCmd.Flags().Int("size", 0, "Size of the volume")
	increaseVolumeCmd.MarkFlagRequired("namespace")
	increaseVolumeCmd.MarkFlagRequired("name")
	increaseVolumeCmd.MarkFlagRequired("size")
	volumeCmd.AddCommand(increaseVolumeCmd)

	deleteVolumeCmd.Flags().String("namespace", "", "Namespace")
	deleteVolumeCmd.Flags().String("name", "", "name")
	deleteVolumeCmd.MarkFlagRequired("namespace")
	deleteVolumeCmd.MarkFlagRequired("name")
	volumeCmd.AddCommand(deleteVolumeCmd)
}