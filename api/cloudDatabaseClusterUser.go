package api

import (
	"context"
	"errors"
)

func (client *Client) CloudDatabaseClusterUserList(input CloudDatabaseClusterResourceInput) ([]CloudDatabaseClusterUserResult, error) {
	resp, err := getCloudDatabaseClusterUsers(context.Background(), *client.client, input)
	if err != nil {
		return []CloudDatabaseClusterUserResult{}, err
	}

	var users []CloudDatabaseClusterUserResult
	cluster := resp.GetCloudDatabaseCluster()
	for _, user := range cluster.GetUsers() {
		users = append(users, CloudDatabaseClusterUserResult{
			Name:        user.Name,
			Status:      user.Status,
			Permissions: user.Permissions,
			Dsn:         user.Dsn,
			Password:    user.Password,
			Role:        user.Role,
		})
	}

	return users, nil
}

func (client *Client) CloudDatabaseClusterUserGet(input CloudDatabaseClusterResourceInput, name string) (CloudDatabaseClusterUserResult, error) {
	resp, err := client.CloudDatabaseClusterUserList(input)
	if err != nil {
		return CloudDatabaseClusterUserResult{}, err
	}

	var result CloudDatabaseClusterUserResult
	for _, user := range resp {
		if user.Name == name {
			result = user
		}
	}

	if result.Name == "" {
		return CloudDatabaseClusterUserResult{}, errors.New("user not found")
	}

	return result, nil
}

func (client *Client) CloudDatabaseClusterUserModify(input CloudDatabaseClusterUserModifyInput) (CloudDatabaseClusterUserResult, error) {
	resp, err := modifyCloudDatabaseClusterUser(context.Background(), *client.client, input)
	if err != nil {
		return CloudDatabaseClusterUserResult{}, err
	}

	user := CloudDatabaseClusterUserResult{
		Name:        resp.GetCloudDatabaseClusterUserModify().Name,
		Status:      resp.GetCloudDatabaseClusterUserModify().Status,
		Permissions: resp.GetCloudDatabaseClusterUserModify().Permissions,
		Dsn:         resp.GetCloudDatabaseClusterUserModify().Dsn,
		Password:    resp.GetCloudDatabaseClusterUserModify().Password,
		Role:        resp.GetCloudDatabaseClusterUserModify().Role,
	}

	return user, nil
}

func (client *Client) CloudDatabaseClusterUserCreate(input CloudDatabaseClusterUserCreateInput) (CloudDatabaseClusterUserResult, error) {
	resp, err := createCloudDatabaseClusterUser(context.Background(), *client.client, input)
	if err != nil {
		return CloudDatabaseClusterUserResult{}, err
	}

	user := CloudDatabaseClusterUserResult{
		Name:        resp.GetCloudDatabaseClusterUserCreate().Name,
		Status:      resp.GetCloudDatabaseClusterUserCreate().Status,
		Permissions: resp.GetCloudDatabaseClusterUserCreate().Permissions,
		Dsn:         resp.GetCloudDatabaseClusterUserCreate().Dsn,
		Password:    resp.GetCloudDatabaseClusterUserCreate().Password,
		Role:        resp.GetCloudDatabaseClusterUserCreate().Role,
	}

	return user, nil
}

func (client *Client) CloudDatabaseClusterUserDelete(input CloudDatabaseClusterUserResourceInput) (bool, error) {
	resp, err := deleteCloudDatabaseClusterUser(context.Background(), *client.client, input)
	if err != nil {
		return false, err
	}
	return resp.GetCloudDatabaseClusterUserDelete(), nil
}
