package api

import (
	"github.com/Khan/genqlient/graphql"
	"github.com/nexaa-cloud/nexaa-cli/config"
	"net/http"
)

type authedTransport struct {
	key     string
	scheme  string
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	value := t.key
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
