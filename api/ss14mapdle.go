package main

import (
	"fmt"
	"log/slog"
	"os"
	"ss14mapdle/functions"
)

func printFunctionError() {
	slog.Error("[main]\tAvailable functions:")
	for key, value := range functions.Functions {
		slog.Error(fmt.Sprintf("[main]\t\t%s: %s", key, value.Description))
	}
}

func main() {
	if len(os.Args) == 1 {
		slog.Error("[main] Missing required argument")
		slog.Error("[main]\tUsage: ss14mapdle function")
		printFunctionError()
		return
	}

	os.Args = os.Args[1:]

	functionName := functions.FunctionName(os.Args[0])

	function, ok := functions.Functions[functionName]

	if !ok {
		slog.Error(fmt.Sprintf("[main] unknown function: %s", functionName))
		printFunctionError()
		return
	}

	slog.Info(fmt.Sprintf("[main] Starting main function: %v", functionName))

	var args []string

	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	function.F(args)
}
