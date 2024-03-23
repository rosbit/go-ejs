package ejs

import (
	"sync"
	"os"
	"time"
)

type jsCtx struct {
	jsvm *JsVm
	mt   time.Time
}

var (
	jsCtxCache map[string]*jsCtx
	lock *sync.Mutex
)

func InitCache() {
	if lock != nil {
		return
	}
	lock = &sync.Mutex{}
	jsCtxCache = make(map[string]*jsCtx)
}

func LoadFileFromCache(path string, vars map[string]interface{}) (ctx *JsVm, existing bool, err error) {
	return LoadFileFromCacheWithEnvs(path, nil, vars)
}

func LoadFileFromCacheWithEnvs(path string, envs, vars map[string]interface{}) (ctx *JsVm, existing bool, err error) {
	lock.Lock()
	defer lock.Unlock()

	jsC, ok := jsCtxCache[path]

	if !ok {
		if ctx, err = createJSContext(path, envs, vars); err != nil {
			return
		}
		fi, _ := os.Stat(path)
		jsC = &jsCtx{
			jsvm: ctx,
			mt: fi.ModTime(),
		}
		jsCtxCache[path] = jsC
		return
	}

	fi, e := os.Stat(path)
	if e != nil {
		err = e
		return
	}
	mt := fi.ModTime()
	if !jsC.mt.Equal(mt) {
		if ctx, err = createJSContext(path, envs, vars); err != nil {
			return
		}
		jsC.jsvm = ctx
		jsC.mt = mt
	} else {
		existing = true
		ctx = jsC.jsvm
	}
	return
}

func createJSContext(path string, envs, vars map[string]interface{}) (ctx *JsVm, err error) {
	ctx = NewVM(envs)
	err = ctx.LoadFile(path, vars)
	return
}
