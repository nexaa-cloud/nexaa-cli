package graphql

import (
	"fmt"
	"reflect"
	"strings"
)

type QueryBuilder struct {
	strategies []MarshalStrategy
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		strategies: []MarshalStrategy{
			&StructStrategy{},
			&SliceStrategy{},
			&BasicStrategy{},
		},
	}
}

func (qb *QueryBuilder) BuildQuery(queryStruct interface{}, params map[string]Parameter) string {
	val := reflect.ValueOf(queryStruct).Elem()
	typ := reflect.TypeOf(queryStruct).Elem()

	query := qb.buildQueryPart(val, typ, 1)

	if len(params) > 0 {
		var queryParts []string
		for name, param := range params {
			queryParts = append(queryParts, fmt.Sprintf("$%s:%s", name, param.GraphqlType))
		}

		return fmt.Sprintf("query (%s) {\n%s\n}\n", strings.Join(queryParts, ","), query)
	}

	return fmt.Sprintf("query {\n%s\n}\n", query)
}

func (qb *QueryBuilder) BuildMutation(mutationName string, variables map[string]Parameter) string {
	return qb.buildMutationPart(mutationName, variables, "")
}

func (qb *QueryBuilder) BuildMutationWithQuery(mutationName string, variables map[string]Parameter, queryStruct interface{}) string {
	val := reflect.ValueOf(queryStruct).Elem()
	typ := reflect.TypeOf(queryStruct).Elem()

	query := qb.buildQueryPart(val, typ, 2)

	return qb.buildMutationPart(mutationName, variables, query)
}

func (qb *QueryBuilder) buildMutationPart(mutationName string, variables map[string]Parameter, query string) string {
	var sb strings.Builder

	var outerParams []string
	var innerParams []string

	required := ""
	for name, param := range variables {
		required = ""
		if param.Required {
			required = "!"
		}
		outerParams = append(outerParams, fmt.Sprintf("$%s: %s%s", name, param.GraphqlType, required))
		innerParams = append(innerParams, fmt.Sprintf("%s: $%s", name, name))
	}

	sb.WriteString(fmt.Sprintf("mutation (%s) {\n", strings.Join(outerParams, ",")))
	sb.WriteString(fmt.Sprintf("    %s (%s)", mutationName, strings.Join(innerParams, ",")))

	if len(query) > 0 {
		sb.WriteString(" {\n")
		sb.WriteString(query)
		sb.WriteString("\n    }\n")
	}
	sb.WriteString("}\n")

	return sb.String()
}

func (qb *QueryBuilder) buildQueryPart(val reflect.Value, typ reflect.Type, indentLevel int) string {
	var queryParts []string
	indent := strings.Repeat("    ", indentLevel)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)
		fieldType := field.Type
		fieldName := strings.ToLower(field.Name)
		tag := field.Tag.Get("graphql")

		if tag != "" {
			fieldName = tag
		}

		for _, strategy := range qb.strategies {
			if strategy.CanHandle(fieldVal, fieldType) {
				subQuery := strings.Trim(strategy.BuildQueryPart(fieldVal, fieldType, indentLevel, qb), " ")

				queryPart := fmt.Sprintf("%s%s", indent, fieldName)

				if len(subQuery) > 0 {
					queryPart += " " + subQuery
				}

				queryParts = append(queryParts, queryPart)
				break
			}
		}
	}

	return strings.Join(queryParts, "\n")
}
