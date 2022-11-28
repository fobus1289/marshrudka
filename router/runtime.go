package router

import (
	"errors"
	"log"
	"net/http"
)

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if err := recover(); err != nil {
			var currentErr error
			log.Println(err)
			switch e := err.(type) {
			case error:
				currentErr = e
			case string:
				currentErr = errors.New(e)
			}
			runtimeError := s.RuntimeErrorFunc(currentErr)
			w.Header().Add("Content-Type", runtimeError.ContentType)
			w.WriteHeader(runtimeError.Status)
			w.Write(runtimeError.Data)
		}
	}()

	sessionData := reflectMap{}

	routers := s.Routers[r.Method]

	if len(routers) == 0 {

		stop := s.Call(&handlerParam{
			ResponseWriter: w,
			Request:        r,
			SessionData:    sessionData,
			Services:       s.Services,
			HttpParamMap:   nil,
		})

		if !stop {
			http.NotFound(w, r)
		}

		return
	}

	router, status, params := routers.Find(w, r)

	handlerParam := &handlerParam{
		ResponseWriter: w,
		Request:        r,
		SessionData:    sessionData,
		Services:       s.Services,
		HttpParamMap:   params,
	}

	if s.Call(handlerParam) {
		return
	}

	if !status {
		http.NotFound(w, r)
		return
	}

	router.Call(handlerParam)
}
