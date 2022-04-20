package response

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
)

type file struct {
	status  int
	name    string
	headers map[string]string
	data    interface{}
}

type stream struct {
	headers map[string]string
	status  int
	data    interface{}
}

type serveFile struct {
	name    string
	headers map[string]string
}

func (s *serveFile) Send(w http.ResponseWriter, r *http.Request) {

	for k, v := range s.headers {
		w.Header().Set(k, v)
	}

	http.ServeFile(w, r, s.name)
}

func (s *stream) Send(w http.ResponseWriter, r *http.Request) {

	for k, v := range s.headers {
		w.Header().Set(k, v)
	}

	if buff, err := json.Marshal(s.data); err == nil {
		_, _ = w.Write(buff)
	}

}

func (f *file) ServeFile(name string) IServerFile {
	_, _file := filepath.Split(name)
	f.headers["Content-Disposition"] = "attachment; filename=" + strconv.Quote(_file)
	f.name = name
	return &serveFile{
		name:    name,
		headers: f.headers,
	}
}

func (f *file) Stream(data interface{}) IServerFile {
	f.data = data
	return &stream{
		status:  f.status,
		data:    data,
		headers: f.headers,
	}
}

func (f *file) Name(filename string) IFile {
	f.headers["Content-Disposition"] = "attachment; filename=" + strconv.Quote(filename)
	return f
}

func (f *file) Download() IFile {
	f.headers["Content-Type"] = "application/octet-stream"
	return f
}
