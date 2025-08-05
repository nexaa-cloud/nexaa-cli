package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/nexaa-cloud/nexaa-cli/config"
	"github.com/spf13/cobra"
)

var env string

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate shell completion script",
	Long: `To enable shell completion, run:

For Bash:
    source <(nexaa completion bash)

For Zsh:
    source <(nexaa completion zsh)

Or to persist it, save the output to a file and source it in your shell config.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if len(args) > 0 && args[0] == "zsh" {
			err = cmd.Root().GenZshCompletion(os.Stdout)
		} else {
			err = cmd.Root().GenBashCompletion(os.Stdout)
		}
		if err != nil {
			fmt.Printf("Error generating completion: %v\n", err)
			os.Exit(1)
		}
	},
}

var rootCmd = &cobra.Command{
	Use:   "nexaa",
	Short: "A CLI tool to manage cloud resources on the Nexaa Serverless Platform.",
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
			fmt.Println("Run 'nexaa login' to authenticate.")
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
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(clouddatabaseclusterCmd)
	rootCmd.AddCommand(containerCmd)
	rootCmd.AddCommand(containerjobCmd)
	rootCmd.AddCommand(registryCmd)
	rootCmd.AddCommand(namespaceCmd)
	rootCmd.AddCommand(loginCmd)
}
