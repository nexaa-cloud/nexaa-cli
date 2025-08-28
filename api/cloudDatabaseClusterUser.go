package api

import (
	"context"
)

func (client *Client) CloudDatabaseClusterUserList(input CloudDatabaseClusterResourceInput) (getCloudDatabaseClusterUsersCloudDatabaseCluster, error) {
	resp, err := getCloudDatabaseClusterUsers(context.Background(), *client.client, input)
	if err != nil {
		return getCloudDatabaseClusterUsersCloudDatabaseCluster{}, err
	}
	return resp.GetCloudDatabaseCluster(), nil
}
