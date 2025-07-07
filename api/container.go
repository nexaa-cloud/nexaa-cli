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
	result := make([]ContainerResult, len(namespaceResult.Containers))

	for i, container := range namespaceResult.Containers {
		var registryName string
		if container.PrivateRegistry == nil {
			registryName = "public"
		} else {
			registryName = container.PrivateRegistry.Name
		}

		// Environment Variables
		var envVars []ContainerResultEnvironmentVariablesEnvironmentVariable
		if container.EnvironmentVariables != nil {
			envVars = make([]ContainerResultEnvironmentVariablesEnvironmentVariable, len(container.EnvironmentVariables))
			for j, ev := range container.EnvironmentVariables {
				envVars[j] = ContainerResultEnvironmentVariablesEnvironmentVariable{
					Name:   ev.Name,
					Value:  ev.Value,
					Secret: ev.Secret,
				}
			}
		}

		// Ingresses
		var ingresses []ContainerResultIngressesIngress
		if container.Ingresses != nil {
			ingresses = make([]ContainerResultIngressesIngress, len(container.Ingresses))
			for j, v := range container.Ingresses {
				ingresses[j] = ContainerResultIngressesIngress{
					DomainName: v.DomainName,
					Port:       v.Port,
					EnableTLS:  v.EnableTLS,
					Allowlist:  v.Allowlist,
				}
			}
		}

		// Mounts
		var mounts []ContainerResultMountsMount
		if container.Mounts != nil {
			mounts = make([]ContainerResultMountsMount, len(container.Mounts))
			for j, m := range container.Mounts {
				var volumeName string
				if m.Volume.Name != "" {
					volumeName = m.Volume.Name
				}
				mounts[j] = ContainerResultMountsMount{
					Path: m.Path,
					Volume: ContainerResultMountsMountVolume{
						Name: volumeName,
					},
				}
			}
		}

		// Health Check
		var healthCheck *ContainerResultHealthCheck
		if container.HealthCheck != nil {
			healthCheck = &ContainerResultHealthCheck{
				Port: container.HealthCheck.Port,
				Path: container.HealthCheck.Path,
			}
		}

		// Auto Scaling
		var autoScaling *ContainerResultAutoScaling
		if container.AutoScaling != nil {
			var triggers []ContainerResultAutoScalingTriggersAutoScalingTrigger
			if container.AutoScaling.Triggers != nil {
				triggers = make([]ContainerResultAutoScalingTriggersAutoScalingTrigger, len(container.AutoScaling.Triggers))
				for j, t := range container.AutoScaling.Triggers {
					triggers[j] = ContainerResultAutoScalingTriggersAutoScalingTrigger{
						Type:      t.Type,
						Threshold: t.Threshold,
					}
				}
			}
			autoScaling = &ContainerResultAutoScaling{
				Replicas: ContainerResultAutoScalingReplicas{
					Minimum: container.AutoScaling.Replicas.Minimum,
					Maximum: container.AutoScaling.Replicas.Maximum,
				},
				Triggers: triggers,
			}
		}

		result[i] = ContainerResult{
			Name:            		container.Name,
			Image:           		container.Image,
			PrivateRegistry: 		&ContainerResultPrivateRegistry{Name: registryName},
			Resources:       	  	container.Resources,
			EnvironmentVariables: 	envVars,
			Ports:                	container.Ports,
			Ingresses:            	ingresses,
			Mounts:               	mounts,
			HealthCheck:          	healthCheck,
			NumberOfReplicas:     	container.NumberOfReplicas,
			AutoScaling:          	autoScaling,
			Locked:               	container.Locked,
		}
	}

	return result, nil
}



func (client *Client) ListContainerByName(namespace string, containerName string) (*ContainerResult, error) {
	containers, err := client.ListContainers(namespace)
	if err != nil {
		return nil, err
	}

	for _, c := range containers {
		if c.Name == containerName {
			return &c, nil
		}
	}

	return nil, fmt.Errorf("container %q not found in namespace %q", containerName, namespace)
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

