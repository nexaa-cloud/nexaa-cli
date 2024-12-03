package api

import (
	// "github.com/shurcooL/graphql"

	"gitlab.com/Tilaa/tilaa-cli/config"
	"gitlab.com/Tilaa/tilaa-cli/graphql"
)

type Registry struct {
	Id     string
	Name   string
	Source string
}

type RegistryInput struct {
	Namespace int
	Id        string
	Name      string
	Source    string
	Username  string
	Password  string
	Verify    bool
}

func ListRegistries(namespace string) ([]Registry, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	var registryQuery struct {
		Namespace struct {
			Id                string
			Name              string
			PrivateRegistries []struct {
				Id   string
				Name string
			}
		} `graphql:"namespace(id: $id)"`
	}

	params := map[string]graphql.Parameter{
		"id": graphql.NewId(namespace),
	}

	query := client.BuildQuery(&registryQuery, params)
	err := client.Query(query)

	if err != nil {
		return nil, err
	}

	var registries []Registry

	for _, registry := range registryQuery.Namespace.PrivateRegistries {
		registries = append(registries, Registry{
			Id:   string(registry.Id),
			Name: string(registry.Name),
		})
	}

	return registries, nil
}

func CreateRegistry(input RegistryInput) (Registry, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	params := map[string]graphql.Parameter{
		"namespaceId": graphql.NewInt(input.Namespace),
		"name":        graphql.NewString(input.Name),
		"source":      graphql.NewString(input.Source),
		"username":    graphql.NewString(input.Username),
		"password":    graphql.NewString(input.Password),
		"verify":      graphql.NewBool(input.Verify),
	}

	mutation := client.BuildMutation("addPrivateRegistry", params)

	err := client.Mutate(mutation)

	return Registry{}, err
}
