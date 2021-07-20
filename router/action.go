package router

import (
	"net/http"
	"reflect"
	"strconv"
)

type _http struct {
	code int
}

type Throw struct {
	StatusCode  int
	ContentType string
	Data        interface{}
}

type Data struct {
	StatusCode         int
	ContentType        string
	ContentDisposition string
	Data               interface{}
}

type File struct {
	StatusCode         int
	ContentType        string
	ContentDisposition string
	Path               string
	Name               string
}

func (f *File) stream(w http.ResponseWriter, r *http.Request) {
	if f.ContentDisposition != "" {
		w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(f.Name))
	}
	if f.ContentType != "" {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeFile(w, r, f.Path)
}

func (f *File) Download() *File {
	f.ContentDisposition = "attachment; filename=" + strconv.Quote(f.Name)
	f.ContentType = "application/octet-stream"
	return f
}

func Response(code int) *_http {
	return &_http{
		code: code,
	}
}

func (h *_http) Throw() *Throw {
	return &Throw{
		StatusCode: h.code,
	}
}

func (h *_http) Json(data interface{}) *Data {
	return &Data{
		StatusCode:  h.code,
		ContentType: "application/json; charset=utf-8",
		Data:        data,
	}
}

func (h *_http) Text(data interface{}) *Data {
	return &Data{
		StatusCode:  h.code,
		ContentType: "text/plain; charset=utf-8",
		Data:        data,
	}
}

func (h *_http) Html(data interface{}) *Data {
	return &Data{
		StatusCode:  h.code,
		ContentType: "text/html; charset=utf-8",
		Data:        data,
	}
}

func (h *_http) File(path string, name string) *File {
	return &File{
		StatusCode: h.code,
		Path:       path,
		Name:       name,
	}
}

func (h *_http) Stream(filename string, data interface{}) *Data {
	return &Data{
		StatusCode:         h.code,
		ContentType:        "application/octet-stream",
		ContentDisposition: "attachment; filename=" + strconv.Quote(filename),
		Data:               data,
	}
}

func (t *Throw) Json(data interface{}) *Throw {
	return &Throw{
		ContentType: "application/json; charset=utf-8",
		Data:        data,
	}
}

func (t *Throw) Text(data interface{}) *Throw {
	return &Throw{
		ContentType: "text/plain; charset=utf-8",
		Data:        data,
	}
}

func (t *Throw) Html(data interface{}) *Throw {
	return &Throw{
		ContentType: "text/html; charset=utf-8",
		Data:        data,
	}
}

func (d *Drive) Dep(owner interface{}) {

	ownerType := reflect.ValueOf(owner)

	if ownerType.Kind() == reflect.Ptr {
		ownerType = ownerType.Elem()
	}

	for i := 0; i < ownerType.NumField(); i++ {
		fieldType := ownerType.Field(i)

		service := d.services[fieldType.Type()]

		if service.Kind() != reflect.Invalid {
			fieldType.Set(service)
		}
	}
}
