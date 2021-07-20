package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const invalid = reflect.Invalid

type reflectMap map[reflect.Type]reflect.Value

var (
	_httpRes         = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	_httpReq         = reflect.TypeOf(&http.Request{})
	_throw           = reflect.TypeOf(&Throw{})
	_data            = reflect.TypeOf(&Data{})
	_request         = reflect.TypeOf(&Request{})
	methodNotAllowed = []byte("method not allowed")
	expectsJSON      = []byte("expects to receive a JSON object")
)

type FormFile struct {
	Name   string
	Data   []byte
	Size   int64
	perm   os.FileMode
	Header textproto.MIMEHeader
}

func (f *FormFile) SetName(name string) *FormFile {
	f.Name = name
	return f
}

func (f *FormFile) SetPrem(perm os.FileMode) *FormFile {
	f.perm = perm
	return f
}

func (f *FormFile) Store(path string) error {
	path = fmt.Sprintf("%s/%s", strings.TrimSuffix(path, "/"), f.Name)
	return ioutil.WriteFile(path, f.Data, f.perm)
}

type Request struct {
	Response http.ResponseWriter
	*http.Request
	params map[string]string
}

func (r *Request) SetHeader(key, value string) {
	r.Response.Header().Set(key, value)
}

func (r *Request) GetHeader(key string) string {
	return r.Response.Header().Get(key)
}

func (r *Request) Write(data interface{}) {
	_, _ = r.Response.Write(valueBytes(reflect.ValueOf(data)))
}

func (r *Request) Param(key string) string {
	return r.params[key]
}

func (r *Request) TryParamGet(key string, out *string) (ok bool) {
	value := r.params[key]

	if value == "" {
		return false
	}
	out = &value
	return true
}

func (r *Request) TryParamGetInt(key string, out *int64) (ok bool) {
	value, err := strconv.ParseInt(r.params[key], 10, 64)

	if err != nil {
		return false
	}

	out = &value

	return true
}

func (r *Request) TryParamGetUInt(key string, out *uint64) (ok bool) {

	value, err := strconv.ParseUint(r.params[key], 10, 64)

	if err != nil {
		return false
	}

	out = &value

	return true
}

func (r *Request) TryParamGetFloat(key string, out *float64) (ok bool) {
	value, err := strconv.ParseFloat(r.params[key], 10)

	if err != nil {
		return false
	}

	out = &value

	return true
}

func (r *Request) TryParamGetBool(key string, out *bool) (ok bool) {

	value, err := strconv.ParseBool(r.params[key])

	if err != nil {
		return false
	}

	out = &value

	return true
}

func (r *Request) ParamGetInt(key string) (value int64) {
	value, _ = strconv.ParseInt(r.params[key], 10, 64)
	return value
}

func (r *Request) ParamGetUInt(key string) (value uint64) {
	value, _ = strconv.ParseUint(r.params[key], 10, 64)
	return value
}

func (r *Request) ParamGetFloat(key string) (value float64) {
	value, _ = strconv.ParseFloat(r.params[key], 10)
	return value
}

func (r *Request) ParamGetBool(key string) (value bool) {
	value, _ = strconv.ParseBool(r.params[key])
	return value
}

func (r *Request) Query(key string) string {
	return r.URL.Query().Get(key)
}

func (r *Request) TryQueryGet(key string, out *string) (ok bool) {

	value := r.URL.Query().Get(key)

	if value == "" {
		return false
	}
	out = &value
	return true
}

func (r *Request) TryQueryGetInt(key string, out *int64) (ok bool) {

	value, err := strconv.ParseInt(r.URL.Query().Get(key), 10, 64)

	if err != nil {
		return false
	}

	out = &value

	return true
}

func (r *Request) TryQueryGetUInt(key string, out *uint64) (ok bool) {

	value, err := strconv.ParseUint(r.URL.Query().Get(key), 10, 64)

	if err != nil {
		return false
	}

	out = &value

	return true
}

func (r *Request) TryQueryGetFloat(key string, out *float64) (ok bool) {
	value, err := strconv.ParseFloat(r.URL.Query().Get(key), 10)

	if err != nil {
		return false
	}

	out = &value

	return true
}

