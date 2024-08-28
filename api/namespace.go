package api

import (
	"context"

	"github.com/shurcooL/graphql"
	"gitlab.com/Tilaa/tilaa-cli/config"
)

type Namespace struct {
	Name string
	Id   string
}

var namespaceQuery struct {
	Namespaces []struct {
		Id   graphql.String
		Name graphql.String
	}
}

func ListNamespaces() ([]Namespace, error) {
	initGraphQLClientWithToken(config.AccessToken)

	err := client.Query(context.Background(), &namespaceQuery, nil)
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
