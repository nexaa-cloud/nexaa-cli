package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/nexaa-cloud/nexaa-cli/api"
	"github.com/spf13/cobra"
)

var cloudDatabaseClusterDatabaseCmd = &cobra.Command{
	Use:     "database_cluster_database",
	Short:   "Manage database of a cloud database cluster",
	Aliases: []string{"dcd"},
}

var createCloudDatabaseClusterDatabaseCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new database inside a cloud database cluster",
	Run: func(cmd *cobra.Command, args []string) {
		namespaceName, _ := cmd.Flags().GetString("namespace")
		clusterName, _ := cmd.Flags().GetString("cluster")
		dbName, _ := cmd.Flags().GetString("name")
		dbDescription, _ := cmd.Flags().GetString("description")
		if dbDescription == "" {
			dbDescription = ""
		}
		client := api.NewClient()
		input := api.CloudDatabaseClusterDatabaseCreateInput{
			Cluster: api.CloudDatabaseClusterResourceInput{
				Name:      clusterName,
				Namespace: namespaceName,
			},
			Database: api.DatabaseInput{
				Name:        dbName,
				Description: &dbDescription,
				State:       api.StatePresent,
			},
		}
		result, err := client.CloudDatabaseClusterDatabaseCreate(input)
		if err != nil {
			log.Fatalf("Failed to create database: %v", err)
			return
		}
		fmt.Println("Created database: ", result.Name)
	},
}

var listCloudDatabaseClusterDatabasesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all databases in a cloud database cluster",
	Run: func(cmd *cobra.Command, args []string) {
		namespaceName, _ := cmd.Flags().GetString("namespace")
		clusterName, _ := cmd.Flags().GetString("cluster")
		client := api.NewClient()
		input := api.CloudDatabaseClusterResourceInput{
			Name:      clusterName,
			Namespace: namespaceName,
		}
		cluster, err := client.CloudDatabaseClusterDatabaseList(input)
		if err != nil {
			log.Fatalf("Failed to list cloud database clusters: %v", err)
		}
		if len(cluster.GetDatabases()) == 0 {
			fmt.Println("No databases found in cloud database cluster.")
			return
		}
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
		fmt.Fprintln(writer, "DATABASE NAME \tDESCRIPTION \t")
		for _, database := range cluster.GetDatabases() {
			if database.Description == nil {
				database.Description = new(string)
			}

			fmt.Fprintf(writer, "%s \t%s \t\n", database.Name, *database.Description)
		}
		writer.Flush()
	},
}

var deleteCloudDatabaseClusterDatabaseCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a database from cloud database cluster",
	Run: func(cmd *cobra.Command, args []string) {
		namespaceName, _ := cmd.Flags().GetString("namespace")
		clusterName, _ := cmd.Flags().GetString("cluster")
		name, _ := cmd.Flags().GetString("name")
		client := api.NewClient()
		input := api.CloudDatabaseClusterDatabaseResourceInput{
			Cluster: api.CloudDatabaseClusterResourceInput{
				Name:      clusterName,
				Namespace: namespaceName,
			},
			Name: name,
		}
		_, err := client.CloudDatabaseClusterDatabaseDelete(input)
		if err != nil {
			log.Fatalf("Failed to delete cloud database cluster: %v", err)
		}
		fmt.Println("Deleted cloud database cluster.")
	},
}

func init() {

	listCloudDatabaseClusterDatabasesCmd.Flags().String("namespace", "", "Name of the namespace cluster belongs to")
	listCloudDatabaseClusterDatabasesCmd.Flags().String("cluster", "", "Name of the cluster")
	listCloudDatabaseClusterDatabasesCmd.MarkFlagRequired("namespace")
	listCloudDatabaseClusterDatabasesCmd.MarkFlagRequired("cluster")
	cloudDatabaseClusterDatabaseCmd.AddCommand(listCloudDatabaseClusterDatabasesCmd)

	deleteCloudDatabaseClusterDatabaseCmd.Flags().String("namespace", "", "Name of the namespace cluster belongs to")
	deleteCloudDatabaseClusterDatabaseCmd.Flags().String("cluster", "", "Name of the cluster")
	deleteCloudDatabaseClusterDatabaseCmd.Flags().String("name", "", "Name of the database we want to delete")
	deleteCloudDatabaseClusterDatabaseCmd.MarkFlagRequired("namespace")
	deleteCloudDatabaseClusterDatabaseCmd.MarkFlagRequired("cluster")
	deleteCloudDatabaseClusterDatabaseCmd.MarkFlagRequired("name")
	cloudDatabaseClusterDatabaseCmd.AddCommand(deleteCloudDatabaseClusterDatabaseCmd)

	createCloudDatabaseClusterDatabaseCmd.Flags().String("namespace", "", "Name of namespace to create the cluster in")
	createCloudDatabaseClusterDatabaseCmd.Flags().String("cluster", "", "Name of the cluster")
	createCloudDatabaseClusterDatabaseCmd.Flags().String("name", "", "Name of the database we want to create")
	createCloudDatabaseClusterDatabaseCmd.Flags().String("description", "", "Description of the database")
	createCloudDatabaseClusterDatabaseCmd.MarkFlagRequired("namespace")
	createCloudDatabaseClusterDatabaseCmd.MarkFlagRequired("cluster")
	createCloudDatabaseClusterDatabaseCmd.MarkFlagRequired("name")
	cloudDatabaseClusterDatabaseCmd.AddCommand(createCloudDatabaseClusterDatabaseCmd)
}
