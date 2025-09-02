package api

import (
	"context"
	"fmt"
)

func toContainerResult(container ContainerResult) (ContainerResult, error) {
	// Environment Variables
	var envVars []EnvironmentVariableResult
	if container.EnvironmentVariables != nil {
		envVars = make([]EnvironmentVariableResult, len(container.EnvironmentVariables))
		for j, ev := range container.EnvironmentVariables {
			envVars[j] = EnvironmentVariableResult{
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
				State:      v.State,
			}
		}
	}

	// Mounts
	var mounts []ContainerMounts
	if container.Mounts != nil {
		mounts = make([]ContainerMounts, len(container.Mounts))
		for j, m := range container.Mounts {
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
	registry := container.PrivateRegistry

	return ContainerResult{
		Name:                 container.Name,
		Image:                container.Image,
		PrivateRegistry:      registry,
		Resources:            container.Resources,
		EnvironmentVariables: envVars,
		Ports:                container.Ports,
		Ingresses:            ingresses,
		Mounts:               mounts,
		HealthCheck:          healthCheck,
		NumberOfReplicas:     container.NumberOfReplicas,
		AutoScaling:          autoScaling,
		Locked:               container.Locked,
		State:                container.State,
	}, nil
}

func (client *Client) ListContainers(namespace string) ([]ContainerResult, error) {
	containerResponse, err := containerList(context.Background(), *client.client, namespace)
	if err != nil {
		return []ContainerResult{}, err
	}

	namespaceResult := containerResponse.GetNamespace()
	result := make([]ContainerResult, len(namespaceResult.Containers))

	for _, container := range namespaceResult.Containers {
		var c, _ = toContainerResult(container.ContainerResult)
		result = append(result, c)
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

	res, err := toContainerResult(container.Container)

	return res, err
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
