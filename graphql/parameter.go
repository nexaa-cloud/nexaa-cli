package graphql

type Parameter struct {
	GraphqlType  string
	GraphqlValue any
	Required     bool
}

func NewId(value any) Parameter {
	return Parameter{
		GraphqlType:  "ID",
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

func NewBool(value any) Parameter {
	return Parameter{
		GraphqlType:  "Boolean",
		GraphqlValue: value,
		Required:     true,
	}
}

func NewComplexParameter(graphqlType string, value any) Parameter {
	return Parameter{
		GraphqlType:  graphqlType,
		GraphqlValue: value,
		Required:     true,
	}
}

func Optional(param Parameter) Parameter {
	param.Required = false
	return param
}
