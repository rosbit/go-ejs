package main

import (
	"github.com/rosbit/go-ejs"
	"os"
	"fmt"
)

func goAdd(a, b int) int {
	return a+b
}

type A struct {
	Name string `json:"name"`
	Age int `json:"age"`
}

var (
	vars = map[string]interface{}{
		"add": goAdd,
		"a":&A{Name:"rosbit", Age:10},
		"m": map[string]interface{}{
			"p": map[string]interface{}{
				"n": "name",
				"a": 10,
			},
			"a": []int{1, 2, 3},
		},
	}
	envs = map[string]interface{}{"TZ": "Asia/Shanghai", "age": 10}
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s <js-file>[ <func>]\n", os.Args[0])
		return
	}

	jsFile := os.Args[1]
	var jsFunc string
	if len(os.Args) >= 3 {
		jsFunc = os.Args[2]
	}

	jsVM := ejs.NewVM(envs)
	jsVM.AddVars(vars)
	res, err := jsVM.EvalFile(jsFile)

	if len(jsFunc) > 0 {
		res, err = jsVM.CallFunc(jsFunc, 1, 3)
	}
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf(" => %v\n", res)
}
