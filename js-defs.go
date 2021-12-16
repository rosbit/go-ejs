package ejs

import (
	"github.com/dop251/goja"
)

type JsVm struct {
	vm *goja.Runtime
	env map[string]interface{}
}

const globalJsFuncs = `
function _CALL(funcName) {
	var args = Array.prototype.slice.call(arguments)
	return _cb_(funcName, 'call', args.slice(1))
}
function _ENV(name) {
	return _cb_(name, 'env')
}

var console = {log:print,warn:print,error:print,info:print}
var js = {call:_CALL, env:_ENV}
`
