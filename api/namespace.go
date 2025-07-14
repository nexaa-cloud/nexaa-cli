package api

import "context"

func (client *Client) NamespacesList() ([]NamespaceResult, error) {
	namespaceResponse, err := namespaceList(context.Background(), *client.client)
	if err != nil {
		return []NamespaceResult{}, err
	}

	namespaceResult := namespaceResponse.GetNamespaces()

	result := make([]NamespaceResult, len(namespaceResult))
	for i, namespace := range namespaceResult {
		result[i] = NamespaceResult{
			Name: namespace.Name,
			Description: namespace.Description,
		}
	}
	
	return  result, nil
}

func (client *Client) NamespaceListByName(name string) (NamespaceResult, error) {
	namespaceResponse, err := namespaceListByName(context.Background(), *client.client, name)
	if err != nil {
		return NamespaceResult{}, err
	}

	namespaceResult := namespaceResponse.GetNamespace()

	result := NamespaceResult{
		Name: namespaceResult.Name,
		Description: namespaceResult.Description,
	}

	return result, nil
}

func (client *Client) NamespaceCreate(input NamespaceCreateInput) (NamespaceResult, error) {
	namespaceCreateResponse, err := namespaceCreate(context.Background(), *client.client, input)
	if err != nil {
		return NamespaceResult{}, err
	}

	return namespaceCreateResponse.GetNamespaceCreate(), nil
}

func (client *Client) NamespaceDelete(name string) (bool, error) {
	namespaceDeleteResponse, err := namespaceDelete(context.Background(), *client.client, name)
	if err != nil {
		return false, err
	}

	return namespaceDeleteResponse.GetNamespaceDelete(), nil
}