package graphql

import "reflect"

type MarshalStrategy interface {
	CanHandle(field reflect.Value, fieldType reflect.Type) bool
	BuildQueryPart(field reflect.Value, fieldType reflect.Type, indentLevel int, qb *QueryBuilder) string
}
