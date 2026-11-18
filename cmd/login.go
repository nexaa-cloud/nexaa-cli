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

		envPassword := os.Getenv("NEXAA_PASSWORD")
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

			if string(bytePassword) == "" {
				log.Fatalf("Received empty password. Aborting login...")
			}

			password = string(bytePassword)
		}

		fmt.Println("WARNING: Logging in with username and password is not recommended except for managing API keys.")
		fmt.Println("WARNING: We recommend to log out when your done or manually delete the " + config.TOKEN_FILE)
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
		fmt.Println() // Print a newline after API key input for clarity
		if err != nil {
			log.Fatalf("Failed to read API key: %v", err)
		}

		if string(byteApiKey) == "" {
			log.Fatalf("Received empty API key. Aborting login...")
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
