package functions

type FunctionName string

type Function struct {
	F           func(params []string)
	Description string
}

var Functions map[FunctionName]Function

func init() {
	Functions = make(map[FunctionName]Function)
}

func AddFunction(functionName FunctionName, f func(params []string), description string) {
	Functions[functionName] = Function{
		F:           f,
		Description: description,
	}
}
