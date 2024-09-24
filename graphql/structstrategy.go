package graphql

import (
	"fmt"
	"reflect"
	"strings"
)

type StructStrategy struct{}

func (s *StructStrategy) CanHandle(field reflect.Value, fieldType reflect.Type) bool {
	return field.Kind() == reflect.Struct
}

func (s *StructStrategy) BuildQueryPart(field reflect.Value, fieldType reflect.Type, indentLevel int, qb *QueryBuilder) string {
	subQuery := qb.buildQueryPart(field, fieldType, indentLevel+1)
	indent := strings.Repeat("    ", indentLevel)
	return fmt.Sprintf("{\n%s\n%s}", subQuery, indent)
}
