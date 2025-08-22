package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	GRAPHQL_URL        string
	KEYCLOAK_URL       string
	KEYCLOAK_CLIENT_ID string
	KEYCLOAK_REALM     string
	TOKEN_FILE         string
)

// Initialize sets up the environment configuration using individual environment variables
// Defaults to production values, can be overridden with environment variables
func Initialize() {
	// Try to load .env file if it exists, ignore errors if it doesn't exist
	_ = godotenv.Load()

	// Set production defaults
	GRAPHQL_URL = getEnvWithDefault("NEXAA_GRAPHQL_URL", "https://graphql.tilaa.com/graphql/platform")
	KEYCLOAK_URL = getEnvWithDefault("NEXAA_KEYCLOAK_URL", "https://auth.tilaa.com")
	KEYCLOAK_CLIENT_ID = "cloud-tilaa"
	KEYCLOAK_REALM = "tilaa"
	TOKEN_FILE = getEnvWithDefault("NEXAA_TOKEN_FILE", "./auth.json")
}

// getEnvWithDefault returns the environment variable value or the default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
