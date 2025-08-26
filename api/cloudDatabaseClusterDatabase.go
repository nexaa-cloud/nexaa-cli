package api

import (
	"context"
)

func (client *Client) CloudDatabaseClusterDatabaseList(input CloudDatabaseClusterResourceInput) (getCloudDatabaseClusterDatabasesCloudDatabaseCluster, error) {
	resp, err := getCloudDatabaseClusterDatabases(context.Background(), *client.client, input)
	if err != nil {
		return getCloudDatabaseClusterDatabasesCloudDatabaseCluster{}, err
	}
	return resp.GetCloudDatabaseCluster(), nil
}
