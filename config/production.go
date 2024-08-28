//go:build !dev

package config

const (
	GRAPHQL_URL        = "https://graphql.tilaa.com"
	KEYCLOAK_URL       = "https://auth.tilaa.com"
	KEYCLOAK_CLIENT_ID = "cloud-tilaa"
	KEYCLOAK_REALM     = "tilaa"
	TOKEN_FILE         = "~/.tilaa/auth.json"
)
