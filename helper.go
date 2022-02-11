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

func LoadFileFromCache(path string, envs, vars map[string]interface{}) (ctx *JsVm, err error) {
	lock.Lock()
	defer lock.Unlock()

	jsC, ok := jsCtxCache[path]

	if !ok {
		ctx = NewVM(envs)
		if err = ctx.LoadFile(path, vars); err != nil {
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
	if jsC.mt.Before(mt) {
		if err = jsC.jsvm.LoadFile(path, vars); err != nil {
			return
		}
		jsC.mt = mt
	}
	ctx = jsC.jsvm
	return
}
