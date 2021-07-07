package marshrudka

import (
	"net/http"
)

const (
	JSON       = "application/json; charset=utf-8"
	TEXT_HTML  = "text/html; charset=utf-8"
	TEXT_PLAIN = "text/plain; charset=utf-8"
)

type Stop bool

type throw struct {
	StatusCode  int
	Data        interface{}
	ContentType string
}

type response struct {
	StatusCode  int
	Data        interface{}
	ContentType string
}

type _http struct {
	code int
}

func Response(code int) *_http {
	return &_http{
		code: code,
	}
}

func (h *_http) Json(data interface{}) *response {
	return &response{
		StatusCode:  h.code,
		Data:        data,
		ContentType: JSON,
	}
}

func (h *_http) Throw(content string, data interface{}) *throw {
	return &throw{
		StatusCode:  h.code,
		Data:        data,
		ContentType: content,
	}
}

type Request struct {
	HttpResponseWriter http.ResponseWriter
	HttpRequest        *http.Request
	Params             map[string]string
}

func (r *Request) Query(key string) string {
	return r.HttpRequest.URL.Query().Get(key)
}

func (r *Request) Param(key string) string {
	return r.Params[key]
}
