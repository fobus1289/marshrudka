package router

import (
	"encoding/json"
	request2 "github.com/fobus1289/marshrudka/router/request"
	"github.com/fobus1289/marshrudka/router/response"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type handler struct {
	Server     *server
	Params     []reflect.Type
	ParseParam map[reflect.Type]func(http.ResponseWriter, *http.Request, *handler) reflect.Value
	Router     *router
	Last       bool
	Call       func([]reflect.Value) []reflect.Value
}

type handlers []*handler

func (hs *handlers) AddRange(_handlers handlers) {
	if hs == nil {
		*hs = handlers{}
	}

	*hs = append(*hs, _handlers...)
}

func (hs handlers) SetRouter(r *router) {
	for _, h := range hs {
		h.Router = r
	}
}

func (hs handlers) Next(w http.ResponseWriter, r *http.Request, pm reflectMap) bool {

	for _, h := range hs {
		if !h.Next(w, r, pm) {
			return false
		}
	}

	return true
}

func (h *handler) Next(w http.ResponseWriter, r *http.Request, pm reflectMap) bool {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			_err, _ := err.(error)
			{
				if _err == nil {
					_err = whatWentWrongErr
				}
			}
			_, _ = w.Write(writeJson(h.Server.runtimeError(_err)))
		}
	}()

	var params []reflect.Value

	for _, param := range h.Params {

		if value := pm[param]; value.IsValid() {
			params = append(params, value)
			continue
		}

		var fn = h.ParseParam[param]

		if fn == nil {
			continue
		}

		if value := fn(w, r, h); value.IsValid() {
			params = append(params, value)
			pm[param] = value
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(writeJson(h.Server.bodyEOF()))
			return false
		}

	}

	var ret = h.Call(params)

	if stop := h.stop(ret, w, r); stop {
		return false
	}

	for _, value := range ret {
		if value.Kind() == reflect.Interface {
			if value = value.Elem(); value.Kind() == reflect.Invalid {
				continue
			}
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
		case response.IPrepare:
			if t.IsAbort() || h.Last {
				w.Header().Set("Content-Type", t.GetContentType())
				w.WriteHeader(t.GetStatusCode())
				_, _ = w.Write(t.Marshal())
				return true
			}
		case response.IServerFile:
			if h.Last {
				t.Send(w, r)
				return true
			}
		}
	}

	if h.Last {
		var retValue = ret[0]
		if retValue.Kind() != reflect.Invalid {
			_, _ = w.Write(write(retValue))
		}
		return true
	}

	return false
}

func write(outValue reflect.Value) []byte {

	if outValue.Kind() == reflect.Ptr || outValue.Kind() == reflect.Interface {
		if outValue = outValue.Elem(); !outValue.IsValid() {
			return []byte{}
		}
	}

	var out = outValue.Interface()

	switch t := out.(type) {
	case string:
		return []byte(t)
	case bool:
		return []byte(strconv.FormatBool(t))
	case int8, int16, int32, int, int64:
		return []byte(strconv.FormatInt(outValue.Int(), 10))
	case uint8, uint16, uint32, uint, uint64:
		return []byte(strconv.FormatUint(outValue.Uint(), 10))
	case float32, float64:
		return []byte(strconv.FormatFloat(outValue.Float(), 'f', -1, 64))
	case complex64, complex128:
		return []byte(strconv.FormatComplex(outValue.Complex(), 'f', -1, 128))
	}

	return writeJson(out)
}

func writeJson(out interface{}) []byte {

	if out == nil {
		return []byte{}
	}

	buff, err := json.Marshal(out)

	if err != nil {
		return []byte{}
	}

	return buff
}

func read(param reflect.Type, req *http.Request) reflect.Value {

	var contentType = req.Header.Get("Content-Type")
	var requestParser = request2.NewBodyParser(param, req)

	if strings.HasPrefix(contentType, "multipart/form-data") ||
		strings.HasSuffix(contentType, "application/x-www-form-urlencoded") {
		return requestParser.Form()
	}

	switch contentType {
	case "application/json":
		return requestParser.Json()
	case "application/xml":
		return requestParser.Xml()
	default:
		return reflect.Value{}
	}

}
