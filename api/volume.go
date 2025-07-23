package api

import (
	"context"
	"fmt"
)

func (client *Client) ListVolumes(namespace string) ([]VolumeResult, error) {
	volumeResponse, err := volumeList(context.Background(), *client.client, namespace)
	if err != nil {
		return []VolumeResult{}, err
	}

	namespaceResult := volumeResponse.GetNamespace()

	result := make([]VolumeResult, len(namespaceResult.Volumes))
	for i, vol := range namespaceResult.Volumes {
		result[i] = VolumeResult{
			Name:   vol.Name,
			Size:   vol.Size,
			Usage:  vol.Usage,
			State:  vol.State,
			Locked: vol.Locked,
		}
	}

	return result, nil
}

func (client *Client) ListVolumeByName(namespace string, volumeName string) (*VolumeResult, error) {
	volumes, err := client.ListVolumes(namespace)
	if err != nil {
		return nil, err
	}

	for _, vol := range volumes {
		if vol.Name == volumeName {
			return &vol, nil
		}
	}

	return nil, fmt.Errorf("volume %q not found in namespace %q", volumeName, namespace)
}

func (client *Client) VolumeCreate(input VolumeCreateInput) (VolumeResult, error) {
	volumeCreateResponse, err := volumeCreate(context.Background(), *client.client, input)
	if err != nil {
		return VolumeResult{}, err
	}

	return volumeCreateResponse.GetVolumeCreate(), nil
}

func (client *Client) VolumeIncrease(input VolumeModifyInput) (VolumeResult, error) {
	volumeIncreaseResponse, err := volumeIncrease(context.Background(), *client.client, input)
	if err != nil {
		return VolumeResult{}, nil
	}

	return volumeIncreaseResponse.GetVolumeIncrease(), nil
}

func (client *Client) VolumeDelete(namespace string, volumeName string) (bool, error) {
	volumeDeleteResponse, err := volumeDelete(context.Background(), *client.client, namespace, volumeName)
	if err != nil {
		return false, err
	}

	return volumeDeleteResponse.GetVolumeDelete(), nil
}
