package router

import (
	"github.com/fobus1289/marshrudka/router/response"
	"net/http"
	"sync"
)

func Response() response.IResponse {
	return response.Response()
}

func NewServer() IServer {
	return &server{
		shutdown: &shutdown{
			Mutex:   &sync.Mutex{},
			request: make(chan *http.Request, 50),
		},
		clients:           0,
		services:          reflectMap{},
		HandlersInterface: handlersInterface{},
		routers:           routers{},
		runtimeError: func(err error) interface{} {
			return whatWentWrongErr.Error()
		},
		bodyEOF: func() interface{} {
			return string(emptyBody)
		},
	}
}
