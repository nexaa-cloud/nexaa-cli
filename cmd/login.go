package cmd

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/nexaa-cloud/nexaa-cli/api"
	"github.com/nexaa-cloud/nexaa-cli/config"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// loginCmd defines the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login using username and password or API key",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")

		envPassword := os.Getenv("TILAA_PASSWORD")
		password, _ := cmd.Flags().GetString("password")

		if password == "" {
			password = envPassword
		}

		if username == "" && password == "" {
			loginApiKey(cmd)
			return
		}

		// If the password is not provided via the command-line argument, prompt for it via stdin
		if password == "" {
			fmt.Print("Enter password: ")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			fmt.Println() // Print a newline after password input for clarity
			if err != nil {
				log.Fatalf("Failed to read password: %v", err)
			}
			password = string(bytePassword)
		}

		err := api.Login(username, password)
		if err != nil {
			log.Fatalf("Login failed: %v", err)
		} else {
			fmt.Println("Login successful, access token stored in " + config.TOKEN_FILE)
		}
	},
}

func loginApiKey(cmd *cobra.Command) {
	fmt.Println("Retrieving API Key")
	envApiKey := os.Getenv("NEXAA_APIKEY")
	apiKey, _ := cmd.Flags().GetString("apikey")
	if apiKey == "" {
		apiKey = envApiKey
	}

	if apiKey == "" {
		fmt.Print("API key not found. Please enter API key: ")
		byteApiKey, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println() // Print a newline after password input for clarity
		if err != nil {
			log.Fatalf("Failed to read password: %v", err)
		}
		apiKey = string(byteApiKey)
	}

	fmt.Println("Logging in with API key")

	err := api.LoginApiKey(apiKey)
	if err != nil {
		log.Fatalf("Login failed with API key failed: %v", err)
	} else {
		fmt.Println("Login successful, access token stored in " + config.TOKEN_FILE)
	}
}

func init() {
	loginCmd.Flags().StringP("username", "u", "", "Username for authentication (optional, without username an API will be prompted)")
	loginCmd.Flags().StringP("password", "p", "", "Password for authentication (optional, will be prompted if not provided)")
	loginCmd.Flags().StringP("apikey", "a", "", "API key for authentication on Nexaa")
}
