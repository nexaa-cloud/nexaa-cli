package cmd

import (
	"log"

	"github.com/nexaa-cloud/nexaa-cli/api"

	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout of the nexaa CLI",
	Run: func(cmd *cobra.Command, args []string) {
		err := api.Logout()
		if err != nil {
			log.Fatalf("Logout failed: %v", err)
		}
	},
}
