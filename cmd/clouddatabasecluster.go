package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/nexaa-cloud/nexaa-cli/api"
	"github.com/spf13/cobra"
)

var clouddatabaseclusterCmd = &cobra.Command{
	Use:     "databasecluster",
	Short:   "Manage cloud database clusters",
	Aliases: []string{"dc"},
}

var createCloudDatabaseClusterCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new cloud database cluster",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")
		planID, _ := cmd.Flags().GetString("plan")
		version, _ := cmd.Flags().GetString("version")
		dbType, _ := cmd.Flags().GetString("type")

		input := api.CloudDatabaseClusterCreateInput{
			Name:      name,
			Namespace: namespace,
			Plan:      planID,
			Spec: api.CloudDatabaseClusterSpecInput{
				Type:    dbType,
				Version: version,
			},
		}

		client := api.NewClient()
		result, err := client.CloudDatabaseClusterCreate(input)
		if err != nil {
			log.Fatalf("Failed to create cloud database cluster: %v", err)
			return
		}
		log.Println("Created cloud database cluster: ", result.Name)
	},
}

var listCloudDatabaseClustersCmd = &cobra.Command{
	Use:   "list",
	Short: "List all cloud database clusters",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient()
		clusters, err := client.CloudDatabaseClusterList()
		if err != nil {
			log.Fatalf("Failed to list cloud database clusters: %v", err)
		}
		if len(clusters) == 0 {
			fmt.Println("No cloud database clusters found.")
			return
		}
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
		fmt.Fprintln(writer, "NAME\tDATABASES\tNAMESPACE\tUSERS\t")
		for _, c := range clusters {
			dbCount := 0
			if c.Databases != nil {
				dbCount = len(c.Databases)
			}
			nsName := c.Namespace.Name
			userCount := 0
			if c.Users != nil {
				userCount = len(c.Users)
			}
			fmt.Fprintf(writer, "%s\t%d\t%s\t%d\t\n", c.Name, dbCount, nsName, userCount)
		}
		writer.Flush()
	},
}

func getDatabaseNames(databases []api.CloudDatabaseClusterResultDatabasesDatabase) string {
	names := make([]string, len(databases))
	for i, db := range databases {
		names[i] = db.Name
	}
	return strings.Join(names, ", ")
}

func getUsernames(users []api.CloudDatabaseClusterResultUsersDatabaseUser) string {
	usernames := make([]string, len(users))
	for i, user := range users {
		usernames[i] = user.Name
	}
	return strings.Join(usernames, ", ")
}

var listCloudDatabaseClusterCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a specific cloud database cluster",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")
		if namespace == "" || name == "" {
			log.Fatal("Namespace and name are required to get a cloud database cluster.")
			return
		}
		client := api.NewClient()
		input := api.CloudDatabaseClusterResourceInput{
			Namespace: namespace,
			Name:      name,
		}
		cluster, err := client.CloudDatabaseClusterGet(input)
		if err != nil {
			log.Fatalf("Failed to get cloud database cluster: %v", err)
			return
		}
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
		DatabaseNames := getDatabaseNames(cluster.Databases)
		Usernames := getUsernames(cluster.Users)
		fmt.Fprintln(writer, "NAME\tNAMESPACE\tPLAN\tTYPE\tVERSION\tDATABASES\tUSERS\t")
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t\n",
			cluster.Name, cluster.Namespace.Name, cluster.Plan.Id, cluster.Spec.Type, cluster.Spec.Version, DatabaseNames, Usernames)
		writer.Flush()
	},
}

var deleteCloudDatabaseClusterCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a cloud database cluster",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")
		client := api.NewClient()
		input := api.CloudDatabaseClusterResourceInput{
			Namespace: namespace,
			Name:      name,
		}
		result, err := client.CloudDatabaseClusterDelete(input)
		if err != nil {
			log.Fatalf("Failed to delete cloud database cluster: %v", err)
			return
		}
		if !result {
			log.Fatalf("Could not delete cloud database cluster with name: %q", name)
		}
		log.Println("Deleted cloud database cluster with name: ", name)
	},
}

var listCloudDatabaseClusterPlansCmd = &cobra.Command{
	Use:   "list-plans",
	Short: "List available plans for cloud database clusters",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient()
		plans, err := client.CloudDatabaseClusterListPlans()
		if err != nil {
			log.Fatalf("Failed to list cloud database cluster plans: %v", err)
		}
		if len(plans) == 0 {
			fmt.Println("No cloud database cluster plans found.")
			return
		}
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
		fmt.Fprintln(writer, "ID\tNAME\tCPU\tSTORAGE\tRAM\tCURRENCY\tPRICE\t")
		for _, p := range plans {
			fmt.Fprintf(writer, "%s\t%s\t%d\t%dGB\t%dGB\t%s\t%.2f\t\n", p.Id, p.Name, p.Cpu, p.Storage, int(p.Memory), *p.Price.Currency, float64(*p.Price.Amount)/100)
		}
		writer.Flush()
	},
}

