package api

import (
	"context"

	"github.com/shurcooL/graphql"
	"gitlab.com/Tilaa/tilaa-cli/config"
)

type Container struct {
	Name      string
	Image     string
	Namespace string
	State     string
}

var containerQuery struct {
	Namespace struct {
		Id         graphql.String
		Name       graphql.String
		Containers []struct {
			Id    graphql.String
			Name  graphql.String
			Image graphql.String
			State graphql.String
		}
	} `graphql:"namespace(id: $id)"`
}

func ListContainers(namespace string) ([]Container, error) {
	initGraphQLClientWithToken(config.AccessToken)

	params := map[string]any{
		"id": graphql.ID(namespace),
	}

	err := client.Query(context.Background(), &containerQuery, params)
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
