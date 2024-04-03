package ejs

import (
	"github.com/dop251/goja"
	"sync"
	"fmt"
	"os"
	"time"
)

func NewVM(env map[string]interface{}) *JsVm {
	js := &JsVm{env: env, lock: &sync.Mutex{}}
	js.createJSContext(nil)
	return js
}

func NewContext() *JsVm {
	return NewVM(nil)
}

func (js *JsVm) BeginSafeCall() {
	js.lock.Lock()
}

func (js *JsVm) EndSafeCall() {
	js.lock.Unlock()
}

func (js *JsVm) LoadFile(path string, vars map[string]interface{}) (err error) {
	_, err = js.EvalFile(path, vars)
	return
}

func (js *JsVm) LoadScript(script string, vars map[string]interface{}) (err error) {
	_, err = js.Eval(script, vars)
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

func (js *JsVm) GetGlobal(name string) (res interface{}, err error) {
	js.lock.Lock()
	defer js.lock.Unlock()

	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); ok {
				return
			}
			err = fmt.Errorf("panic in GetGlobal(%s): %v", name, r)
		}
	}()

	v := js.vm.Get(name)
	res = v.Export()
	return
}

func (js *JsVm) EvalFile(path string, vars ...map[string]interface{}) (res interface{}, err error) {
	js.lock.Lock()
	defer js.lock.Unlock()

	b, e := os.ReadFile(path)
	if e != nil {
		err = e
		return
	}
	if len(vars) > 0 {
		js.AddVars(vars[0])
	}
	p, e := goja.Compile(path, string(b), true)
	if e != nil {
		err = e
		return
	}

	v, e := js.vm.RunProgram(p)
	if e != nil {
		err = e
		return
	}
	res = v.Export()
	return
}

func (js *JsVm) Eval(script string, vars ...map[string]interface{}) (res interface{}, err error) {
	js.lock.Lock()
	defer js.lock.Unlock()

	if len(vars) > 0 {
		js.AddVars(vars[0])
	}
	v, e := js.vm.RunString(script)
	if e != nil {
		err = e
		return
	}
	res = v.Export()
	return
}

func (js *JsVm) CallFunc(funcName string, args ...interface{}) (res interface{}, err error) {
	js.lock.Lock()
	defer js.lock.Unlock()

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

// @param funcVarPtr  in format `var funcVar func(....) ...; funcVarPtr = &funcVar`
func (js *JsVm) BindFunc(funcName string, funcVarPtr interface{}) (err error) {
	js.lock.Lock()
	defer js.lock.Unlock()

	return js.vm.ExportTo(js.vm.Get(funcName), funcVarPtr)
}

func (js *JsVm) BindFuncs(funcName2FuncVarPtr map[string]interface{}) (err error) {
	for funcName, funcVarPtr := range funcName2FuncVarPtr {
		if err = js.BindFunc(funcName, funcVarPtr); err != nil {
			return
		}
	}
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

func formatTimestamp(tm int64, layout ...string) string {
	var l string
	if len(layout) > 0 && len(layout[0]) > 0 {
		l = layout[0]
	} else {
		l = time.DateTime
	}
	return time.Unix(tm, 0).Format(l)
}

func (js *JsVm) createJSContext(vars map[string]interface{}) {
	js.vm = goja.New()
	reqReg.Enable(js.vm)
	js.vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
	js.AddVars(vars)
	js.AddVars(map[string]interface{}{
		"_cb_": js.callback,
		"print": fmt.Println,
		"formatTimestamp": formatTimestamp,
		"sprintf": fmt.Sprintf,
	})
	js.vm.RunString(globalJsFuncs)
}
