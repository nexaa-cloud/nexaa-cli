package api

import (
	"gitlab.com/tilaa/tilaa-cli/config"
	"gitlab.com/tilaa/tilaa-cli/graphql"
)

type Registry struct {
	Namespace 	string
	Name   		string
	Source 		string
	Username	string
	Locked 		bool
}

type RegistryInput struct {
	Namespace string
	Name      string
	Source    string
	Username  string
	Password  string
	Verify    bool
}

type RegistryResponse struct {
	Name 		string				`json:"name"`
	Namespace 	NamespaceResponse	`json:"namespace"`
	Source 		string				`json:"source"`
	Username	string				`json:"username"`
	Locked 		bool				`json:"locked"`
}

func ListRegistries(namespace string) ([]Registry, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	var registryQuery struct {
		Namespace struct {
			Name              string
			PrivateRegistries []struct {
				Name 		string
				Source		string
				Username	string
				Locked		bool
			}
		} `graphql:"namespace(name: $name)"`
	}

	params := map[string]graphql.Parameter{
		"name": graphql.NewString(namespace),
	}

	query := client.BuildQuery(&registryQuery, params)
	err := client.Query(query)

	if err != nil {
		return nil, err
	}

	var registries []Registry

	for _, registry := range registryQuery.Namespace.PrivateRegistries {
		registries = append(registries, Registry{
			Name: string(registry.Name),
			Source: string(registry.Source),
			Username: string(registry.Username),
			Locked: bool(registry.Locked),
		})
	}

	return registries, nil
}

func ListRegistryByName(namespace string, registry string) (*Registry, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	var registryQuery struct {
		Namespace struct {
			Name 	string
			PrivateRegistries []struct {
				Id			string
				Name 		string
				Source 		string
				Username 	string
				Locked		bool
			} 
		} `graphql:"namespace(name: $name)"`
	}

	params := map[string]graphql.Parameter{
		"name": graphql.NewString(namespace),
	}

	query := client.BuildQuery(&registryQuery, params)
	err := client.Query(query)

	if err != nil {
		return nil, err
	}

	var reg Registry

	for _, item := range registryQuery.Namespace.PrivateRegistries {
		if item.Name == registry {
			reg.Name = item.Name
			reg.Source = item.Source
			reg.Username = item.Username
		}
	}

	reg.Namespace = registryQuery.Namespace.Name

	return &reg, err
}

func CreateRegistry(input RegistryInput) (Registry, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	createRegistryInput := map[string]any{
		"namespace": input.Namespace,
		"name": input.Name,
		"source": input.Source,
		"username": input.Username,
		"password": input.Password,
		"verify": input.Verify,
	}

	params := map[string]graphql.Parameter{
		"registryInput": graphql.NewComplexParameter("RegistryCreateInput", createRegistryInput),
	}

	var resp RegistryResponse

	mutation := client.BuildMutationWithQuery("registryConnectionCreate", params, &resp)

	err := client.Mutate(mutation)

	var registry Registry
	registry.Name = resp.Name
	registry.Namespace = resp.Namespace.Name
	registry.Source = resp.Source
	registry.Username = resp.Username
	registry.Locked = resp.Locked	

	return registry, err
}

func DeleteRegistry(namespace string, name string) error {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	registryInput := map[string]any{
		"namespace": namespace,
		"name": name,
	}

	params := map[string]graphql.Parameter{
		"registryConnection": graphql.NewComplexParameter("DeleteRegistryConnectionInput", registryInput),
	}

	mutation := client.BuildMutation("registryConnectionDelete", params)

	err := client.Mutate(mutation)
	if err != nil {
		return err
	}

	return nil
}
