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

		return fmt.Sprintf("query (%s) {\n%s\n}", strings.Join(queryParts, ","), query)
	}

	return fmt.Sprintf("query {\n%s\n}", query)
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
				subQuery := strategy.BuildQueryPart(fieldVal, fieldType, indentLevel, qb)
				queryParts = append(queryParts, fmt.Sprintf("%s%s %s", indent, fieldName, subQuery))
				break
			}
		}
	}

	return strings.Join(queryParts, "\n")
}
