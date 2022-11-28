package router

import (
	"net/http"
)

type IMethod interface {
	GET(actionPath string, handlers ...any) IRouter
	POST(actionPath string, handlers ...any) IRouter
	PUT(actionPath string, handlers ...any) IRouter
	PATCH(actionPath string, handlers ...any) IRouter
	DELETE(actionPath string, handlers ...any) IRouter
	MATCH(actionPath string, methods []string, handlers ...any) IRouter
	ANY(actionPath string, handlers ...any) IRouter
}

type IServer interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Use(handlers ...any) IServer
	Group(actionPath string, handlers ...any) IGroup
	AddScoped(scoped any) IServer
	AddSingleton(singleton any) IServer
	DeserializeError(cb func(error) *RuntimeError)
	UseService()
	IMethod
}

func NewServer() IServer {

	deserializeErrorFunc := &RuntimeError{
		Status:      http.StatusBadRequest,
		ContentType: "text/plain; charset=utf-8",
	}

	runtimeErrorFunc := &RuntimeError{
		Status:      http.StatusInternalServerError,
		ContentType: "text/plain; charset=utf-8",
	}

	return &server{
		Services: reflectMapFunc{},
		Routers:  map[string]routers{},
		Call: func(_ *handlerParam) bool {
			return false
		},
		DeserializeErrorFunc: func(err error) *RuntimeError {
			return deserializeErrorFunc
		},
		RuntimeErrorFunc: func(err error) *RuntimeError {
			return runtimeErrorFunc
		},
	}
}

type server struct {
	Routers              map[string]routers
	Handlers             []*handler
	Call                 Call
	Scopeds              []any
	Singletons           []any
	Services             reflectMapFunc
	PersistentService    reflectMap
	DeserializeErrorFunc func(error) *RuntimeError
	RuntimeErrorFunc     func(error) *RuntimeError
}

func (s *server) Use(handlers ...any) IServer {

	s.Handlers = newHandlers(handlers, s)

	calles := moreHandler(s.Handlers)

	s.Call = func(param *handlerParam) (stop bool) {
		for _, call := range calles {
			return call(param)
		}
		return
	}

	return s
}

func (s *server) DeserializeError(cb func(error) *RuntimeError) {
	if cb != nil {
		s.DeserializeErrorFunc = cb
	}
}
