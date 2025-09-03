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

type destructedUser struct {
	Name         string
	DatabaseName string
	Permission   string
}

var cloudDatabaseClusterUserCmd = &cobra.Command{
	Use:     "database_cluster_user",
	Short:   "Manage user of a cloud database cluster",
	Aliases: []string{"dcu"},
}

var createCloudDatabaseClusterUserCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user inside a cloud database cluster",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		clusterName, _ := cmd.Flags().GetString("cluster")
		userName, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		permPairs, _ := cmd.Flags().GetStringSlice("permission")

		resource := api.CloudDatabaseClusterResourceInput{
			Name:      clusterName,
			Namespace: namespace,
		}

		client := api.NewClient()

		permsissions := make([]api.DatabaseUserPermissionInput, 0)
		for _, p := range permPairs {

			parts := strings.SplitN(p, ":", 2)
			if len(parts) != 2 {
				log.Fatalf("Invalid permission format %q. Use database:permission, e.g., mydb:read_write", p)
				return
			}
			name := parts[0]
			Permission := parts[1]

			permsissions = append(permsissions, api.DatabaseUserPermissionInput{
				DatabaseName: name,
				Permission:   api.DatabasePermission(Permission),
				State:        api.StatePresent,
			})
		}

		input := api.CloudDatabaseClusterUserCreateInput{
			Cluster: resource,
			User: api.DatabaseUserInput{
				Name:        userName,
				Password:    &password,
				State:       api.StatePresent,
				Permissions: permsissions,
			},
		}

		user, err := client.CloudDatabaseClusterUserCreate(input)
		if err != nil {
			log.Fatalf("Failed to create user %q in cluster %q/%q: %v", userName, namespace, clusterName, err)
			return
		}
		fmt.Printf("User %q created.\n", user.Name)
	},
}

var modifyCloudDatabaseClusterUserCmd = &cobra.Command{
	Use:   "modify",
	Short: "Modify a new user inside a cloud database cluster",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		clusterName, _ := cmd.Flags().GetString("cluster")
		userName, _ := cmd.Flags().GetString("user")
		newPassword, _ := cmd.Flags().GetString("password")
		addPerms, _ := cmd.Flags().GetStringSlice("add-permission")
		removePerms, _ := cmd.Flags().GetStringSlice("remove-permission")

		resource := api.CloudDatabaseClusterResourceInput{
			Name:      clusterName,
			Namespace: namespace,
		}

		client := api.NewClient()

		var parsed []api.DatabaseUserPermissionInput
		for _, p := range addPerms {
			parts := strings.SplitN(p, ":", 2)
			if len(parts) != 2 {
				log.Fatalf("Invalid add-permission format %q. Use database:permission", p)
				return
			}
			name := parts[0]
			var databasePermission api.DatabasePermission
			if parts[1] == "read_write" {
				databasePermission = api.DatabasePermissionReadWrite
			} else {
				databasePermission = api.DatabasePermissionReadOnly
			}
			parsed = append(parsed, api.DatabaseUserPermissionInput{DatabaseName: name, Permission: databasePermission, State: api.StatePresent})
		}
		for _, p := range removePerms {
			// remove only needs database name; permission is ignored if not provided
			db := p
			if strings.Contains(p, ":") {
				db = strings.SplitN(p, ":", 2)[0]
			}
			parsed = append(parsed, api.DatabaseUserPermissionInput{DatabaseName: db, Permission: "", State: api.StateAbsent})
		}

		UserInput := api.DatabaseUserInput{
			Name:        userName,
			Permissions: parsed,
			Password:    &newPassword,
			State:       api.StatePresent,
		}

		input := api.CloudDatabaseClusterUserModifyInput{
			Cluster: &resource,
			User:    &UserInput,
		}

		user, err := client.CloudDatabaseClusterUserModify(input)
		if err != nil {
			log.Fatalf("Failed to modify user %q in cluster %q/%q: %v", userName, namespace, clusterName, err)
			return
		}
		fmt.Printf("User %q updated.\n", user.Name)
	},
}

var listCloudDatabaseClusterUserCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users in a cloud database cluster",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		name, _ := cmd.Flags().GetString("cluster")
		resource := api.CloudDatabaseClusterResourceInput{
			Name:      name,
			Namespace: namespace,
		}

		client := api.NewClient()
		users, err := client.CloudDatabaseClusterUserList(resource)
		if err != nil {
			log.Fatalf("Failed to list cloud database cluster users: %v", err)
		}
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)

		var destructedUsers []destructedUser
		for _, user := range users {
			if len(user.Permissions) == 0 {
				var permission string
				if user.Role == "admin" {
					permission = "Database admin"
				} else {
					permission = "No Permission"
				}
				destructedUsers = append(destructedUsers, destructedUser{
					Name:         user.Name,
					DatabaseName: "",
					Permission:   permission,
				})
				continue
			}
			for _, permission := range user.Permissions {
				destructedUsers = append(destructedUsers, destructedUser{
					Name:         user.Name,
					DatabaseName: permission.DatabaseName,
					Permission:   string(permission.GetPermission()),
				})
			}
		}

		fmt.Fprintln(writer, "NAME\tDATABASES\tPERMISSION")
		for _, user := range destructedUsers {
			fmt.Fprintf(writer, "%s\t%s\t%s\n", user.Name, user.DatabaseName, user.Permission)
		}
		writer.Flush()
	},
}

var deleteCloudDatabaseClusterUserCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a user from cloud database cluster",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		clusterName, _ := cmd.Flags().GetString("cluster")
		userName, _ := cmd.Flags().GetString("user")

		cluster := api.CloudDatabaseClusterResourceInput{
			Name:      clusterName,
			Namespace: namespace,
		}

		client := api.NewClient()
		_, err := client.CloudDatabaseClusterUserDelete(api.CloudDatabaseClusterUserResourceInput{
			Name:    userName,
			Cluster: cluster,
		})
		if err != nil {
			log.Fatalf("Failed to delete user %q from cluster %q/%q: %v", userName, namespace, clusterName, err)
			return
		}
		fmt.Printf("User %q deleted.\n", userName)
	},
}

func init() {
	// list
	listCloudDatabaseClusterUserCmd.Flags().StringP("namespace", "n", "", "Namespace name")
	listCloudDatabaseClusterUserCmd.Flags().String("cluster", "", "Name of the cluster")
	listCloudDatabaseClusterUserCmd.MarkFlagRequired("namespace")
	listCloudDatabaseClusterUserCmd.MarkFlagRequired("cluster")
	cloudDatabaseClusterUserCmd.AddCommand(listCloudDatabaseClusterUserCmd)

	// create
	createCloudDatabaseClusterUserCmd.Flags().StringP("namespace", "n", "", "Namespace name")
	createCloudDatabaseClusterUserCmd.Flags().String("cluster", "", "Name of the cluster")
	createCloudDatabaseClusterUserCmd.Flags().String("user", "", "Username to create")
	createCloudDatabaseClusterUserCmd.Flags().String("password", "", "Password for the user")
	createCloudDatabaseClusterUserCmd.Flags().StringSlice("permission", []string{}, "Permissions in the form database:permission (repeatable)")
	createCloudDatabaseClusterUserCmd.MarkFlagRequired("namespace")
	createCloudDatabaseClusterUserCmd.MarkFlagRequired("cluster")
	createCloudDatabaseClusterUserCmd.MarkFlagRequired("user")
	createCloudDatabaseClusterUserCmd.MarkFlagRequired("password")
	cloudDatabaseClusterUserCmd.AddCommand(createCloudDatabaseClusterUserCmd)

	// modify
	modifyCloudDatabaseClusterUserCmd.Flags().StringP("namespace", "n", "", "Namespace name")
	modifyCloudDatabaseClusterUserCmd.Flags().String("cluster", "", "Name of the cluster")
	modifyCloudDatabaseClusterUserCmd.Flags().String("user", "", "Username to modify")
	modifyCloudDatabaseClusterUserCmd.Flags().String("password", "", "New password for the user (optional)")
	modifyCloudDatabaseClusterUserCmd.Flags().StringSlice("add-permission", []string{}, "Add permission in the form database:permission (repeatable)")
	modifyCloudDatabaseClusterUserCmd.Flags().StringSlice("remove-permission", []string{}, "Remove permission for a database. Accepts database or database:permission (repeatable)")
	modifyCloudDatabaseClusterUserCmd.MarkFlagRequired("namespace")
	modifyCloudDatabaseClusterUserCmd.MarkFlagRequired("cluster")
	modifyCloudDatabaseClusterUserCmd.MarkFlagRequired("user")
	cloudDatabaseClusterUserCmd.AddCommand(modifyCloudDatabaseClusterUserCmd)

	// delete
	deleteCloudDatabaseClusterUserCmd.Flags().StringP("namespace", "n", "", "Namespace name")
	deleteCloudDatabaseClusterUserCmd.Flags().String("cluster", "", "Name of the cluster")
	deleteCloudDatabaseClusterUserCmd.Flags().String("user", "", "Username to delete")
	deleteCloudDatabaseClusterUserCmd.MarkFlagRequired("namespace")
	deleteCloudDatabaseClusterUserCmd.MarkFlagRequired("cluster")
	deleteCloudDatabaseClusterUserCmd.MarkFlagRequired("user")
	cloudDatabaseClusterUserCmd.AddCommand(deleteCloudDatabaseClusterUserCmd)
}
