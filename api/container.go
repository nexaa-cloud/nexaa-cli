package api

import (
	"gitlab.com/Tilaa/tilaa-cli/config"
	"gitlab.com/Tilaa/tilaa-cli/graphql"
)

type Container struct {
	Name      string
	Image     string
	Namespace string
	State     string
}

type ContainerInput struct {
	Name      string
	Image     string
	Namespace int
	Http      string
	Https     string
	Env       []string
	SecretEnv []string
	Ports     []string
}

type Ingress struct {
	DomainName string `json:"domainName"`
	Port       int    `json:"port"`
	EnableTLS  bool   `json:"enableTLS"`
}

var containerQuery struct {
	Namespace struct {
		Id         string
		Name       string
		Containers []struct {
			Id    string
			Name  string
			Image string
			State string
		}
	} `graphql:"namespace(id: $id)"`
}

func ListContainers(namespace string) ([]Container, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	params := map[string]graphql.Parameter{
		"id": graphql.NewId(namespace),
	}

	query := client.BuildQuery(&containerQuery, params)

	err := client.Query(query)
	if err != nil {
		return nil, err
	}

	var containers []Container

	namespace = string(containerQuery.Namespace.Name)

	for _, container := range containerQuery.Namespace.Containers {
		containers = append(containers, Container{
			Namespace: namespace,
			Name:      string(container.Name),
			Image:     string(container.Image),
			State:     string(container.State),
		})
	}

	return containers, nil
}

func CreateContainer(input ContainerInput) (Container, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	createContainerInput := map[string]any{
		"namespaceId":             input.Namespace,
		"resourceSpecificationId": 91,
		"image":                   input.Image,
		"name":                    input.Name,
		"ports":                   input.Ports,
	}

	ingresses := getIngresses(input)
	if len(ingresses) > 0 {
		createContainerInput["ingresses"] = ingresses
	}

	params := map[string]graphql.Parameter{
		"containerInput": graphql.NewComplexParameter("CreateContainerInput", createContainerInput),
	}

	mutation := client.BuildMutation("createContainer", params)

	err := client.Mutate(mutation)

	return Container{}, err
}

func getIngresses(input ContainerInput) []Ingress {
	var ingresses []Ingress

	if input.Https != "" {
		ingresses = append(ingresses, Ingress{
			Port:       80,
			DomainName: input.Https,
			EnableTLS:  true,
		})
	}

	return ingresses
}
