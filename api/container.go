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
	Id        int
	Name      string
	Image     string
	Namespace int
	Http      string
	HttpPort  int
	Https     string
	HttpsPort int
	Registry  int
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
	}

	if len(input.Ports) > 0 {
		createContainerInput["ports"] = input.Ports
	}

	if input.Registry != 0 {
		createContainerInput["privateRegistryId"] = input.Registry
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

func ModifyContainer(input ContainerInput) (Container, error) {
	client := graphql.NewClient(config.GRAPHQL_URL, config.AccessToken)

	modifyContainerInput := map[string]any{
		"containerId": input.Id,
	}

	if len(input.Image) > 0 {
		modifyContainerInput["image"] = input.Image
	}

	if len(input.Ports) > 0 {
		modifyContainerInput["ports"] = input.Ports
	}

	if input.Registry != 0 {
		modifyContainerInput["privateRegistryId"] = input.Registry
	}

	ingresses := getIngresses(input)
	if len(ingresses) > 0 {
		modifyContainerInput["ingresses"] = ingresses
	}

	params := map[string]graphql.Parameter{
		"containerInput": graphql.NewComplexParameter("ConfigureContainerInput", modifyContainerInput),
	}

	mutation := client.BuildMutation("modifyContainer", params)

	err := client.Mutate(mutation)

	return Container{}, err
}

func getIngresses(input ContainerInput) []Ingress {
	var ingresses []Ingress

	if input.Http != "" && input.HttpPort != 0 {
		ingresses = append(ingresses, Ingress{
			Port:       input.HttpPort,
			DomainName: input.Http,
			EnableTLS:  false,
		})
	}

	if input.Https != "" && input.HttpsPort != 0 {
		ingresses = append(ingresses, Ingress{
			Port:       input.HttpsPort,
			DomainName: input.Https,
			EnableTLS:  true,
		})
	}

	return ingresses
}
