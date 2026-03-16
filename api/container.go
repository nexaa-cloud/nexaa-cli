package api

import (
	"context"
	"fmt"
)

func (client *Client) ListContainers(namespace string) ([]ContainerResult, error) {
	containerResponse, err := containerList(context.Background(), *client.client, namespace)
	if err != nil {
		return []ContainerResult{}, err
	}

	namespaceResult := containerResponse.GetNamespace()
	if len(namespaceResult.Containers) == 0 {
		return []ContainerResult{}, nil
	}

	result := make([]ContainerResult, 0, len(namespaceResult.Containers))
	for _, container := range namespaceResult.Containers {
		result = append(result, container.ContainerResult)
	}

	return result, nil
}

func (client *Client) ListContainerByName(namespace string, containerName string) (ContainerResult, error) {
	container, err := containerByName(context.Background(), *client.client, namespace, containerName)
	if err != nil {
		return ContainerResult{}, err
	}

	if container == nil {
		return ContainerResult{}, fmt.Errorf("container %q not found in namespace %q", containerName, namespace)
	}

	return container.Container, err
}

func (client *Client) ContainerCreate(input ContainerCreateInput) (ContainerResult, error) {
	containerCreateResponse, err := containerCreate(context.Background(), *client.client, input)
	if err != nil {
		return ContainerResult{}, err
	}

	return containerCreateResponse.GetContainerCreate(), nil
}

func (client *Client) ContainerModify(input ContainerModifyInput) (ContainerResult, error) {
	containerModifyResponse, err := containerModify(context.Background(), *client.client, input)
	if err != nil {
		return ContainerResult{}, err
	}

	return containerModifyResponse.GetContainerModify(), nil
}

func (client *Client) ContainerDelete(namespace string, containerName string) (bool, error) {
	containerDeleteResponse, err := containerDelete(context.Background(), *client.client, namespace, containerName)
	if err != nil {
		return false, err
	}

	return containerDeleteResponse.GetContainerDelete(), nil
}
