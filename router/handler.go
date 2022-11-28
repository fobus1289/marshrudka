package router

import (
	"net/http"
	"reflect"
)

type handlerParam struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	SessionData    reflectMap
	Services       reflectMapFunc
	HttpParamMap   map[string]string
}

type handlers []*handler

type handler struct {
	Action   Call
	Params   paramFuncs
	Call     func([]reflect.Value) []reflect.Value
	OutTypes []reflect.Type
}
