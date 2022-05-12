package request

import (
	"encoding/json"
	"encoding/xml"
	"github.com/gorilla/schema"
	"net/http"
	"reflect"
	"strconv"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

var (
	decoder = schema.NewDecoder()
)

func NewBodyParser(p reflect.Type, r *http.Request) IRequestParser {
	return &bodyParser{
		Type:    p,
		Request: r,
	}
}

type bodyParser struct {
	Type    reflect.Type
	Request *http.Request
}

func (r *bodyParser) getValue() reflect.Value {
	if r.Type.Kind() == reflect.Ptr {
		return reflect.New(r.Type.Elem())
	}
	return reflect.New(r.Type)
}

func (r *bodyParser) Json() reflect.Value {

	var value = r.getValue()

	if value.Kind() == reflect.Invalid {
		return value
	}

	if err := json.NewDecoder(r.Request.Body).Decode(value.Interface()); err != nil {
		return reflect.Value{}
	}

	if r.Type.Kind() == reflect.Ptr {
		return value
	}

	return value.Elem()
}

func (r *bodyParser) Xml() reflect.Value {

	var value = r.getValue()

	if value.Kind() == reflect.Invalid {
		return value
	}

	if err := xml.NewDecoder(r.Request.Body).Decode(value.Interface()); err != nil {
		return reflect.Value{}
	}

	if r.Type.Kind() == reflect.Ptr {
		return value
	}

	return value.Elem()
}

func (r *bodyParser) Form() reflect.Value {

	var (
		value   = r.getValue()
		element = value.Elem()
	)

	switch element.Kind() {
	case reflect.Map, reflect.Struct:
	case reflect.Interface:
		if element.NumMethod() != 0 {
			return reflect.Value{}
		}
	default:
		return reflect.Value{}
	}

	_ = r.Request.ParseMultipartForm(defaultMaxMemory)

	if len(r.Request.PostForm) == 0 {
		return reflect.Value{}
	}

	var (
		form = r.Request.PostForm
	)

	if element.Kind() == reflect.Struct {
		if err := decoder.Decode(value.Interface(), form); err != nil {
			return reflect.Value{}
		}

		if r.Type.Kind() == reflect.Ptr {
			return value
		} else {
			return value.Elem()
		}
	}

	var out interface{}
	switch element.Interface().(type) {
	case map[string]string:
		var m = map[string]string{}
		for k, v := range form {
			m[k] = v[0]
		}
		out = m
	case map[string]interface{}:
		var m = map[string]interface{}{}
		for k, v := range form {
			m[k] = v[0]
		}
		out = m
	default:
		var m = map[string]interface{}{}
		for k, v := range form {
			m[k] = v[0]
		}
		out = m
	}

	element.Set(reflect.ValueOf(out))

	if r.Type.Kind() == reflect.Ptr {

		return value
	}

	return element
}

func convertForm(field reflect.Value, convertVale string) {

	var result reflect.Value

	switch field.Interface().(type) {
	case bool:
		val, _ := strconv.ParseBool(convertVale)
		field.SetBool(val)
		return
	case uint, uint8, uint16, uint32, uint64:
		val, _ := strconv.ParseUint(convertVale, 10, 64)
		field.SetUint(val)
		return
	case int, int8, int16, int32, int64:
		val, _ := strconv.ParseInt(convertVale, 10, 64)
		field.SetInt(val)
		return
	case float64, float32:
		val, _ := strconv.ParseFloat(convertVale, 64)
		field.SetFloat(val)
		return
	case complex64, complex128:
		var val, _ = strconv.ParseComplex(convertVale, 128)
		field.SetComplex(val)
		return
	case string:
		field.SetString(convertVale)
		return
	case *string:
		result = reflect.ValueOf(&convertVale)
	case *bool:
		var boolRef, _ = strconv.ParseBool(convertVale)
		result = reflect.ValueOf(&boolRef)
	case *uint8:
		if uint64Ref, err := strconv.ParseUint(convertVale, 10, 8); err == nil {
			var uint8Ref = uint8(uint64Ref)
			result = reflect.ValueOf(&uint8Ref)
		}
	case *uint16:
		if uint64Ref, err := strconv.ParseUint(convertVale, 10, 16); err == nil {
			var uint16Ref = uint16(uint64Ref)
			result = reflect.ValueOf(&uint16Ref)
		}
	case *uint:
		if uint64Ref, err := strconv.ParseUint(convertVale, 10, 32); err == nil {
			var uintRef = uint(uint64Ref)
			result = reflect.ValueOf(&uintRef)
		}
	case *uint32:
		if uint64Ref, err := strconv.ParseUint(convertVale, 10, 32); err == nil {
			var uint32Ref = uint32(uint64Ref)
			result = reflect.ValueOf(&uint32Ref)
		}
	case *uint64:
		if uint64Ref, err := strconv.ParseUint(convertVale, 10, 64); err == nil {
			result = reflect.ValueOf(&uint64Ref)
		}
	case *int8:
		if int64Ref, err := strconv.ParseUint(convertVale, 10, 8); err == nil {
			var int8Ref = int8(int64Ref)
			result = reflect.ValueOf(&int8Ref)
		}
	case *int16:
		if int64Ref, err := strconv.ParseUint(convertVale, 10, 16); err == nil {
			var int16Ref = int16(int64Ref)
			result = reflect.ValueOf(&int16Ref)
		}
	case *int:
		if int64Ref, err := strconv.ParseUint(convertVale, 10, 32); err == nil {
			var int32Ref = int(int64Ref)
			result = reflect.ValueOf(&int32Ref)
		}
	case *int32:
		if int64Ref, err := strconv.ParseUint(convertVale, 10, 32); err == nil {
			var int32Ref = int32(int64Ref)
			result = reflect.ValueOf(&int32Ref)
		}

	case *int64:
		if uint64Ref, err := strconv.ParseUint(convertVale, 10, 64); err == nil {
			result = reflect.ValueOf(&uint64Ref)
		}
	case *float32:
		if float64Ref, err := strconv.ParseFloat(convertVale, 32); err == nil {
			var float32Ref = float32(float64Ref)
			result = reflect.ValueOf(&float32Ref)
		}

	case *float64:
		if float64Ref, err := strconv.ParseFloat(convertVale, 64); err == nil {
			result = reflect.ValueOf(&float64Ref)
		}
	case *complex64:
		if complex128Ref, err := strconv.ParseComplex(convertVale, 64); err == nil {
			var complex64Ref = complex64(complex128Ref)
			result = reflect.ValueOf(&complex64Ref)
		}
	case *complex128:
		if complex128Ref, err := strconv.ParseComplex(convertVale, 128); err == nil {
			result = reflect.ValueOf(&complex128Ref)
		}
	}
	if result.Kind() != reflect.Invalid {
		field.Set(result)
	}
}
