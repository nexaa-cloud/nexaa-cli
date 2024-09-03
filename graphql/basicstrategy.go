package graphql

import (
	"fmt"
	"reflect"
	"strings"
)

type BasicStrategy struct{}

func (s *BasicStrategy) CanHandle(field reflect.Value, fieldType reflect.Type) bool {
	return field.Kind() != reflect.Struct && field.Kind() != reflect.Slice
}

func (s *BasicStrategy) BuildQueryPart(field reflect.Value, fieldType reflect.Type, indentLevel int, qb *QueryBuilder) string {
	indent := strings.Repeat("    ", indentLevel)
	return fmt.Sprintf("%s", indent)
}
