package api

import (
	// "github.com/shurcooL/graphql"

	"gitlab.com/tilaa/tilaa-cli/config"
	"gitlab.com/tilaa/tilaa-cli/graphql"
)

type Registry struct {
	Id     		string
	Namespace 	string
	Name   		string
	Source 		string
	Username	string
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

func ListRegistryByName(namespace string, registryname string) (*Registry, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	var registryQuery struct {
		PrivateRegistry []struct {
			Id			string
			Name 		string
			Source 		string
			Username 	string
		} `graphql:"namespace(name: $name)"`
	}

	params := map[string]graphql.Parameter{
		"namespace": graphql.NewString(namespace),
	}

	query := client.BuildQuery(&registryQuery, params)
	err := client.Query(query)
	if err != nil {
		return nil, err
	}

	var registry Registry

	for _, item := range registryQuery.PrivateRegistry {
		if item.Name == registryname {
			registry.Id = item.Id
			registry.Name = item.Name
			registry.Source = item.Source
			registry.Username = item.Username
		}
	}

	return &registry, err
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
