package router

import (
	"net/http"
	"strings"
)

type IRouter interface {
}

type router struct {
	Path              string
	Paths             []string
	HttpUrlValidators []func(string) (bool, bool)
	Handlers          handlers
	Call              Call
	Services          reflectMapFunc
}

type routers []*router

func (rs routers) Find(w http.ResponseWriter, r *http.Request) (*router, bool, map[string]string) {

	var paths []string

	urlPath := r.URL.Path
	{
		if urlPath != "/" {
			urlPath = strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/"), "/")
			paths = strings.Split(urlPath, "/")
		}
	}

	for _, router := range rs {

		ok, params := router.HasUrl(urlPath, paths)

		if !ok {
			continue
		}

		return router, true, params
	}

	return nil, false, nil
}

func (r *router) HasUrl(urlPath string, paths []string) (bool, map[string]string) {

	if urlPath == r.Path {
		return true, nil
	}

	httpUrlValidators := r.HttpUrlValidators

	if len(paths) != len(httpUrlValidators) {
		return false, nil
	}

	var params = map[string]string{}

	for i := 0; i < len(paths); i++ {

		validator := httpUrlValidators[i]

		path := paths[i]

		valid, hasParam := validator(path)

		if !valid {
			return false, nil
		}

		if hasParam {
			params[r.Paths[i][1:]] = path
		}

	}

	return true, params
}

type Routers = routers
type Router = router

func (rs routers) Len() int {
	return len(rs)
}

func (rs routers) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs routers) Less(i, j int) bool {
	return strings.Count(rs[i].Path, ":") == 0
}

func (s *server) GET(actionPath string, handlers ...any) IRouter {
	return newRouter(handlers, s, []string{http.MethodGet}, actionPath)
}

func (s *server) POST(actionPath string, handlers ...any) IRouter {
	return newRouter(handlers, s, []string{http.MethodPost}, actionPath)
}

func (s *server) PUT(actionPath string, handlers ...any) IRouter {
	return newRouter(handlers, s, []string{http.MethodPut}, actionPath)
}

func (s *server) PATCH(actionPath string, handlers ...any) IRouter {
	return newRouter(handlers, s, []string{http.MethodPatch}, actionPath)
}

func (s *server) DELETE(actionPath string, handlers ...any) IRouter {
	return newRouter(handlers, s, []string{http.MethodDelete}, actionPath)
}

func (s *server) MATCH(actionPath string, methods []string, handlers ...any) IRouter {
	return newRouter(handlers, s, methods, actionPath)
}

func (s *server) ANY(actionPath string, handlers ...any) IRouter {
	return newRouter(handlers, s, []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	}, actionPath)
}
