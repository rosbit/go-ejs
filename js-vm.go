package ejs

import (
	"github.com/dop251/goja"
	"fmt"
	"os"
)

func NewVM(env map[string]interface{}) *JsVm {
	js := &JsVm{env: env}
	js.createJSContext(nil)
	return js
}

func (js *JsVm) LoadFile(path string, vars map[string]interface{}) (err error) {
	b, e := os.ReadFile(path)
	if e != nil {
		err = e
		return
	}
	return js.LoadScript(string(b), vars)
}

func (js *JsVm) LoadScript(script string, vars map[string]interface{}) (err error) {
	js.createJSContext(vars)
	_, err = js.vm.RunString(script)
	return
}

func (js *JsVm) AddVars(vars map[string]interface{}) {
	for n, v := range vars {
		js.vm.Set(n, v)
	}
}

func (js *JsVm) AddVar(name string, val interface{}) {
	js.vm.Set(name, val)
}

func (js *JsVm) EvalFile(path string) (res interface{}, err error) {
	b, e := os.ReadFile(path)
	if e != nil {
		err = e
		return
	}
	return js.Eval(string(b))
}

func (js *JsVm) Eval(script string) (res interface{}, err error) {
	v, e := js.vm.RunString(script)
	if e != nil {
		err = e
		return
	}
	res = v.Export()
	return
}

func (js *JsVm) CallFunc(funcName string, args ...interface{}) (res interface{}, err error) {
	f, ok := goja.AssertFunction(js.vm.Get(funcName))
	if !ok {
		err = fmt.Errorf("function name %s is not found in JS script", funcName)
		return
	}
	valArgs := make([]goja.Value, len(args))
	for i, arg := range args {
		valArgs[i] = js.vm.ToValue(arg)
	}
	v, e := f(goja.Undefined(), valArgs...)
	if e != nil {
		err = e
		return
	}
	res = v.Export()
	return
}

func (js *JsVm) getEnv(key string) interface{} {
	if len(key) == 0 || len(js.env) == 0 {
		return goja.Undefined()
	}
	if v, ok := js.env[key]; ok {
		return v
	}
	return goja.Undefined()
}

func (js *JsVm) callGoFunc(funcName string, args ...interface{}) interface{} {
	switch funcName {
	case "dump":
		return js.dumpContext(args)
	default:
		fmt.Printf("funcName: %s, args: %#v\n", funcName, args)
		return 0
	}
}

func (js *JsVm) dumpContext(args ...interface{}) interface{} {
	fmt.Fprintf(os.Stderr, "env: %#v\n", js.env)
	return goja.Undefined()
}

func (js *JsVm) callback(name string, op string, args []interface{}) interface{} {
	switch op {
	case "call":
		return js.callGoFunc(name, args...)
	case "env":
		return js.getEnv(name)
	}
	return goja.Undefined()
}

func (js *JsVm) createJSContext(vars map[string]interface{}) {
	js.vm = goja.New()
	js.vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	js.AddVars(vars)
	js.vm.Set("_cb_", js.callback)
	js.vm.Set("print", fmt.Println)
	js.vm.RunString(globalJsFuncs)
}
