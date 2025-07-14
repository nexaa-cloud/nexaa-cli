package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/tilaa/tilaa-cli/config"
)

var env string

var rootCmd = &cobra.Command{
	Use:   "tilaa",
	Short: "A CLI tool to manage cloud resources on the Tilaa Serverless Platform.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := config.Initialize(env); err != nil {
			fmt.Printf("Error initializing environment: %v\n", err)
			os.Exit(1)
		}

		if err := config.LoadConfig(); err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		if config.AccessToken == "" && cmd.Name() != "login" {
			fmt.Println("ERROR: No access token found, please login first.")
			fmt.Println("Run 'tilaa login' to authenticate.")
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&env, "env", config.GetEnvironment(), "Environment (dev, prod)")

	rootCmd.AddCommand(containerCmd)
	rootCmd.AddCommand(containerjobCmd)
	rootCmd.AddCommand(registryCmd)
	rootCmd.AddCommand(namespaceCmd)
	rootCmd.AddCommand(loginCmd)
}
