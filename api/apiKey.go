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

func (client *Client) ApiKeysList() ([]ApiKeyResult, error) {
	resp, err := getApiKeys(context.Background(), *client.client)
	if err != nil {
		return []ApiKeyResult{}, err
	}

	result := resp.GetApiKeys()

	apiKeys := make([]ApiKeyResult, len(result))
	for i, apiKey := range result {
		apiKeys[i] = ApiKeyResult{
			Identifier:     apiKey.Identifier,
			Name:           apiKey.Name,
			Authorizations: apiKey.Authorizations,
			StartDate:      apiKey.StartDate,
			EndDate:        apiKey.EndDate,
			Status:         apiKey.Status,
		}
	}
	return apiKeys, nil
}
