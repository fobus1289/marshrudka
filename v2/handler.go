package v2

import (
	"encoding/json"
	"net/http"
	"reflect"
)

type handler struct {
	Server *server
	Params []reflect.Type
	Router *router
	Last   bool
	Call   func([]reflect.Value) []reflect.Value
}

type handlers []*handler

func (hs *handlers) AddRange(_handlers handlers) {
	if hs == nil {
		*hs = handlers{}
	}

	*hs = append(*hs, _handlers...)
}

func (hs handlers) Next(w http.ResponseWriter, r *http.Request, pm paramsMap) bool {

	for _, h := range hs {
		if !h.Next(w, r, pm) {
			return false
		}
	}

	return true
}

func (hs handlers) SetRouter(r *router) {
	for _, h := range hs {
		h.Router = r
	}
}

func (h *handler) Next(w http.ResponseWriter, r *http.Request, pm paramsMap) bool {

	var params []reflect.Value

	for _, param := range h.Params {

		if prm := pm[param]; prm.Kind() != reflect.Invalid {
			params = append(params, prm)
			continue
		}

		service := h.Server.services[param]

		if service.Kind() != reflect.Invalid {
			pm[service.Type()] = service
			params = append(params, service)
			continue
		}

		switch param {
		case httpRes:
			var res = reflect.ValueOf(w)
			pm[res.Type()] = res
			params = append(params, res)
			continue
		case httpReq:
			var req = reflect.ValueOf(r)
			pm[req.Type()] = req
			params = append(params, req)
			continue
		case request:
			var _request = h.getRequestParam(w, r)
			var req = reflect.ValueOf(_request)
			pm[req.Type()] = req
			params = append(params, req)
			continue
		}

		if object := read(param, r); object.Kind() != reflect.Invalid {
			pm[object.Type()] = object
			params = append(params, object)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(expectsJSON)
			return false
		}

	}

	var ret = h.Call(params)

	if stop := h.stop(ret, w, r); stop {
		return false
	}

	for _, value := range ret {
		if value.Kind() == reflect.Invalid {
			continue
		}
		pm[value.Type()] = value
	}

	return true
}

func (h *handler) stop(ret []reflect.Value, w http.ResponseWriter, r *http.Request) (stop bool) {

	if len(ret) < 1 && h.Last {
		return true
	}

	for _, value := range ret {
		switch t := value.Interface().(type) {
		case *throw:
			w.Header().Set("Content-Type", t.contentType)
			w.WriteHeader(t.status)
			_, _ = w.Write(write(t.data))
			return true
		case *response:
			if h.Last {
				w.Header().Set("Content-Type", t.contentType)
				w.WriteHeader(t.status)
				_, _ = w.Write(write(t.data))
				return true
			}
		}
	}

	if h.Last {
		var retValue = ret[0]
		if retValue.Kind() != reflect.Invalid {
			_, _ = w.Write(write(retValue.Interface()))
			return true
		}
		return true
	}

	return false
}

func (h *handler) getRequestParam(w http.ResponseWriter, r *http.Request) *Request {
	var _request = &Request{
		Response: w,
		Request:  r,
		params:   map[string]string{},
	}

	var route = h.Router
	var params = route.Params

	if len(params) < 1 {
		return _request
	}

	var httpParams = h.Router.Match.FindStringSubmatch(r.URL.Path)

	if len(httpParams) <= 2 {
		return _request
	}

	httpParams = httpParams[2:]

	if len(httpParams) != len(params) {
		return _request
	}

	for i, param := range params {
		_request.params[param] = httpParams[i]
	}

	return _request
}

func write(data interface{}) []byte {
	buff, err := json.Marshal(data)

	if err != nil {
		return nil
	}

	return buff
}

func read(param reflect.Type, request *http.Request) reflect.Value {

	bodyObject := reflect.New(param)

	if err := json.NewDecoder(request.Body).Decode(bodyObject.Interface()); err != nil {
		return reflect.Value{}
	}

	return bodyObject.Elem()
}
