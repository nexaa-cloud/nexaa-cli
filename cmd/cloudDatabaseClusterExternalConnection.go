package cmd

import (
	"fmt"
	"log"

	"github.com/nexaa-cloud/nexaa-cli/api"
	"github.com/spf13/cobra"
)

var cloudDatabaseClusterEnableExternalConnectionCmd = &cobra.Command{
	Use:     "external-connection",
	Short:   "Enable or Disable external connection on a cloud database cluster",
}

var enableCloudDatabaseClusterExternalConnectionCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable external connection on a cloud database cluster",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		clusterName, _ := cmd.Flags().GetString("cluster")
		allowedIp, _ := cmd.Flags().GetStringArray("allowed-ip")
		
		allowList := make([]api.AllowListInput, 0)
		for _, ip := range(allowedIp) {
			allowList = append(allowList, api.AllowListInput{Ip:ip})
		}

		resource := api.CloudDatabaseClusterModifyInput{
			Name:      clusterName,
			Namespace: namespace,
			ExternalConnection: &api.ExternalConnectionInput {
				SharedIp: true,
				Ports: []api.ExternalConnectionPortInput {
					api.ExternalConnectionPortInput{AllowList: allowList},
				} ,
			},
		}

		client := api.NewClient()

		cluster, err := client.CloudDatabaseClusterModify(resource)
		if err != nil {
			log.Fatalf("Failed to enable external connection in cluster %q/%q: %v", namespace, clusterName, err)
			return
		}
		fmt.Printf("External connection enabled. Reachable at:\n")
		fmt.Printf("Ipv4: %q:%q \n", cluster.ExternalConnection.Ipv4, cluster.ExternalConnection.Ports[0].ExternalPort)
		fmt.Printf("Ipv6: %q:%q \n", cluster.ExternalConnection.Ipv6, cluster.ExternalConnection.Ports[0].ExternalPort)
	},
}

var disableCloudDatabaseClusterExternalConnectionCmd = &cobra.Command{
	Use:   "disable",
	Short: "disable external connection on a cloud database cluster",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		clusterName, _ := cmd.Flags().GetString("cluster")
		
		resource := api.CloudDatabaseClusterModifyInput{
			Name:      clusterName,
			Namespace: namespace,
			ExternalConnection: &api.ExternalConnectionInput {
				State: api.StateAbsent,
			},
		}

		client := api.NewClient()

		cluster, err := client.CloudDatabaseClusterModify(resource)
		if err != nil {
			log.Fatalf("Failed to disable external connection in cluster %q/%q: %v", namespace, clusterName, err)
			return
		}
		fmt.Printf("External connection disabled in: %q/%q. \n", cluster.Namespace.Name, cluster.Name)
	},
}

func init() {
	enableCloudDatabaseClusterExternalConnectionCmd.Flags().String("namespace", "", "Namespace")
	enableCloudDatabaseClusterExternalConnectionCmd.Flags().String("cluster", "", "Name of the cluster")
	enableCloudDatabaseClusterExternalConnectionCmd.Flags().StringArray("allowed-ip", []string{"0.0.0.0/0", "::/0"}, "Allowed ip for the connection")
	enableCloudDatabaseClusterExternalConnectionCmd.MarkFlagRequired("namespace")
	enableCloudDatabaseClusterExternalConnectionCmd.MarkFlagRequired("cluster")
	cloudDatabaseClusterEnableExternalConnectionCmd.AddCommand(enableCloudDatabaseClusterExternalConnectionCmd)

	disableCloudDatabaseClusterExternalConnectionCmd.Flags().String("namespace", "", "Namespace")
	disableCloudDatabaseClusterExternalConnectionCmd.Flags().String("cluster", "", "Name of the cluster")
	disableCloudDatabaseClusterExternalConnectionCmd.MarkFlagRequired("namespace")
	disableCloudDatabaseClusterExternalConnectionCmd.MarkFlagRequired("cluster")	
	cloudDatabaseClusterEnableExternalConnectionCmd.AddCommand(disableCloudDatabaseClusterExternalConnectionCmd)	
}
