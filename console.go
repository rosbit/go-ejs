package ejs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/dop251/goja"
)

// ANSI Color in Terminal
const (
	cReset  = "\033[0m"
	cRed    = "\033[31m"
	cGreen  = "\033[32m"
	cYellow = "\033[33m"
	cBlue   = "\033[34m"
	cCyan   = "\033[36m"
	cGray   = "\033[90m"
)

// format JS Value as a string
func formatValue(val goja.Value) string {
	if val == nil || val == goja.Undefined() {
		return cGray + "undefined" + cReset
	}
	if val == goja.Null() {
		return cGray + "null" + cReset
	}

	switch val.ExportType().Kind() {
	case reflect.String:
		return cRed + val.String() + cReset
	case reflect.Bool:
		return cRed + fmt.Sprintf("%v", val.Export()) + cReset
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return cYellow + fmt.Sprintf("%v", val.Export()) + cReset
	default:
		// make Object, Array as a JSON string
		exported := val.Export()
		b, err := json.Marshal(exported)
		if err != nil {
			return cCyan + fmt.Sprintf("[Complex Object: %v]", exported) + cReset
		}
		return cCyan + string(b) + cReset
	}
}

func joinArgs(args []goja.Value) string {
	if len(args) == 0 {
		return ""
	}
	parts := make([]string, len(args))
	for i, arg := range args {
		parts[i] = formatValue(arg)
	}
	return strings.Join(parts, " ")
}

func consolePrint(call goja.FunctionCall) goja.Value {
	fmt.Println(joinArgs(call.Arguments))
	return goja.Undefined()
}

/*
func consolePrint(call goja.FunctionCall) goja.Value {
	args := make([]interface{}, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = arg // arg.Export()
	}
	fmt.Println(args...)
	return goja.Undefined()
}*/
