package graphql

type Parameter struct {
	GraphqlType  string
	GraphqlValue any
	Required     bool
}

func NewId(value any) Parameter {
	return Parameter{
		GraphqlType:  "ID!",
		GraphqlValue: value,
		Required:     true,
	}
}

func NewInt(value any) Parameter {
	return Parameter{
		GraphqlType:  "Int",
		GraphqlValue: value,
		Required:     true,
	}
}

func NewString(value any) Parameter {
	return Parameter{
		GraphqlType:  "String",
		GraphqlValue: value,
		Required:     true,
	}
}

func Optional(param Parameter) Parameter {
	param.Required = false
	return param
}
