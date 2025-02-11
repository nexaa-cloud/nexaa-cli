package api

import (
	"github.com/Khan/genqlient/graphql"
	"gitlab.com/Tilaa/tilaa-cli/config"
	"net/http"
)

type authedTransport struct {
	key     string
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.key)
	return t.wrapped.RoundTrip(req)
}

type Client struct {
	client *graphql.Client
}

func NewClient() *Client {
	httpClient := http.Client{
		Transport: &authedTransport{
			key:     config.AccessToken,
			wrapped: http.DefaultTransport,
		},
	}

	client := graphql.NewClient(config.GRAPHQL_URL, &httpClient)

	return &Client{client: &client}
}
