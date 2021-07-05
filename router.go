package marshrudka

import (
	"net/http"
	"regexp"
)

type router struct {
	path       string
	uri        *regexp.Regexp
	method     string
	middleware []func(responseWriter http.ResponseWriter, request *http.Request) bool
	actions    actions
	notFound   func(responseWriter http.ResponseWriter, request *http.Request)
}

func (r *router) Middleware(middlewares func(responseWriter http.ResponseWriter, request *http.Request) bool) *router {
	r.middleware = append(r.middleware, middlewares)
	return r
}

type routers []*router

func (r *routers) Add(router *router) {
	*r = append(*r, router)
}

type group struct {
	Path    string
	actions []interface{}
	*drive
}
