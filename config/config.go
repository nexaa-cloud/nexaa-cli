package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var (
	AccessToken  string // OAuth Access Token
	ExpiresAt    int64  // Token expiration time as Unix timestamp
	RefreshToken string // OAuth Refresh Token
)

// Config represents the structure of the configuration stored on disk
type Config struct {
	AccessToken  string `json:"access_token"`
	ExpiresAt    int64  `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
}

// SaveConfig writes the current configuration to disk
func SaveConfig() error {
	configData := Config{
		AccessToken:  AccessToken,
		ExpiresAt:    ExpiresAt,
		RefreshToken: RefreshToken,
	}

	// Serialize the config struct to JSON
	data, err := json.MarshalIndent(configData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %v", err)
	}

	// Write the JSON data to the token file
	err = os.WriteFile(TOKEN_FILE, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write config to file: %v", err)
	}

	return nil
}

// LoadConfig reads the configuration from disk
func LoadConfig() error {
	// Check if the token file exists
	if _, err := os.Stat(TOKEN_FILE); os.IsNotExist(err) {
		// File does not exist, so keep the values empty
		return nil
	}

	// Read the content of the token file
	data, err := os.ReadFile(TOKEN_FILE)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	// Deserialize the JSON content into the config struct
	var configData Config
	err = json.Unmarshal(data, &configData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config JSON: %v", err)
	}

	// Update the global variables with the loaded configuration
	AccessToken = configData.AccessToken
	ExpiresAt = configData.ExpiresAt
	RefreshToken = configData.RefreshToken

	return nil
}

// IsTokenExpired checks if the current access token is expired
func IsTokenExpired() bool {
	// Compare current time with ExpiresAt
	return time.Now().Unix() >= ExpiresAt
}
