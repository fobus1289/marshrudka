package router

import (
	"net/http"
)

type IServer interface {
	IRouter
	IUse
	IDependency
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	HasClient() bool
	Run(addr string) error
	RunTLS(addr, certFile, keyFile string) error
	RunAsync(addr string) error
	RunAsyncTLS(addr, certFile, keyFile string) error
	RuntimeError(handler func(err error) interface{})
	BodyParseError(handler func() interface{})
}

type IGroup interface {
	IUse
	IRouter
}

type IRouter interface {
	Group(path string, handlers ...interface{}) IGroup
	GET(path string, handlers ...interface{}) IMatch
	FileServer(method, path, dir string, handlers ...interface{})
	POST(path string, handlers ...interface{}) IMatch
	PUT(path string, handlers ...interface{}) IMatch
	PATCH(path string, handlers ...interface{}) IMatch
	DELETE(path string, handlers ...interface{}) IMatch
	ANY(path string, handlers ...interface{}) IMatch
	MATCH(path string, methods []string, handlers ...interface{}) IMatch
}

type IMatch interface {
	Where(key, pattern string) IMatch
	WhereIn(pattern map[string]string) IMatch
	StripPrefix(prefix string) IMatch
}

type IUse interface {
	Use(handlers ...interface{})
}

type IDependency interface {
	SetService(v interface{}) bool
	GetService(out interface{}) bool
	SetServices(services ...interface{}) bool
	GetServices(services ...interface{}) bool
	FillServiceFields(service interface{}) bool
	FillServicesFields(services ...interface{}) bool
}

type IService interface {
	Constructor(services ...IService)
}
