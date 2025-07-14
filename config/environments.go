package config

import (
	"fmt"
	"os"
)

type EnvConfig struct {
	GraphQLURL       string
	KeycloakURL      string
	KeycloakClientID string
	KeycloakRealm    string
	TokenFile        string
}

var environments = map[string]EnvConfig{
	"dev": {
		GraphQLURL:       "https://staging-graphql.tilaa.com/graphql/platform",
		KeycloakURL:      "https://staging-auth.tilaa.com",
		KeycloakClientID: "cloud-tilaa",
		KeycloakRealm:    "tilaa",
		TokenFile:        "./auth.json",
	},
	"prod": {
		GraphQLURL:       "https://graphql.tilaa.com/graphql/platform",
		KeycloakURL:      "https://auth.tilaa.com",
		KeycloakClientID: "cloud-tilaa",
		KeycloakRealm:    "tilaa",
		TokenFile:        "./auth.json",
	},
}

var (
	GRAPHQL_URL        string
	KEYCLOAK_URL       string
	KEYCLOAK_CLIENT_ID string
	KEYCLOAK_REALM     string
	TOKEN_FILE         string
)

// Initialize sets up the environment configuration
func Initialize(env string) error {
	if env == "" {
		env = "prod" // default
	}

	config, exists := environments[env]
	if !exists {
		return fmt.Errorf("unknown environment: %s", env)
	}

	GRAPHQL_URL = config.GraphQLURL
	KEYCLOAK_URL = config.KeycloakURL
	KEYCLOAK_CLIENT_ID = config.KeycloakClientID
	KEYCLOAK_REALM = config.KeycloakRealm
	TOKEN_FILE = config.TokenFile

	return nil
}

// GetEnvironment returns current environment from env var or default
func GetEnvironment() string {
	env := os.Getenv("TILAA_ENV")
	if env == "" {
		return "prod"
	}
	return env
}
