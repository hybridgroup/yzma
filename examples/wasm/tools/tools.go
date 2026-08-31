//go:build js && wasm

package main

import (
	"fmt"
	"strconv"

	"github.com/hybridgroup/yzma/pkg/message"
)

// toolDefinitions gives the tools that the model can call. They make their
// answer from the arguments alone, thus a test always gets the same result.
func toolDefinitions() []message.ToolDefinition {
	return []message.ToolDefinition{
		{
			Type: "function",
			Function: message.ToolFunctionDefinition{
				Name:        "get_weather",
				Description: "Get the current weather of a place",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"location": map[string]interface{}{
							"type":        "string",
							"description": "The name of the city",
						},
					},
					"required": []string{"location"},
				},
			},
		},
		{
			Type: "function",
			Function: message.ToolFunctionDefinition{
				Name:        "calculate",
				Description: "Do arithmetic on two numbers",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"a": map[string]interface{}{
							"type":        "number",
							"description": "The first number",
						},
						"b": map[string]interface{}{
							"type":        "number",
							"description": "The second number",
						},
						"op": map[string]interface{}{
							"type":        "string",
							"description": "One of add, subtract, multiply, divide",
						},
					},
					"required": []string{"a", "b", "op"},
				},
			},
		},
	}
}

// executeToolCall runs one tool call and gives its result as text.
func executeToolCall(call message.ToolCall) (string, error) {
	switch call.Function.Name {
	case "get_weather":
		location, ok := call.Function.Arguments["location"]
		if !ok || location == "" {
			return "", fmt.Errorf("get_weather needs a location")
		}
		return fmt.Sprintf("The weather in %s is 22 degrees Celsius and sunny.", location), nil

	case "calculate":
		a, err := argumentAsFloat(call.Function.Arguments, "a")
		if err != nil {
			return "", err
		}
		b, err := argumentAsFloat(call.Function.Arguments, "b")
		if err != nil {
			return "", err
		}
		return calculate(a, b, call.Function.Arguments["op"])

	default:
		return "", fmt.Errorf("unknown function %s", call.Function.Name)
	}
}

// calculate does one operation on two numbers.
func calculate(a, b float64, op string) (string, error) {
	var result float64

	switch op {
	case "add", "+":
		result = a + b
	case "subtract", "-":
		result = a - b
	case "multiply", "*":
		result = a * b
	case "divide", "/":
		if b == 0 {
			return "", fmt.Errorf("cannot divide by zero")
		}
		result = a / b
	default:
		return "", fmt.Errorf("unknown operation %s", op)
	}

	return strconv.FormatFloat(result, 'f', -1, 64), nil
}

// argumentAsFloat reads one argument as a number. Every argument of a tool call
// is text, thus this must parse it.
func argumentAsFloat(args map[string]string, key string) (float64, error) {
	value, ok := args[key]
	if !ok {
		return 0, fmt.Errorf("no argument %s", key)
	}

	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("argument %s is not a number: %s", key, value)
	}

	return f, nil
}
