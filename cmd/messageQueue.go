package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/nexaa-cloud/nexaa-cli/api"
	"github.com/spf13/cobra"
)

var messageQueueCmd = &cobra.Command{
	Use:     "queue",
	Short:   "Manage message queues",
	Aliases: []string{"mq"},
}

var listMessageQueuesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all message queues",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient()
		queues, err := client.MessageQueueList()
		if err != nil {
			log.Fatalf("Failed to list message queues: %v", err)
		}
		if len(queues) == 0 {
			fmt.Println("No message queues found.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "NAME\t NAMESPACE\t STATE\t LOCKED\t")

		for _, queue := range queues {
			fmt.Fprintf(w, "%s\t %s\t %s\t %t\t\n", queue.Name, queue.Namespace.Name, queue.State, queue.Locked)
		}

		w.Flush()
	},
}

var getMessageQueueCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a message queue",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")
		client := api.NewClient()

		input := api.MessageQueueResourceInput{
			Name:      name,
			Namespace: namespace,
		}

		queue, err := client.MessageQueueGet(input)
		if err != nil {
			log.Fatalf("Failed to get message queue: %v", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "NAME\t NAMESPACE\t STATE\t LOCKED\t")
		fmt.Fprintf(w, "%s\t %s\t %s\t %t\t\n", queue.Name, queue.Namespace.Name, queue.State, queue.Locked)
		w.Flush()
	},
}

var createMessageQueueCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new message queue",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")
		plan, _ := cmd.Flags().GetString("plan")
		queueType, _ := cmd.Flags().GetString("type")
		version, _ := cmd.Flags().GetString("version")
		allowlistStr, _ := cmd.Flags().GetString("allowlist")

		// Parse allowlist
		var allowList []api.AllowListInput
		if allowlistStr != "" {
			ips := splitAndTrim(allowlistStr)
			for _, ip := range ips {
				allowList = append(allowList, api.AllowListInput{
					Ip:    ip,
					State: api.StatePresent,
				})
			}
		}

		input := api.MessageQueueCreateInput{
			Name:      name,
			Namespace: namespace,
			Plan:      plan,
			Spec: api.MessageQueueSpecInput{
				Type:    queueType,
				Version: version,
			},
			AllowList: allowList,
		}

		client := api.NewClient()
		queue, err := client.MessageQueueCreate(input)
		if err != nil {
			log.Fatalf("Failed to create message queue: %v", err)
		}

		log.Println("Created message queue:", queue.Name)
	},
}

var deleteMessageQueueCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a message queue",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("name")

		client := api.NewClient()

		input := api.ResourceNameInput{
			Name:      name,
			Namespace: namespace,
		}

		result, err := client.MessageQueueDelete(input)
		if err != nil {
			log.Fatalf("Failed to delete message queue: %v", err)
		}

		if !result {
			log.Fatalf("Could not delete message queue with name: %s", name)
		}
		log.Println("Deleted message queue with name:", name)
	},
}

var listMessageQueuePlansCmd = &cobra.Command{
	Use:   "plans",
	Short: "List all available message queue plans",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient()
		plans, err := client.MessageQueuePlans()
		if err != nil {
			log.Fatalf("Failed to list message queue plans: %v", err)
		}

		if len(plans) == 0 {
			fmt.Println("No message queue plans found.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "ID\t NAME\t GROUP\t CPU\t MEMORY (GB)\t REPLICAS\t STORAGE\t")

		for _, plan := range plans {
			fmt.Fprintf(w, "%s\t %s\t %s\t %d\t %.2f\t %d\t %d\t\n",
				plan.Id, plan.Name, plan.Group, plan.Cpu, plan.Memory, plan.Replicas, plan.Storage)
		}

		w.Flush()
	},
}

var listMessageQueueVersionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "List all available message queue versions",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient()
		versions, err := client.MessageQueueVersions()
		if err != nil {
			log.Fatalf("Failed to list message queue versions: %v", err)
		}

		if len(versions) == 0 {
			fmt.Println("No message queue versions found.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "TYPE\t VERSION\t PATCH LEVEL\t")

		for _, version := range versions {
			fmt.Fprintf(w, "%s\t %s\t %s\t\n", version.Type, version.Version, version.PatchLevelVersion)
		}

		w.Flush()
	},
}

func init() {
	// List command
	messageQueueCmd.AddCommand(listMessageQueuesCmd)

	// Get command
	getMessageQueueCmd.Flags().String("namespace", "", "Namespace")
	getMessageQueueCmd.Flags().String("name", "", "Name of the message queue")
	getMessageQueueCmd.MarkFlagRequired("namespace")
	getMessageQueueCmd.MarkFlagRequired("name")
	messageQueueCmd.AddCommand(getMessageQueueCmd)

	// Create command
	createMessageQueueCmd.Flags().String("namespace", "", "Namespace")
	createMessageQueueCmd.Flags().String("name", "", "Name for the message queue")
	createMessageQueueCmd.Flags().String("plan", "", "Plan ID for the message queue")
	createMessageQueueCmd.Flags().String("type", "", "Type of the message queue (e.g., RabbitMQ)")
	createMessageQueueCmd.Flags().String("version", "", "Version of the message queue")
	createMessageQueueCmd.Flags().String("allowlist", "", "Comma-separated list of IP addresses or CIDR ranges (e.g., 192.168.1.1,10.0.0.0/24)")
	createMessageQueueCmd.MarkFlagRequired("namespace")
	createMessageQueueCmd.MarkFlagRequired("name")
	createMessageQueueCmd.MarkFlagRequired("plan")
	createMessageQueueCmd.MarkFlagRequired("type")
	createMessageQueueCmd.MarkFlagRequired("version")
	messageQueueCmd.AddCommand(createMessageQueueCmd)

	// Delete command
	deleteMessageQueueCmd.Flags().String("namespace", "", "Namespace")
	deleteMessageQueueCmd.Flags().String("name", "", "Name of the message queue")
	deleteMessageQueueCmd.MarkFlagRequired("namespace")
	deleteMessageQueueCmd.MarkFlagRequired("name")
	messageQueueCmd.AddCommand(deleteMessageQueueCmd)

	// Plans command
	messageQueueCmd.AddCommand(listMessageQueuePlansCmd)

	// Versions command
	messageQueueCmd.AddCommand(listMessageQueueVersionsCmd)
}
