package v2

import "net/http"

type IServer interface {
	IRouter
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Run(addr string) error
	RunTLS(addr, certFile, keyFile string) error
	RunAsync(addr string) error
	RunAsyncTLS(addr, certFile, keyFile string) error
}

type IRouter interface {
	Group(path string, handlers ...interface{}) IRouter
	GET(path string, handlers ...interface{})
	POST(path string, handlers ...interface{})
	PUT(path string, handlers ...interface{})
	PATCH(path string, handlers ...interface{})
	DELETE(path string, handlers ...interface{})
	ANY(path string, handlers ...interface{})
	MATCH(path string, methods []string, handlers ...interface{})
}

type IUse interface {
	Use(handlers ...interface{})
}

type IResponse interface {
	Error(status int) ISend
	Ok(status int) ISend
}

type ISend interface {
	Json(data interface{}) ISend
	Text(data interface{}) ISend
	Html(data interface{}) ISend
	Xml(data interface{}) ISend
}
