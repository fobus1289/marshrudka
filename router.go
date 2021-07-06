package marshrudka

import (
	"net/http"
	"regexp"
)

type router struct {
	path     string
	uri      *regexp.Regexp
	method   string
	params   []string
	actions  actions
	notFound func(responseWriter http.ResponseWriter, request *http.Request)
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
