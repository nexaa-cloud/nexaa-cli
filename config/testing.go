//go:build dev

package config

const (
	GRAPHQL_URL        = "https://staging-graphql.tilaa.com"
	KEYCLOAK_URL       = "https://staging-auth.tilaa.com"
	KEYCLOAK_CLIENT_ID = "cloud-tilaa"
	KEYCLOAK_REALM     = "tilaa"
	TOKEN_FILE         = "./auth.json"
)
