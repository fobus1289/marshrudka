package router

import (
	"fmt"
	request2 "github.com/fobus1289/marshrudka/router/request"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type router struct {
	Path        string
	Match       *regexp.Regexp
	WhereMatch  *regexp.Regexp
	HandlerFunc http.HandlerFunc
	Params      []string
	Methods     map[string]bool
	Handlers    handlers
}

type routers []*router

func (rs routers) Find(res http.ResponseWriter, req *http.Request) bool {

	for _, route := range rs {
		if route.Has(res, req) {
			return true
		}
	}

	return false
}

func (r *router) Has(res http.ResponseWriter, req *http.Request) bool {

	var (
		isMatch  = r.Match.MatchString(req.URL.Path)
		isMethod = r.Methods[req.Method]
	)

	if !isMethod && isMatch {
		res.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = res.Write(methodNotAllowed)
		return true
	}

	if !isMethod || !isMatch {
		return false
	}

	if r.WhereMatch != nil && !r.WhereMatch.MatchString(req.URL.Path) {
		return false
	}

	r.HandlerFunc(res, req)

	return true
}

func (r *router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var pm = reflectMap{}
	r.Handlers.Next(res, req, pm)
}

func (r *router) Where(key, pattern string) IMatch {
	var path = r.Path

	if key != "*" {

		if !request2.InArray(key, r.Params) {
			log.Fatalln(fmt.Sprintf("%s dont exist", key))
		}

		key = fmt.Sprintf(":%s", key)
	}

	path = strings.Replace(path, key, pattern, -1)

	r.WhereMatch = regexp.MustCompile(fmt.Sprintf("^(/?%s/?)$", path))

	return r
}

func (r *router) WhereIn(pattern map[string]string) IMatch {

	var path = r.Path

	for k, v := range pattern {
		if !strings.HasPrefix(v, "(") && !strings.HasSuffix(v, ")") {
			v = fmt.Sprintf("(%s)", v)
		}

		if k != "*" {
			k = fmt.Sprintf(":%s", k)
		}

		var value = fmt.Sprintf("%s", v)
		path = strings.Replace(path, k, value, 1)
	}

	r.WhereMatch = regexp.MustCompile(fmt.Sprintf("^(/?%s/?)$", path))

	return r
}

func (r *router) StripPrefix(prefix string) IMatch {
	prefix = strings.Replace(prefix, " ", "", -1)

	if !strings.HasPrefix(prefix, "/") {
		prefix = fmt.Sprintf("/%s", prefix)
	}

	if !strings.HasSuffix(prefix, "/") {
		prefix = fmt.Sprintf("%s/", prefix)
	}

	r.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		req.URL.Path = strings.TrimPrefix(req.URL.Path, prefix)
		req.URL.RawPath = strings.TrimPrefix(req.URL.RawPath, prefix)
		r.ServeHTTP(res, req)
	}

	return r
}
