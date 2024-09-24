package api

import (
	// "github.com/shurcooL/graphql"
	"gitlab.com/Tilaa/tilaa-cli/config"
	"gitlab.com/Tilaa/tilaa-cli/graphql"
)

type Namespace struct {
	Name string
	Id   string
}

func ListNamespaces() ([]Namespace, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	var namespaceQuery struct {
		Namespaces []struct {
			Id   string
			Name string
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

func CreateNamespace(name string, description string) error {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	params := map[string]graphql.Parameter{
		"customerId":              graphql.NewInt(1),
		"pricingPlanId":           graphql.NewInt(1),
		"resourceSpecificationId": graphql.NewInt(1),
		"name":                    graphql.NewString(name),
		"description":             graphql.NewString(description),
	}

	mutation := client.BuildMutation("createNamespace", params)

	err := client.Mutate(mutation)
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
