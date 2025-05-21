package api

import (
	// "github.com/shurcooL/graphql"

	"gitlab.com/tilaa/tilaa-cli/config"
	"gitlab.com/tilaa/tilaa-cli/graphql"
)

type Namespace struct {
	Name 		string
	Description string
}

type NamespaceInput struct {
	Name 		string
	Description string
}

type NamespaceResponse struct {
	Name 			string		`json:"name"`
	Description		string		`json:"description"`
}

func ListNamespaces() ([]Namespace, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	var namespaceQuery struct {
		Namespaces []struct {
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
			Name: string(namespace.Name),
			Description: string(namespace.Description),
		})
	}

	return namespaces, nil
}

func ListNamespaceByName(name string) (*Namespace, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	var namespaceQuery struct {
		Namespace struct {
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

	namespace.Name = namespaceQuery.Namespace.Name
	namespace.Description = namespaceQuery.Namespace.Description

	return &namespace, nil
}

func CreateNamespace(input NamespaceInput) (Namespace, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	//customerId, err := GetAccountId()

	createNamespaceInput := map[string]any{
		"name":				input.Name,
		"description":		input.Description,
	}

	params := map[string]graphql.Parameter{
		"namespaceInput": graphql.NewComplexParameter("NamespaceCreateInput", createNamespaceInput),
	}

	var resp NamespaceResponse

	mutation := client.BuildMutationWithQuery("namespaceCreate", params, &resp)

	err := client.Mutate(mutation)

	var namespace Namespace
	namespace.Name = resp.Name
	namespace.Description = resp.Description

	return namespace, err
}


func DeleteNamespace(name string) error {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	deleteNamespaceInput := map[string]any{
		"name" : name,
	}

	params := map[string]graphql.Parameter{
		"namespace": graphql.NewComplexParameter("DeleteNamespaceInput", deleteNamespaceInput),
	}

	mutation := client.BuildMutation("namespaceDelete", params)

	err := client.Mutate(mutation)
	if err != nil {
		return err
	}

	return nil
}
