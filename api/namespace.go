package api

import (
	// "github.com/shurcooL/graphql"

	"gitlab.com/tilaa/tilaa-cli/config"
	"gitlab.com/tilaa/tilaa-cli/graphql"
)

type Namespace struct {
	Name 		string
	Id   		string
	Description string
}

func ListNamespaces() ([]Namespace, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	var namespaceQuery struct {
		Namespaces []struct {
			Id   		string
			Name 		string
			Description string
		}
	}

	params := map[string]graphql.Parameter{}

	query := client.BuildQuery(&namespaceQuery, params)
	err := client.Query(query)

	if err != nil {
		return nil, err
	}

	var namespaces []Namespace

	for _, namespace := range namespaceQuery.Namespaces {
		namespaces = append(namespaces, Namespace{
			Id:   string(namespace.Id),
			Name: string(namespace.Name),
		})
	}

	return namespaces, nil
}

func ListNamespaceByName(name string) (*Namespace, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	var namespaceQuery struct {
		Namespace struct {
			Id			string
			Name 		string
			Description	string
		} `graphql:"namespace(name: $name)"`
	}

	params := map[string]graphql.Parameter{
		"name": graphql.NewString(name),
	}

	query := client.BuildQuery(&namespaceQuery, params)
	err := client.Query(query)

	if err != nil {
		return nil, err
	}

	var namespace Namespace

	namespace.Id = namespaceQuery.Namespace.Id
	namespace.Name = namespaceQuery.Namespace.Name
	namespace.Description = namespaceQuery.Namespace.Description

	return &namespace, nil
}

func CreateNamespace(name string, description string) error {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	customerId, err := GetAccountId()
	if err != nil {
		return err
	}

	params := map[string]graphql.Parameter{
		"customerId":    graphql.NewInt(customerId),
		"pricingPlanId": graphql.NewInt(2),
		"name":          graphql.NewString(name),
		"description":   graphql.NewString(description),
	}

	mutation := client.BuildMutation("createNamespace", params)

	err = client.Mutate(mutation)
	if err != nil {
		return err
	}

	return nil
}


func DeleteNamespace(id int) error {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	params := map[string]graphql.Parameter{
		"id": graphql.NewInt(id),
	}

	mutation := client.BuildMutation("deleteNamespace", params)

	err := client.Mutate(mutation)
	if err != nil {
		return err
	}

	return nil
}
