package cmd

import (
	"fmt"
	"log"
	"syscall"

	"gitlab.com/Tilaa/tilaa-cli/api"
	"gitlab.com/Tilaa/tilaa-cli/config"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// loginCmd defines the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login using username and password",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		if username == "" {
			fmt.Println("Username is required.")
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

func init() {
	loginCmd.Flags().StringP("username", "u", "", "Username for authentication")
	loginCmd.Flags().StringP("password", "p", "", "Password for authentication (optional, will be prompted if not provided)")
}