func (r *Request) TryQueryGetBool(key string, out *bool) (ok bool) {

	value, err := strconv.ParseBool(r.URL.Query().Get(key))

	if err != nil {
		return false
	}

	out = &value

	return true
}

func (r *Request) QueryGetInt(key string) (value int64) {
	value, _ = strconv.ParseInt(r.URL.Query().Get(key), 10, 64)
	return value
}

func (r *Request) QueryGetUInt(key string) (value uint64) {
	value, _ = strconv.ParseUint(r.URL.Query().Get(key), 10, 64)
	return value
}

func (r *Request) QueryGetFloat(key string) (value float64) {
	value, _ = strconv.ParseFloat(r.URL.Query().Get(key), 10)
	return value
}

func (r *Request) QueryGetBool(key string) (value bool) {
	value, _ = strconv.ParseBool(r.URL.Query().Get(key))
	return value
}

func (r *Request) FormFile(key string) (*FormFile, error) {

	_, file, err := r.Request.FormFile(key)

	if err != nil {
		return nil, err
	}

	var formFile = &FormFile{
		Name:   file.Filename,
		perm:   os.FileMode(0644),
		Data:   make([]byte, file.Size),
		Size:   file.Size,
		Header: file.Header,
	}

	fOpen, err := file.Open()

	if err != nil {
		return nil, err
	}

	defer func(fOpen multipart.File) {
		err = fOpen.Close()
	}(fOpen)

	_, err = fOpen.Read(formFile.Data)

	if err != nil {
		return nil, err
	}

	return formFile, nil
}

func isThrow(val reflect.Value, w http.ResponseWriter) bool {
	switch t := val.Interface().(type) {
	case Throw:
		data := valueBytes(reflect.ValueOf(t.Data))
		w.WriteHeader(t.StatusCode)
		w.Header().Set("Content-Type", t.ContentType)
		_, _ = w.Write(data)
		return true
	case *Throw:
		data := valueBytes(reflect.ValueOf(t.Data))
		w.WriteHeader(t.StatusCode)
		w.Header().Set("Content-Type", t.ContentType)
		_, _ = w.Write(data)
		return true
	}
	return false
}

func isResponse(val reflect.Value, w http.ResponseWriter) bool {
	switch t := val.Interface().(type) {
	case Data:
		data := valueBytes(reflect.ValueOf(t.Data))
		if t.ContentDisposition != "" {
			w.Header().Set("Content-Disposition", t.ContentDisposition)
		}
		if t.ContentType != "" {
			w.Header().Set("Content-Type", t.ContentType)
		}
		_, _ = w.Write(data)
		return true
	case *Data:
		data := valueBytes(reflect.ValueOf(t.Data))
		if t.ContentDisposition != "" {
			w.Header().Set("Content-Disposition", t.ContentDisposition)
		}
		if t.ContentType != "" {
			w.Header().Set("Content-Type", t.ContentType)
		}
		_, _ = w.Write(data)
		return true
	}
	return false
}

func isFileResponse(val reflect.Value, w http.ResponseWriter, r *http.Request) bool {
	switch t := val.Interface().(type) {
	case File:
		t.stream(w, r)
		return true
	case *File:
		t.stream(w, r)
		return true
	}
	return false
}

func valueBytes(value reflect.Value) []byte {

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	switch value.Kind() {

	case reflect.Bool:
		var boolBit = "false"
		if value.Bool() {
			boolBit = "true"
		}
		return []byte(boolBit)
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return []byte(strconv.FormatInt(value.Int(), 10))
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return []byte(strconv.FormatUint(value.Uint(), 10))
	case reflect.Float32,
		reflect.Float64:
		return []byte(strconv.FormatFloat(value.Float(), 'f', -1, 64))
	case reflect.String:
		return []byte(value.String())
	case reflect.Struct, reflect.Slice, reflect.Interface, reflect.Map:
		toByte, _ := json.Marshal(value.Interface())
		return toByte
	}
	return nil
}

func setOther(param reflect.Type, request *http.Request, w http.ResponseWriter) *reflect.Value {

	_value := reflect.New(param)

	err := json.NewDecoder(request.Body).Decode(_value.Interface())

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(expectsJSON)
		log.Println(err)
		return nil
	}

	val := _value.Elem()

	return &val
}
