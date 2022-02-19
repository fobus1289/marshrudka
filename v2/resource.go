package v2

import "net/http"

const (
	contentJson = "application/json; charset=utf-8"
	contentText = "text/plain"
	contentHtml = "text/html"
	contentXml  = "application/xml"
)

type Request struct {
	Response http.ResponseWriter
	*http.Request
	params map[string]string
}

type throw struct {
	status      int
	contentType string
	data        interface{}
}

type response struct {
	status      int
	contentType string
	data        interface{}
}

func (r *response) Error(status int) ISend {
	return &throw{
		status: status,
	}
}

func (r *response) Ok(status int) ISend {
	r.status = status
	return r
}

func Response() IResponse {
	return &response{}
}

func (t *throw) Json(data interface{}) ISend {
	t.contentType = contentJson
	t.data = data
	return t
}

func (t *throw) Text(data interface{}) ISend {
	t.contentType = contentText
	t.data = data
	return t
}

func (t *throw) Html(data interface{}) ISend {
	t.contentType = contentHtml
	t.data = data
	return t
}

func (t *throw) Xml(data interface{}) ISend {
	t.contentType = contentXml
	t.data = data
	return t
}

func (r *response) Json(data interface{}) ISend {
	r.contentType = contentJson
	r.data = data
	return r
}

func (r *response) Text(data interface{}) ISend {
	r.contentType = contentText
	r.data = data
	return r
}

func (r *response) Html(data interface{}) ISend {
	r.contentType = contentHtml
	r.data = data
	return r
}

func (r *response) Xml(data interface{}) ISend {
	r.contentType = contentXml
	r.data = data
	return r
}
