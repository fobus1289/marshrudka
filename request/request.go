package request

import (
	"net/http"
)

type IDeserialize interface {
	Json(in any) error
	Xml(in any) error
	FormData(in any) error
}

type IRequest interface {
	IParam
	IQueryParam
	IFormFile
	IDeserialize
	IHeader
}

type IFormFile interface {
	GetFile(formKey string) IFiles
	Files() IFileMap
}

type IParam interface {
	HasParam(key string) bool
	Param(key string) string
	ParamDefault(key string, _default string) string
	ParamInt(key string) int64
	ParamIntDefault(key string, _default int64) int64
	ParamFloat(key string) float64
	ParamFloatDefault(key string, _default float64) float64
	ParamIntClamp(key string, min, max int64) int64
	ParamFloatClamp(key string, min, max float64) float64
	TryGetParam(key string, in any) bool
}

type IQueryParam interface {
	HasQuery(key string) bool
	Query(key string) string
	QueryDefault(key string, _default string) string
	QueryInt(key string) int64
	QueryIntDefault(key string, _default int64) int64
	QueryFloat(key string) float64
	QueryFloatDefault(key string, _default float64) float64
	QueryIntClamp(key string, min, max int64) int64
	QueryFloatClamp(key string, min, max float64) float64
	TryQueryGetQuery(key string, in any) bool
}

type IHeader interface {
	GetHeader(key string) string
	Authorization() string
}

type IAuth interface {
	User(user any) error
	Header(header IHeader) error
}

type request struct {
	*http.Request
	ParamMap  map[string]string
	FormFiles fileMap
}

func NewRequest(r *http.Request, paramMap map[string]string) IRequest {
	return &request{
		Request:  r,
		ParamMap: paramMap,
	}
}
