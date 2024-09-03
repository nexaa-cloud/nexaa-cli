package graphql

import (
	"fmt"
	"reflect"
	"strings"
)

type SliceStrategy struct{}

func (s *SliceStrategy) CanHandle(field reflect.Value, fieldType reflect.Type) bool {
	return field.Kind() == reflect.Slice
}

func (s *SliceStrategy) BuildQueryPart(field reflect.Value, fieldType reflect.Type, indentLevel int, qb *QueryBuilder) string {
	var subQuery string
	if field.Len() > 0 {
		subQuery = qb.buildQueryPart(field.Index(0), fieldType.Elem(), indentLevel+1)
	} else {
		elemType := reflect.New(fieldType.Elem()).Elem()
		subQuery = qb.buildQueryPart(elemType, elemType.Type(), indentLevel+1)
	}
	indent := strings.Repeat("    ", indentLevel)
	return fmt.Sprintf("{\n%s\n%s}", subQuery, indent)
}
