package cmd

import (
	"fmt"
	"log"

	"github.com/nexaa-cloud/nexaa-cli/api"
	"github.com/spf13/cobra"
)

var messageQueueEnableExternalConnectionCmd = &cobra.Command{
	Use:   "external-connection",
	Short: "Enable or Disable external connection on a cloud database cluster",
}

var enableMessageQueueExternalConnectionCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable external connection on a cloud database cluster",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		clusterName, _ := cmd.Flags().GetString("cluster")
		allowedIp, _ := cmd.Flags().GetStringArray("allowed-ip")

		allowList := make([]api.AllowListInput, 0)
		for _, ip := range allowedIp {
			allowList = append(allowList, api.AllowListInput{Ip: ip, State: api.StatePresent})
		}

		resource := api.MessageQueueModifyInput{
			Name:      clusterName,
			Namespace: namespace,
			ExternalConnection: &api.ExternalConnectionInput{
				State:    api.StatePresent,
				SharedIp: true,
				Ports: []api.ExternalConnectionPortInput{
					{
						AllowList: allowList,
						State:     api.StatePresent,
					},
				},
			},
		}

		client := api.NewClient()

		cluster, err := client.MessageQueueModify(resource)
		if err != nil {
			log.Fatalf("Failed to enable external connection in cluster %q/%q: %v", namespace, clusterName, err)
			return
		}
		fmt.Printf("External connection enabled. Reachable at:\n")
		fmt.Printf("Ipv4: %s:%d \n", cluster.ExternalConnection.Ipv4, cluster.ExternalConnection.Ports[0].ExternalPort)
		fmt.Printf("Ipv6: %s:%d \n", cluster.ExternalConnection.Ipv6, cluster.ExternalConnection.Ports[0].ExternalPort)
	},
}

var disableMessageQueueExternalConnectionCmd = &cobra.Command{
	Use:   "disable",
	Short: "disable external connection on a cloud database cluster",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		clusterName, _ := cmd.Flags().GetString("cluster")

		resource := api.MessageQueueModifyInput{
			Name:      clusterName,
			Namespace: namespace,
			ExternalConnection: &api.ExternalConnectionInput{
				State: api.StateAbsent,
				Ports: []api.ExternalConnectionPortInput{},
			},
		}

		client := api.NewClient()

		cluster, err := client.MessageQueueModify(resource)
		if err != nil {
			log.Fatalf("Failed to disable external connection in cluster %q/%q: %v", namespace, clusterName, err)
			return
		}
		fmt.Printf("External connection disabled in: %s/%s. \n", cluster.Namespace.Name, cluster.Name)
	},
}

func init() {
	enableMessageQueueExternalConnectionCmd.Flags().String("namespace", "", "Namespace")
	enableMessageQueueExternalConnectionCmd.Flags().String("cluster", "", "Name of the cluster")
	enableMessageQueueExternalConnectionCmd.Flags().StringArray("allowed-ip", []string{"0.0.0.0/0", "::/0"}, "Allowed ip for the connection")
	enableMessageQueueExternalConnectionCmd.MarkFlagRequired("namespace")
	enableMessageQueueExternalConnectionCmd.MarkFlagRequired("cluster")
	messageQueueEnableExternalConnectionCmd.AddCommand(enableMessageQueueExternalConnectionCmd)

	disableMessageQueueExternalConnectionCmd.Flags().String("namespace", "", "Namespace")
	disableMessageQueueExternalConnectionCmd.Flags().String("cluster", "", "Name of the cluster")
	disableMessageQueueExternalConnectionCmd.MarkFlagRequired("namespace")
	disableMessageQueueExternalConnectionCmd.MarkFlagRequired("cluster")
	messageQueueEnableExternalConnectionCmd.AddCommand(disableMessageQueueExternalConnectionCmd)
}
