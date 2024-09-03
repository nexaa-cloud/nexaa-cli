package graphql

type Parameter struct {
	GraphqlType  string
	GraphqlValue any
}

func NewId(value any) Parameter {
	return Parameter{
		GraphqlType:  "ID!",
		GraphqlValue: value,
	}
}
