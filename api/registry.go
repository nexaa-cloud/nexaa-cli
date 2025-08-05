package api

import (
	"context"
	"fmt"
)

func (client *Client) ListRegistries(namespace string) ([]RegistryResult, error) {
	registryResponse, err := registryList(context.Background(), *client.client, namespace)
	if err != nil {
		return []RegistryResult{}, err
	}

	namespaceResult := registryResponse.GetNamespace()

	result := make([]RegistryResult, len(namespaceResult.PrivateRegistries))
	for i, registry := range namespaceResult.PrivateRegistries {
		result[i] = RegistryResult{
			Name:     registry.Name,
			Source:   registry.Source,
			Username: registry.Username,
			State:    registry.State,
			Locked:   registry.Locked,
		}
	}

	return result, nil
}

func (client *Client) ListRegistryByName(namespace string, registryName string) (*RegistryResult, error) {
	registries, err := client.ListRegistries(namespace)
	if err != nil {
		return nil, err
	}

	for _, registry := range registries {
		if registry.Name == registryName {
			return &registry, nil
		}
	}

	return nil, fmt.Errorf("registry %q not found in namespace %q", registryName, namespace)
}

func (client *Client) RegistryCreate(input RegistryCreateInput) (RegistryResult, error) {
	registryCreateResponse, err := registryCreate(context.Background(), *client.client, input)
	if err != nil {
		return RegistryResult{}, err
	}

	return registryCreateResponse.GetRegistryConnectionCreate(), nil
}

func (client *Client) RegistryDelete(namespace string, registryName string) (bool, error) {
	registryDeleteResponse, err := registryDelete(context.Background(), *client.client, namespace, registryName)
	if err != nil {
		return false, err
	}

	return registryDeleteResponse.GetRegistryConnectionDelete(), nil
}
