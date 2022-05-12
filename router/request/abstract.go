package request

import (
	"io"
	"net/http"
	"os"
	"reflect"
)

type files map[string]IFileContainer

type IFiles interface {
	Get(name string) IFileContainer
	Len() int
}

type IFormFile interface {
	Get(formKey string) IFileContainer
	Files() IFiles
}

type IFileInfo interface {
	Size() int64
	Name() string
	ContentType() string
	Extension() string
}

type IFile interface {
	Read(writer io.Writer) IFile
	Store(dir string, storagePath *string) IFile
	Rollback() IFile
	RandomFileName() IFile
	Info() IFileInfo
	SetNewName(name string) IFile
	GetNewName() string
	SetPrem(perm os.FileMode) IFile
	IsValid() bool
	Error() error
}

type IFileContainer interface {
	Files() []IFile
	GetFirst() IFile
	Errors() []error
	RollbackAll() IFileContainer
	StoreAll(dir string, storagePaths *[]string) IFileContainer
	Count() int
	Has() bool
	HasMultiple() bool
}

type IRequestParser interface {
	Json() reflect.Value
	Xml() reflect.Value
	Form() reflect.Value
}

type IRequest interface {
	IParam
	IQueryParam
	IBody
	FormFile() IFormFile
	Request() *http.Request
	Response() http.ResponseWriter
}

type IParam interface {
	Param(key string) String
	HasParam(key string) bool
	TryGetParam(key string, in interface{}) bool
}

type IQueryParam interface {
	Query(key string) String
	HasQuery(key string) bool
	TryGetQuery(key string, in interface{}) bool
}

type IBody interface {
	Json(interface{}) error
	Xml(interface{}) error
}
