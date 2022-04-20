package response

import "net/http"

type IResponse interface {
	Ok(status int) ISend
	Abort(status int) ISend
	File() IFile
}

type ISend interface {
	Prepare() IPrepare
	Json(data interface{}) ISend
	Text(data interface{}) ISend
	Html(data interface{}) ISend
	Xml(data interface{}) ISend
}

type IPrepare interface {
	SetStatus(status int)
	GetStatusCode() int
	GetContentType() string
	SetContentType(contentType string)
	GetData() interface{}
	SetData(data interface{})
	Marshal() []byte
	GetMarshal() func(data interface{}) []byte
	SetMarshal(func(data interface{}) []byte)
	IsAbort() bool
}

type IGeneralResponse interface {
	IResponse
	ISend
	IPrepare
}

type IFile interface {
	ServeFile(name string) IServerFile
	Stream(data interface{}) IServerFile
	Name(filename string) IFile
	Download() IFile
}

type IServerFile interface {
	Send(w http.ResponseWriter, r *http.Request)
}
