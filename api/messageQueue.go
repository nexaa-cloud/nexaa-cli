package api

import (
	"context"
)

func (client *Client) MessageQueueList() ([]MessageQueueResult, error) {
	resp, err := messageQueuesGet(context.Background(), *client.client)
	if err != nil {
		return []MessageQueueResult{}, err
	}

	result := resp.GetMessageQueues()

	return result, nil
}

func (client *Client) MessageQueueGet(input MessageQueueResourceInput) (MessageQueueResult, error) {
	resp, err := messageQueueGet(context.Background(), *client.client, input)
	if err != nil {
		return MessageQueueResult{}, err
	}
	return resp.GetMessageQueue(), nil
}

func (client *Client) MessageQueueCreate(input MessageQueueCreateInput) (MessageQueueResult, error) {
	resp, err := messageQueueCreate(context.Background(), *client.client, input)
	if err != nil {
		return MessageQueueResult{}, err
	}
	return resp.GetMessageQueueCreate(), nil
}

func (client *Client) MessageQueueModify(input MessageQueueModifyInput) (MessageQueueResult, error) {
	resp, err := messageQueueModify(context.Background(), *client.client, input)
	if err != nil {
		return MessageQueueResult{}, err
	}
	return resp.GetMessageQueueModify(), nil
}

func (client *Client) MessageQueueDelete(input MessageQueueResourceInput) (bool, error) {
	resp, err := messageQueueDelete(context.Background(), *client.client, input)
	if err != nil {
		return false, err
	}
	return resp.GetMessageQueueDelete(), nil
}

func (client *Client) MessageQueuePlans() ([]MessageQueuePlanResult, error) {
	resp, err := messageQueuePlansGet(context.Background(), *client.client)
	if err != nil {
		return []MessageQueuePlanResult{}, err
	}

	result := resp.GetMessageQueuePlans()

	return result, nil
}

func (client *Client) MessageQueueVersions() ([]MessageQueueVersionResult, error) {
	resp, err := messageQueueVersionsGet(context.Background(), *client.client)
	if err != nil {
		return []MessageQueueVersionResult{}, err
	}

	result := resp.GetMessageQueueVersions()

	return result, nil
}

func (client *Client) MessageQueueAdminCredentials(input MessageQueueResourceInput, username string) (MessageQueueUserCredentialsResult, error) {
	resp, err := messageQueueUserCredentialsGet(context.Background(), *client.client, input, username)
	if err != nil {
		return MessageQueueUserCredentialsResult{}, err
	}
	return resp.GetMessageQueueUserCredentials(), nil
}
