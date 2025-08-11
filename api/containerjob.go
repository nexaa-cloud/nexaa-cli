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
		var registryName string
		if job.PrivateRegistry == nil {
			registryName = "public"
		} else {
			registryName = job.PrivateRegistry.Name
		}

		var envVars []ContainerJobResultEnvironmentVariablesEnvironmentVariable
		if job.EnvironmentVariables != nil {
			envVars = make([]ContainerJobResultEnvironmentVariablesEnvironmentVariable, len(job.EnvironmentVariables))
			for j, ev := range job.EnvironmentVariables {
				envVars[j] = ContainerJobResultEnvironmentVariablesEnvironmentVariable{
					Name:   ev.Name,
					Value:  ev.Value,
					Secret: ev.Secret,
				}
			}
		}

		var mounts []ContainerJobResultMountsMount
		if job.Mounts != nil {
			mounts = make([]ContainerJobResultMountsMount, len(job.Mounts))
			for j, m := range job.Mounts {
				var volumeName string
				if m.Volume.Name != "" {
					volumeName = m.Volume.Name
				}
				mounts[j] = ContainerJobResultMountsMount{
					Path: m.Path,
					Volume: ContainerJobResultMountsMountVolume{
						Name: volumeName,
					},
				}
			}
		}

		result[i] = ContainerJobResult{
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
