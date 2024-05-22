package common

import "net/http"

var upsertBeforeHooks = map[string]func(r *http.Request, in any) (out any, err error){}
var deleteHooks = map[string]func(r *http.Request, in any) (out any, err error){}
var batchDeleteHooks = map[string]func(r *http.Request, in any) (out any, err error){}

func GetUpsertBeforeHook(key string) (f func(r *http.Request, in any) (out any, err error), exists bool) {
	f, exists = upsertBeforeHooks[key]
	return
}

func RegisterUpsertBeforeHook(key string, f func(r *http.Request, in any) (out any, err error)) {
	upsertBeforeHooks[key] = f
}

func GetDeleteBeforeHook(key string) (f func(r *http.Request, in any) (out any, err error), exists bool) {
	f, exists = deleteHooks[key]
	return
}

func RegisterDeleteBeforeHook(key string, f func(r *http.Request, in any) (out any, err error)) {
	deleteHooks[key] = f
}

func GetBatchDeleteBeforeHook(key string) (f func(r *http.Request, in any) (out any, err error), exists bool) {
	f, exists = batchDeleteHooks[key]
	return
}

func RegisterBatchDeleteBeforeHook(key string, f func(r *http.Request, in any) (out any, err error)) {
	batchDeleteHooks[key] = f
}
