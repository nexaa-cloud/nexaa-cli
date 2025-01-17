package graphql

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildQuery(t *testing.T) {
	qb := NewQueryBuilder()

	var query struct {
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

	var sb strings.Builder

	sb.WriteString("query ($id:ID!) {\n")
	sb.WriteString("    namespace(id: $id) {\n")
	sb.WriteString("        id\n")
	sb.WriteString("        name\n")
	sb.WriteString("        containers {\n")
	sb.WriteString("            id\n")
	sb.WriteString("            name\n")
	sb.WriteString("            image\n")
	sb.WriteString("            state\n")
	sb.WriteString("        }\n")
	sb.WriteString("    }\n")
	sb.WriteString("}\n")

	expectedQuery := sb.String()

	params := map[string]Parameter{
		"id": NewId("1"),
	}

	result := qb.BuildQuery(&query, params)

	assert.Equal(t, expectedQuery, result)
}

func TestBuildMutation(t *testing.T) {
	qb := NewQueryBuilder()

	var query struct {
		Id   int
		Name string
	}

	params := map[string]Parameter{
		"name": NewString("Testing"),
	}

	var sb strings.Builder

	sb.WriteString("mutation ($name: String!) {\n")
	sb.WriteString("    createNamespace (name: $name) {\n")
	sb.WriteString("        id\n")
	sb.WriteString("        name\n")
	sb.WriteString("    }\n")
	sb.WriteString("}\n")

	actual := qb.BuildMutationWithQuery("createNamespace", params, &query)

	assert.Equal(t, sb.String(), actual)
}

func TestBuildMutationWithOptionalParameter(t *testing.T) {
	qb := NewQueryBuilder()

	var query struct {
		Id   int
		Name string
	}

	params := map[string]Parameter{
		"name": Optional(NewString("Testing")),
	}

	var sb strings.Builder

	sb.WriteString("mutation ($name: String) {\n")
	sb.WriteString("    createNamespace (name: $name) {\n")
	sb.WriteString("        id\n")
	sb.WriteString("        name\n")
	sb.WriteString("    }\n")
	sb.WriteString("}\n")

	actual := qb.BuildMutationWithQuery("createNamespace", params, &query)

	assert.Equal(t, sb.String(), actual)
}

func TestMutationWithComplexTypes(t *testing.T) {
	qb := NewQueryBuilder()

	var sb strings.Builder

	sb.WriteString("mutation ($CreateContainerInput: ContainerInput!) {\n")
	sb.WriteString("    createContainer (CreateContainerInput: $CreateContainerInput)\n")
	sb.WriteString("}\n")

	expected := sb.String()

	var query struct{}

	params := map[string]Parameter{
		"CreateContainerInput": NewComplexParameter("ContainerInput", map[string]any{
			"namespaceId": 15,
		}),
	}

	actual := qb.BuildMutationWithQuery("createContainer", params, &query)

	assert.Equal(t, expected, actual)
}
