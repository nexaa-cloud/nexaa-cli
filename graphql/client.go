package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	Endpoint string
	Token    string
}

// NewClient creates a new GraphQL client with the provided endpoint and optional Bearer token.
func NewClient(endpoint, token string) *Client {
	return &Client{
		Endpoint: endpoint,
		Token:    token,
	}
}

// Query represents a GraphQL query.
type Query struct {
	Query      string
	Variables  map[string]any
	ReturnData interface{}
}

// Mutation represents a GraphQL mutation.
type Mutation struct {
	Query      string
	Variables  map[string]any
	ReturnData interface{}
}

// BuildMutation initializes a mutation request.
func (c *Client) BuildMutation(mutationName string, params map[string]Parameter) *Mutation {
	qb := NewQueryBuilder()
	query := qb.BuildMutation(mutationName, params)

	variables := map[string]any{}
	for key, value := range params {
		variables[key] = value.GraphqlValue
	}

	return &Mutation{
		Query:      query,
		Variables:  variables,
		ReturnData: nil,
	}
}

func (c *Client) BuildMutationWithQuery(mutationName string, params map[string]Parameter, queryStruct interface{}) *Mutation {
	qb := NewQueryBuilder()
	query := qb.BuildMutationWithQuery(mutationName, params, queryStruct)

	variables := map[string]any{}
	for key, value := range params {
		variables[key] = value.GraphqlValue
	}

	return &Mutation{
		Query:      query,
		Variables:  variables,
		ReturnData: queryStruct,
	}
}

// BuildQuery initializes a query request.
func (c *Client) BuildQuery(queryStruct interface{}, params map[string]Parameter) *Query {
	qb := NewQueryBuilder()
	query := qb.BuildQuery(queryStruct, params)

	variables := map[string]any{}

	for key, value := range params {
		variables[key] = value.GraphqlValue
	}

	return &Query{
		Query:      query,
		Variables:  variables,
		ReturnData: queryStruct,
	}
}

// Mutate executes a GraphQL mutation.
func (c *Client) Mutate(mutation *Mutation) error {
	return c.executeRequest(mutation.Query, mutation.Variables, mutation.ReturnData)
}

// Query executes a GraphQL query.
func (c *Client) Query(query *Query) error {
	return c.executeRequest(query.Query, query.Variables, query.ReturnData)
}

// executeRequest performs the HTTP request to the GraphQL server.
func (c *Client) executeRequest(query string, variables map[string]any, returnData interface{}) error {
	payload, err := c.preparePayload(query, variables)
	if err != nil {
		return err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	// Perform the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("unauthorized request. Please log in")
	}
	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GraphQL request failed: %s", string(body))
	}

	// Parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read HTTP response: %v", err)
	}

	return c.processResponse(body, returnData)
}

// preparePayload marshals the GraphQL query and variables into a JSON payload.
func (c *Client) preparePayload(query string, variables map[string]any) ([]byte, error) {
	payload := map[string]any{
		"query":     query,
		"variables": variables,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GraphQL request: %v", err)
	}

	return payloadBytes, nil
}

// processResponse handles both successful data responses and GraphQL errors.
func (c *Client) processResponse(body []byte, returnData interface{}) error {
	// Temporary struct to capture both "data" and "errors"
	var response struct {
		Data   json.RawMessage `json:"data"`
		Errors []struct {
			Message    string           `json:"message"`
			Locations  []map[string]any `json:"locations"`
			Extensions map[string]any   `json:"extensions,omitempty"`
		} `json:"errors,omitempty"`
	}

	// Unmarshal the body into the response struct
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshal GraphQL response: %v", err)
	}

	// Handle GraphQL errors if they exist
	if len(response.Errors) > 0 {
		// Convert the error details into a readable format
		var errMessages []string
		for _, e := range response.Errors {
			errMessages = append(errMessages, e.Message)
		}
		return fmt.Errorf("GraphQL errors: %s", errMessages)
	}

	if returnData != nil {
		// Unmarshal the "data" part into the provided returnData struct
		if err := json.Unmarshal(response.Data, returnData); err != nil {
			return fmt.Errorf("failed to unmarshal GraphQL data: %v", err)
		}
	}

	return nil
}
