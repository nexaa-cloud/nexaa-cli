package api

import (
	"context"
	"fmt"
)

func toContainerJobResult(job ContainerJobResult) (ContainerJobResult, error) {
	var registryName string
	if job.PrivateRegistry == nil {
		registryName = "public"
	} else {
		registryName = job.PrivateRegistry.Name
	}

	var envVars []EnvironmentVariableResult
	if job.EnvironmentVariables != nil {
		envVars = make([]EnvironmentVariableResult, len(job.EnvironmentVariables))
		for j, ev := range job.EnvironmentVariables {
			envVars[j] = EnvironmentVariableResult{
				Name:   ev.Name,
				Value:  ev.Value,
				Secret: ev.Secret,
			}
		}
	}

	var mounts []ContainerMounts
	if job.Mounts != nil {
		mounts = make([]ContainerMounts, len(job.Mounts))
		for j, m := range job.Mounts {
			var volumeName string
			if m.Volume.Name != "" {
				volumeName = m.Volume.Name
			}
			mounts[j] = ContainerMounts{
				Path: m.Path,
				Volume: ContainerMountsVolume{
					Name: volumeName,
				},
			}
		}
	}

	return ContainerJobResult{
		Name:                 job.Name,
		Image:                job.Image,
		PrivateRegistry:      &ContainerJobResultPrivateRegistry{Name: registryName},
		Resources:            job.Resources,
		EnvironmentVariables: envVars,
		Command:              job.Command,
		Mounts:               mounts,
		Schedule:             job.Schedule,
		Enabled:              job.Enabled,
		State:                job.State,
	}, nil
}

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
	for _, job := range namespaceResult.ContainerJobs {
		result = append(result, job.ContainerJobResult)
	}
	return result, nil
}

func (client *Client) ContainerJobByName(namespace string, name string) (ContainerJobResult, error) {
	apiResponse, err := containerJobByName(context.Background(), *client.client, namespace, name)

	if err != nil {
		return ContainerJobResult{}, err
	}

	if apiResponse == nil {
		return ContainerJobResult{}, fmt.Errorf("container job %q not found in namespace %q", name, namespace)
	}

	containerJob, err := toContainerJobResult(apiResponse.ContainerJob)

	return containerJob, err
}

func (client *Client) ContainerJobDelete(namespace string, containerJobName string) (bool, error) {
	containerJobDeleteResponse, err := containerJobDelete(context.Background(), *client.client, namespace, containerJobName)
	if err != nil {
		return false, err
	}

	return containerJobDeleteResponse.GetContainerJobDelete(), nil
}
