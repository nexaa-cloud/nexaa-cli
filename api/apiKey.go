package api

import "context"

func (client *Client) AccessTokenCreate() (string, error) {
	resp, err := accessTokenCreate(context.Background(), *client.client)
	if err != nil {
		return "", err
	}

	return resp.GetAccessTokenCreate(), nil
}

func (client *Client) ApiKeyCreate(input ApiKeyInput) (string, error) {
	resp, err := apiKeyCreate(context.Background(), *client.client, input)
	if err != nil {
		return "", err
	}
	return resp.GetApiKeyCreate(), nil
}
