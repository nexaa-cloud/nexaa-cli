package api

import (
	"context"
)

func (client *Client) CloudDatabaseClusterCreate(input CloudDatabaseClusterCreateInput) (CloudDatabaseClusterResult, error) {
	resp, err := cloudDatabaseClusterCreate(context.Background(), *client.client, input)
	if err != nil {
		return CloudDatabaseClusterResult{}, err
	}
	return resp.GetCloudDatabaseClusterCreate(), nil
}

func (client *Client) CloudDatabaseClusterModify(input CloudDatabaseClusterModifyInput) (CloudDatabaseClusterResult, error) {
	resp, err := cloudDatabaseClusterModify(context.Background(), *client.client, input)
	if err != nil {
		return CloudDatabaseClusterResult{}, err
	}
	return resp.GetCloudDatabaseClusterModify(), nil
}

func (client *Client) CloudDatabaseClusterList() ([]CloudDatabaseClusterResult, error) {
	resp, err := getCloudDatabaseClusters(context.Background(), *client.client)
	if err != nil {
		return []CloudDatabaseClusterResult{}, err
	}

	result := resp.GetCloudDatabaseClusters()

	return result, nil
}

func (client *Client) CloudDatabaseClusterDelete(input CloudDatabaseClusterResourceInput) (bool, error) {
	resp, err := cloudDatabaseClusterDelete(context.Background(), *client.client, input)
	if err != nil {
		return false, err
	}
	return resp.GetCloudDatabaseClusterDelete(), nil
}

func (client *Client) CloudDatabaseClusterDatabaseCreate(input CloudDatabaseClusterDatabaseCreateInput) (CloudDatabaseClusterDatabaseResult, error) {
	resp, err := createCloudDatabaseClusterDatabase(context.Background(), *client.client, input)
	if err != nil {
		return CloudDatabaseClusterDatabaseResult{}, err
	}
	return resp.GetCloudDatabaseClusterDatabaseCreate(), nil
}

func (client *Client) CloudDatabaseClusterDatabaseDelete(input CloudDatabaseClusterDatabaseResourceInput) (bool, error) {
	resp, err := deleteCloudDatabaseClusterDatabase(context.Background(), *client.client, input)
	if err != nil {
		return false, err
	}
	return resp.GetCloudDatabaseClusterDatabaseDelete(), nil
}

func (client *Client) CloudDatabaseClusterUserCredentials(cloudDatabase CloudDatabaseClusterResourceInput, userName string) (string, error) {
	resp, err := getCloudDatabaseClusterUserCredentials(context.Background(), *client.client, cloudDatabase, userName)
	if err != nil {
		return "", err
	}
	return resp.GetCloudDatabaseClusterUserCredentials().Dsn, nil
}

func (client *Client) CloudDatabaseClusterListPlans() ([]CloudDatabaseClusterPlan, error) {
	resp, err := clusterPlans(context.Background(), *client.client)
	if err != nil {
		return []CloudDatabaseClusterPlan{}, err
	}

	result := resp.GetCloudDatabaseClusterPlans()

	return result, nil
}

func (client *Client) CloudDatabaseClusterListSpecs() ([]CloudDatabaseClusterSpec, error) {
	resp, err := clusterVersions(context.Background(), *client.client)
	if err != nil {
		return []CloudDatabaseClusterSpec{}, err
	}

	result := resp.GetCloudDatabaseClusterVersions()

	return result, nil
}
