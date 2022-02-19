package v2

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
)

type router struct {
	Path     string
	Match    *regexp.Regexp
	Params   []string
	decoder  func(r io.Reader) json.Decoder
	encoder  func()
	Methods  map[string]bool
	Handlers handlers
}

type routers []*router

func (rs routers) Find(res http.ResponseWriter, req *http.Request) bool {
	var pm = paramsMap{}

	for _, route := range rs {
		if route.Has(res, req, pm) {
			return true
		}
	}

	return false
}

func (r *router) Has(res http.ResponseWriter, req *http.Request, pm paramsMap) bool {
	if !r.Methods[req.Method] {
		return false
	}

	if !r.Match.MatchString(req.URL.Path) {
		return false
	}
	r.Handlers.Next(res, req, pm)
	return true
}
