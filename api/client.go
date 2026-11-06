package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Khan/genqlient/graphql"
	"github.com/nexaa-cloud/nexaa-cli/config"
)

type authedTransport struct {
	key     string
	scheme  string
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	value := t.key
	if t.key == "" {
		return t.wrapped.RoundTrip(req)
	}

	if t.scheme != "" {
		value = t.scheme + " " + t.key
	}
	req.Header.Set("Authorization", value)
	return t.wrapped.RoundTrip(req)
}

type Client struct {
	client *graphql.Client
}

func NewClient() *Client {
	httpClient := http.Client{
		Transport: &authedTransport{
			key:     config.AccessToken,
			scheme:  "Bearer",
			wrapped: http.DefaultTransport,
		},
	}

	client := graphql.NewClient(config.GRAPHQL_URL, &httpClient)

	return &Client{client: &client}
}

func ApiKeyClient() *Client {
	httpClient := http.Client{
		Transport: &authedTransport{
			key:     config.ApiKey,
			scheme:  "",
			wrapped: http.DefaultTransport,
		},
	}

	client := graphql.NewClient(config.GRAPHQL_URL, &httpClient)

	return &Client{client: &client}
}

func AKTokenClient() *Client {
	httpClient := http.Client{
		Transport: &authedTransport{
			key:     config.AccessToken,
			scheme:  "",
			wrapped: http.DefaultTransport,
		},
	}

	client := graphql.NewClient(config.GRAPHQL_URL, &httpClient)

	return &Client{client: &client}
}

func BuildClient() *Client {
	if config.AuthMethod == "apikey" {
		if err := RefreshAKToken(); err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			os.Exit(1)
		}
		return AKTokenClient()
	}
	return NewClient()
}
