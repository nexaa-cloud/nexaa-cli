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

func (client *Client) CloudDatabaseClusterUserModify(input CloudDatabaseClusterUserModifyInput) (modifyCloudDatabaseClusterUserCloudDatabaseClusterUserModifyDatabaseUser, error) {
	resp, err := modifyCloudDatabaseClusterUser(context.Background(), *client.client, input)
	if err != nil {
		return modifyCloudDatabaseClusterUserCloudDatabaseClusterUserModifyDatabaseUser{}, err
	}
	return resp.GetCloudDatabaseClusterUserModify(), nil
}

func (client *Client) CloudDatabaseClusterUserCreate(input CloudDatabaseClusterUserCreateInput) (createCloudDatabaseClusterUserCloudDatabaseClusterUserCreateDatabaseUser, error) {
	resp, err := createCloudDatabaseClusterUser(context.Background(), *client.client, input)
	if err != nil {
		return createCloudDatabaseClusterUserCloudDatabaseClusterUserCreateDatabaseUser{}, err
	}
	return resp.GetCloudDatabaseClusterUserCreate(), nil
}

func (client *Client) CloudDatabaseClusterUserDelete(input CloudDatabaseClusterUserResourceInput) (bool, error) {
	resp, err := deleteCloudDatabaseClusterUser(context.Background(), *client.client, input)
	if err != nil {
		return false, err
	}
	return resp.GetCloudDatabaseClusterUserDelete(), nil
}
