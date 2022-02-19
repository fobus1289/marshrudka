package v2

import (
	"net/http"
	"strings"
)

type server struct {
	services          paramsMap
	HandlersInterface handlersInterface
	routers           routers
}

func NewServer() IServer {
	return &server{
		services:          paramsMap{},
		HandlersInterface: handlersInterface{},
		routers:           routers{},
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !s.routers.Find(w, r) {
		http.NotFound(w, r)
	}
}

func (s *server) send() {

}

func (s *server) Use(handlers ...interface{}) {
	s.HandlersInterface.AddRange(handlers)
}

func (s *server) Group(path string, handlers ...interface{}) IRouter {
	var _handlersInterface = handlersInterface{}
	return &group{
		Path:              trimPrefixAndSuffix(path),
		HandlersInterface: *_handlersInterface.AddRange(s.HandlersInterface).AddRange(handlers),
		Server:            s,
	}
}

func (s *server) GET(path string, handlers ...interface{}) {
	_handlers := parseFunc(s, false, handlers)
	route := &router{
		Path:   path,
		Match:  createRequestRegular(getRegular(path)),
		Params: getPattern(path),
		Methods: map[string]bool{
			http.MethodGet: true,
		},
		Handlers: _handlers,
	}
	_handlers.SetRouter(route)
	s.routers = append(s.routers, route)
}

func (s *server) POST(path string, handlers ...interface{}) {
	_handlers := parseFunc(s, false, handlers)
	route := &router{
		Path:   path,
		Match:  createRequestRegular(getRegular(path)),
		Params: getPattern(path),
		Methods: map[string]bool{
			http.MethodPost: true,
		},
		Handlers: _handlers,
	}
	_handlers.SetRouter(route)
	s.routers = append(s.routers, route)
}

func (s *server) PUT(path string, handlers ...interface{}) {
	_handlers := parseFunc(s, false, handlers)
	route := &router{
		Path:   path,
		Match:  createRequestRegular(getRegular(path)),
		Params: getPattern(path),
		Methods: map[string]bool{
			http.MethodPut: true,
		},
		Handlers: _handlers,
	}
	_handlers.SetRouter(route)
	s.routers = append(s.routers, route)
}

func (s *server) PATCH(path string, handlers ...interface{}) {
	_handlers := parseFunc(s, false, handlers)
	route := &router{
		Path:   path,
		Match:  createRequestRegular(getRegular(path)),
		Params: getPattern(path),
		Methods: map[string]bool{
			http.MethodPatch: true,
		},
		Handlers: _handlers,
	}
	_handlers.SetRouter(route)
	s.routers = append(s.routers, route)
}

func (s *server) DELETE(path string, handlers ...interface{}) {
	_handlers := parseFunc(s, false, handlers)
	route := &router{
		Path:   path,
		Match:  createRequestRegular(getRegular(path)),
		Params: getPattern(path),
		Methods: map[string]bool{
			http.MethodDelete: true,
		},
		Handlers: _handlers,
	}
	_handlers.SetRouter(route)
	s.routers = append(s.routers, route)
}

func (s *server) ANY(path string, handlers ...interface{}) {
	_handlers := parseFunc(s, false, handlers)
	route := &router{
		Path:   path,
		Match:  createRequestRegular(getRegular(path)),
		Params: getPattern(path),
		Methods: map[string]bool{
			http.MethodGet:    true,
			http.MethodPost:   true,
			http.MethodPut:    true,
			http.MethodPatch:  true,
			http.MethodDelete: true,
		},
		Handlers: _handlers,
	}
	_handlers.SetRouter(route)
	s.routers = append(s.routers, route)
}

func (s *server) MATCH(path string, methods []string, handlers ...interface{}) {

	if len(methods) < 1 {
		panic("MATCH error")
	}

	var mapMethods = map[string]bool{}

	for _, method := range methods {
		mapMethods[strings.ToUpper(method)] = true
	}

	_handlers := parseFunc(s, false, handlers)

	route := &router{
		Path:     path,
		Match:    createRequestRegular(getRegular(path)),
		Params:   getPattern(path),
		Methods:  mapMethods,
		Handlers: _handlers,
	}
	_handlers.SetRouter(route)
	s.routers = append(s.routers, route)
}

func (s *server) Run(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *server) RunTLS(addr, certFile, keyFile string) error {
	panic("implement me")
}

func (s *server) RunAsync(addr string) error {
	panic("implement me")
}

func (s *server) RunAsyncTLS(addr, certFile, keyFile string) error {
	panic("implement me")
}
