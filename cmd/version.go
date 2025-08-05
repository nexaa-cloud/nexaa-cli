package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// These variables are set via ldflags during build
// They reference the variables in main.go
var (
	Version   = "dev"
	BuildDate = "unknown"
	CommitSHA = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Print version information for the Nexaa CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("nexaa version %s\n", Version)
		fmt.Printf("Built: %s\n", BuildDate)
		fmt.Printf("Commit: %s\n", CommitSHA)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}