var listCloudDatabaseClusterSpecsCmd = &cobra.Command{
	Use:   "list-specs",
	Short: "List available specs for cloud database clusters",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient()
		specs, err := client.CloudDatabaseClusterListSpecs()
		if err != nil {
			log.Fatalf("Failed to list cloud database cluster specs: %v", err)
		}
		if len(specs) == 0 {
			fmt.Println("No cloud database cluster specs found.")
			return
		}
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
		fmt.Fprintln(writer, "TYPE\tVERSION\t")
		for _, s := range specs {
			fmt.Fprintf(writer, "%s\t%s\t\n", s.Type, s.Version)
		}
		writer.Flush()
	},
}

var getClusterDatabaseUserCredentialsCmd = &cobra.Command{
	Use:   "get-credentials",
	Short: "Get user connection string for a database cluster",
	Run: func(cmd *cobra.Command, args []string) {
		clusterName, _ := cmd.Flags().GetString("cluster")
		namespace, _ := cmd.Flags().GetString("namespace")
		userName, _ := cmd.Flags().GetString("user")
		input := api.CloudDatabaseClusterResourceInput{
			Name:      clusterName,
			Namespace: namespace,
		}
		client := api.NewClient()
		dsn, err := client.CloudDatabaseClusterUserCredentials(input, userName)
		if err != nil {
			log.Fatalf("Failed to get user credentials: %v", err)
			return
		}
		fmt.Printf("DSN for user %s: %s\n", userName, dsn)
	},
}

func init() {
	createCloudDatabaseClusterCmd.Flags().StringP("namespace", "n", "", "Namespace")
	createCloudDatabaseClusterCmd.Flags().String("name", "", "Name for this cluster")
	createCloudDatabaseClusterCmd.Flags().String("plan", "", "ID of the plan to use for this cluster")
	createCloudDatabaseClusterCmd.Flags().String("type", "", "Type of the cluster (e.g., 'postgresql', 'mysql')")
	createCloudDatabaseClusterCmd.Flags().String("version", "", "Version of the database engine (e.g., '14', '15')")
	createCloudDatabaseClusterCmd.MarkFlagRequired("namespace")
	createCloudDatabaseClusterCmd.MarkFlagRequired("name")
	createCloudDatabaseClusterCmd.MarkFlagRequired("plan ID")
	createCloudDatabaseClusterCmd.MarkFlagRequired("type")
	createCloudDatabaseClusterCmd.MarkFlagRequired("version")
	clouddatabaseclusterCmd.AddCommand(createCloudDatabaseClusterCmd)

	clouddatabaseclusterCmd.AddCommand(listCloudDatabaseClustersCmd)

	deleteCloudDatabaseClusterCmd.Flags().String("namespace", "", "Namespace")
	deleteCloudDatabaseClusterCmd.Flags().String("name", "", "Name of the cluster")
	deleteCloudDatabaseClusterCmd.MarkFlagRequired("namespace")
	deleteCloudDatabaseClusterCmd.MarkFlagRequired("name")
	clouddatabaseclusterCmd.AddCommand(deleteCloudDatabaseClusterCmd)

	clouddatabaseclusterCmd.AddCommand(listCloudDatabaseClusterPlansCmd)
	clouddatabaseclusterCmd.AddCommand(listCloudDatabaseClusterSpecsCmd)

	listCloudDatabaseClusterCmd.Flags().StringP("namespace", "n", "", "Namespace name")
	listCloudDatabaseClusterCmd.Flags().String("name", "", "Name of the cluster")
	listCloudDatabaseClusterCmd.MarkFlagRequired("namespace")
	listCloudDatabaseClusterCmd.MarkFlagRequired("name")
	clouddatabaseclusterCmd.AddCommand(listCloudDatabaseClusterCmd)

	getClusterDatabaseUserCredentialsCmd.Flags().String("cluster", "", "Cluster name")
	getClusterDatabaseUserCredentialsCmd.Flags().StringP("namespace", "n", "", "Namespace")
	getClusterDatabaseUserCredentialsCmd.Flags().String("user", "", "User name")
	getClusterDatabaseUserCredentialsCmd.MarkFlagRequired("cluster")
	getClusterDatabaseUserCredentialsCmd.MarkFlagRequired("namespace")
	getClusterDatabaseUserCredentialsCmd.MarkFlagRequired("user")
	clouddatabaseclusterCmd.AddCommand(getClusterDatabaseUserCredentialsCmd)
}
