package api

import (
	"context"
)

func (client *Client) ContainerJobCreate(input ContainerJobCreateInput) (ContainerJobResult, error) {
	containerJobCreateResponse, err := containerJobCreate(context.Background(), *client.client, input)
	if err != nil {
		return ContainerJobResult{}, err
	}

	return containerJobCreateResponse.GetContainerJobCreate(), nil
}

func (client *Client) ContainerJobModify(input ContainerJobModifyInput) (ContainerJobResult, error) {
	containerJobCreateResponse, err := containerJobModify(context.Background(), *client.client, input)
	if err != nil {
		return ContainerJobResult{}, err
	}

	return containerJobCreateResponse.GetContainerJobModify(), nil
}

func (client *Client) ContainerJobList(namespace string) ([]ContainerJobResult, error) {

	containerJobListResponse, err := containerJobList(context.Background(), *client.client, namespace)

	if err != nil {
		return []ContainerJobResult{}, err
	}

	namespaceResult := containerJobListResponse.GetNamespace()

	result := make([]ContainerJobResult, len(namespaceResult.ContainerJobs))
	for i, job := range namespaceResult.ContainerJobs {
		result[i] = ContainerJobResult{
			Name:    job.Name,
			Image:   job.Image,
			Enabled: job.Enabled,
			State:   job.State,
		}
	}
	return result, nil
}

func (client *Client) ContainerJobDelete(namespace string, containerJobName string) (bool, error) {
	containerJobDeleteResponse, err := containerJobDelete(context.Background(), *client.client, namespace, containerJobName)
	if err != nil {
		return false, err
	}

	return containerJobDeleteResponse.GetContainerJobDelete(), nil
}
