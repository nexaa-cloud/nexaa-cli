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

var containerEnableExternalConnectionCmd = &cobra.Command{
	Use:   "external-connection",
	Short: "Enable or Disable external connection on a cloud database cluster",
}

var enableContainerExternalConnectionCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable external connection on a cloud database cluster",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		containerName, _ := cmd.Flags().GetString("name")
		allowedIp, _ := cmd.Flags().GetStringArray("allowed-ip")
		internalPort, _ := cmd.Flags().GetInt("internal-port")
		externalPort, _ := cmd.Flags().GetInt("external-port")
		protocol, _ := cmd.Flags().GetString("protocol")

		if internalPort == 0 {
			log.Fatalf("Internal port must be provided and cannot be 0")
			return
		}

		allowList := make([]api.AllowListInput, 0)
		for _, ip := range allowedIp {
			allowList = append(allowList, api.AllowListInput{Ip: ip, State: api.StatePresent})
		}

		var resource api.ContainerModifyInput

		if externalPort == 0 {
			resource = api.ContainerModifyInput{
				Name:      containerName,
				Namespace: namespace,
				ExternalConnection: &api.ExternalConnectionInput{
					State:    api.StatePresent,
					SharedIp: true,
					Ports: []api.ExternalConnectionPortInput{
						{
							AllowList:    allowList,
							State:        api.StatePresent,
							InternalPort: &internalPort,
							Protocol:     api.Protocol(protocol),
						},
					},
				},
			}
		} else {
			resource = api.ContainerModifyInput{
				Name:      containerName,
				Namespace: namespace,
				ExternalConnection: &api.ExternalConnectionInput{
					State:    api.StatePresent,
					SharedIp: true,
					Ports: []api.ExternalConnectionPortInput{
						{
							AllowList:    allowList,
							State:        api.StatePresent,
							InternalPort: &internalPort,
							ExternalPort: &externalPort,
							Protocol:     api.Protocol(protocol),
						},
					},
				},
			}
		}

		client := api.NewClient()

		container, err := client.ContainerModify(resource)
		if err != nil {
			log.Fatalf("Failed to enable external connection in cluster %q/%q: %v", namespace, containerName, err)
			return
		}

		fmt.Printf("External connection enabled.\n")
		printConnections(container)
	},
}

var disableContainerExternalConnectionCmd = &cobra.Command{
	Use:   "disable",
	Short: "disable external connection on a container",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")
		externalPort, _ := cmd.Flags().GetInt("external-port")

		client := api.NewClient()

		// If externalPort is not provided, it will disable all external connections on the container
		var resource api.ContainerModifyInput
		if externalPort == 0 {
			resource = api.ContainerModifyInput{
				Name:      name,
				Namespace: namespace,
				ExternalConnection: &api.ExternalConnectionInput{
					State: api.StateAbsent,
					Ports: []api.ExternalConnectionPortInput{},
				},
			}
		} else {
			ports := []api.ExternalConnectionPortInput{}

			ports = append(ports, api.ExternalConnectionPortInput{
				ExternalPort: &externalPort,
				State:        api.StateAbsent,
				Protocol:     api.ProtocolUdp,
				AllowList:    []api.AllowListInput{},
			})

			resource = api.ContainerModifyInput{
				Name:      name,
				Namespace: namespace,
				ExternalConnection: &api.ExternalConnectionInput{
					State: api.StatePresent,
					Ports: ports,
				},
			}
		}

		container, err := client.ContainerModify(resource)

		if err != nil {
			log.Fatalf("Failed to disable external connection in cluster %q/%q: %v", namespace, name, err)
			return
		}
		fmt.Printf("External connection disabled in: %s/%s. \n", namespace, name)

		printConnections(container)
	},
}

var listContainerExternalConnectionCmd = &cobra.Command{
	Use:   "list",
	Short: "list external connection on a cloud database cluster",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")

		client := api.NewClient()
		container, err := client.ListContainerByName(namespace, name)

		if err != nil {
			log.Fatalf("Container %q/%q not found: %v", namespace, name, err)
			return
		}

		if container.ExternalConnection == nil || len(container.ExternalConnection.Ports) == 0 {
			log.Printf("No external connections enabled on %q/%q.", namespace, name)
			return
		}

		printConnections(container)
	},
}

func printConnections(container api.ContainerResult) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
	fmt.Fprintln(writer, "IPV4\t IPV6\t EXTERNAL PORT\t INTERNAL PORT\t PROTOCOL\t ALLOWLIST\t")

	for _, port := range container.ExternalConnection.Ports {
		fmt.Fprintf(
			writer,
			"%s\t %s\t %d\t %d\t %s\t %s\n",
			container.ExternalConnection.Ipv4,
			container.ExternalConnection.Ipv4,
			port.ExternalPort,
			*port.InternalPort,
			port.Protocol,
			strings.Join(port.AllowList, ","),
		)
	}

	writer.Flush()
}

func init() {
	enableContainerExternalConnectionCmd.Flags().StringP("namespace", "n", "", "Namespace")
	enableContainerExternalConnectionCmd.Flags().String("name", "", "Name of the container")
	enableContainerExternalConnectionCmd.Flags().StringArray("allowed-ip", []string{"0.0.0.0/0", "::/0"}, "Allowed ip for the connection")
	enableContainerExternalConnectionCmd.Flags().Int("internal-port", 0, "Internal port to enable external connection on")
	enableContainerExternalConnectionCmd.Flags().Int("external-port", 0, "Internal port to enable external connection on")
	enableContainerExternalConnectionCmd.Flags().String("protocol", "TCP", "Protocol for the external connection (tcp or udp)")
	enableContainerExternalConnectionCmd.MarkFlagRequired("namespace")
	enableContainerExternalConnectionCmd.MarkFlagRequired("name")
	enableContainerExternalConnectionCmd.MarkFlagRequired("internal-port")
	containerEnableExternalConnectionCmd.AddCommand(enableContainerExternalConnectionCmd)

	disableContainerExternalConnectionCmd.Flags().StringP("namespace", "n", "", "Namespace")
	disableContainerExternalConnectionCmd.Flags().String("name", "", "Name of the container")
	disableContainerExternalConnectionCmd.Flags().Int("external-port", 0, "External port to disable")
	disableContainerExternalConnectionCmd.MarkFlagRequired("namespace")
	disableContainerExternalConnectionCmd.MarkFlagRequired("name")
	containerEnableExternalConnectionCmd.AddCommand(disableContainerExternalConnectionCmd)

	listContainerExternalConnectionCmd.Flags().StringP("namespace", "n", "", "Namespace")
	listContainerExternalConnectionCmd.Flags().String("name", "", "Name of the container")
	listContainerExternalConnectionCmd.MarkFlagRequired("namespace")
	listContainerExternalConnectionCmd.MarkFlagRequired("name")
	containerEnableExternalConnectionCmd.AddCommand(listContainerExternalConnectionCmd)

}
