package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/Tilaa/tilaa-cli/config"
)

var rootCmd = &cobra.Command{
	Use:   "tilaa",
	Short: "A CLI tool to manage cloud resources on the Tilaa Serverless Platform.",
}

// Execute initializes the CLI
func Execute() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if config.AccessToken == "" && (len(os.Args) < 2 || os.Args[1] != "login") {
		fmt.Println("No access token found, please login first.")
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(containerCmd)
	rootCmd.AddCommand(registryCmd)
	rootCmd.AddCommand(namespaceCmd)
	rootCmd.AddCommand(loginCmd)
}
