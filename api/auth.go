package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gitlab.com/Tilaa/tilaa-cli/config"

	"github.com/shurcooL/graphql"
)

type OAuthResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
}

func Login(username, password string) error {
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("grant_type", "password")
	data.Set("client_id", "cloud-tilaa")

	fmt.Println(config.KEYCLOAK_URL)

	req, err := http.NewRequest("POST", config.KEYCLOAK_URL+"/realms/tilaa/protocol/openid-connect/token", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("invalid credentials")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	var oauthResp OAuthResponse
	if err := json.Unmarshal(body, &oauthResp); err != nil {
		return fmt.Errorf("failed to parse OAuth response: %v", err)
	}

	// Save the token and refresh token in config
	config.AccessToken = oauthResp.AccessToken
	config.RefreshToken = oauthResp.RefreshToken
	config.ExpiresAt = time.Now().Add(time.Duration(oauthResp.ExpiresIn) * time.Second).UnixMicro()

	config.SaveConfig()

	return nil
}

func initGraphQLClientWithToken(token string) {
	httpClient := &http.Client{
		Transport: &oauthTransport{token: token},
	}
	client = graphql.NewClient(config.GRAPHQL_URL, httpClient)
}

type oauthTransport struct {
	token string
}

func (t *oauthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	return http.DefaultTransport.RoundTrip(req)
}
