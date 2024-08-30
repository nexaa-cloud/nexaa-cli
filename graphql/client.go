package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
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

// Mutation represents a GraphQL mutation.
type Mutation struct {
	Query      string
	Variables  map[string]any
	ReturnData any
}

// BuildMutation initializes a mutation request.
func (c *Client) BuildMutation(queryStruct interface{}, variables map[string]any) *Mutation {
	query := buildGraphQLQuery(queryStruct)
	return &Mutation{
		Query:      query,
		Variables:  variables,
		ReturnData: queryStruct,
	}
}

// Query represents a GraphQL query.
type Query struct {
	Query      string
	Variables  map[string]any
	ReturnData any
}

// BuildQuery initializes a query request.
func (c *Client) BuildQuery(queryStruct interface{}, variables map[string]any) *Query {
	query := buildGraphQLQuery(queryStruct)
	fmt.Println(query)
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
	// Build the request payload
	payload := map[string]any{
		"query":     query,
		"variables": variables,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal GraphQL request: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer(payloadBytes))
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

	// Unmarshal into the provided struct
	if err := json.Unmarshal(body, returnData); err != nil {
		return fmt.Errorf("failed to unmarshal GraphQL response: %v", err)
	}

	return nil
}

// buildGraphQLQuery dynamically constructs a GraphQL query or mutation string based on the provided struct.
func buildGraphQLQuery(queryStruct interface{}) string {
	query := buildGraphQLQueryPart(reflect.ValueOf(queryStruct).Elem(), reflect.TypeOf(queryStruct).Elem())
	return fmt.Sprintf("query { %s }", query)
}

// buildGraphQLQueryPart recursively builds the query string for each part of the struct.
func buildGraphQLQueryPart(val reflect.Value, typ reflect.Type) string {
	var queryParts []string

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name
		tag := field.Tag.Get("graphql")

		if tag != "" {
			// Use the tag as the field name if it's specified
			fieldName = tag
		}

		// Handle struct fields by recursively building their sub-query
		if val.Field(i).Kind() == reflect.Struct {
			subQuery := buildGraphQLQueryPart(val.Field(i), val.Field(i).Type())
			queryParts = append(queryParts, fmt.Sprintf("%s { %s }", fieldName, subQuery))
		} else if val.Field(i).Kind() == reflect.Slice {
			// Handle slice fields by generating the query for the first element type
			if val.Field(i).Len() > 0 {
				subQuery := buildGraphQLQueryPart(val.Field(i).Index(0), val.Field(i).Type().Elem())
				queryParts = append(queryParts, fmt.Sprintf("%s { %s }", fieldName, subQuery))
			} else {
				// Generate the query for an empty slice based on the element type
				elemType := reflect.New(val.Field(i).Type().Elem()).Elem()
				subQuery := buildGraphQLQueryPart(elemType, elemType.Type())
				queryParts = append(queryParts, fmt.Sprintf("%s { %s }", fieldName, subQuery))
			}
		} else {
			queryParts = append(queryParts, strings.ToLower(fieldName))
		}
	}

	return strings.Join(queryParts, " ")
}
