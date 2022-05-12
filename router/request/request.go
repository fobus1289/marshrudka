package request

import (
	"net/http"
	"reflect"
	"strconv"
	"sync"
)

func NewRequest(w http.ResponseWriter, r *http.Request) IRequest {
	return &request{
		w:         w,
		r:         r,
		paramsMap: map[string]string{},
		Once:      &sync.Once{},
	}
}

type request struct {
	w         http.ResponseWriter
	r         *http.Request
	params    *Params
	paramsMap map[string]string
	ifs       IFormFile
	*sync.Once
}

func (r *request) FormFile() IFormFile {
	if r.ifs == nil {
		r.ifs = NewFormFile(r.w, r.r)
	}
	return r.ifs
}

func (r *request) Request() *http.Request {
	return r.r
}

func (r *request) Response() http.ResponseWriter {
	return r.w
}

func (r *request) setType(value string, in interface{}) bool {
	var inValue = reflect.ValueOf(in)

	if !inValue.IsValid() || inValue.Kind() != reflect.Ptr {
		return false
	}

	var inTmp = inValue.Elem()

	switch t := inTmp.Interface().(type) {
	case string:
		inTmp.SetString(value)
	case bool:
		var _bool, err = strconv.ParseBool(value)
		if err != nil {
			return false
		}
		inTmp.SetBool(_bool)
	case int8, int16, int32, int, int64:
		var _int, err = strconv.ParseInt(value, 10, getBitSize(t))
		if err != nil {
			return false
		}
		inTmp.SetInt(_int)
	case uint8, uint16, uint32, uint, uint64:
		var _uint, err = strconv.ParseUint(value, 10, getBitSize(t))
		if err != nil {
			return false
		}
		inTmp.SetUint(_uint)
	case float32, float64:
		var _float, err = strconv.ParseFloat(value, getBitSize(t))
		if err != nil {
			return false
		}
		inTmp.SetFloat(_float)
	case complex64, complex128:
		var _complex, err = strconv.ParseComplex(value, getBitSize(t))
		if err != nil {
			return false
		}
		inTmp.SetComplex(_complex)
	default:
		return false
	}

	return true
}

func getBitSize(in interface{}) int {
	switch in.(type) {
	case int8, uint8:
		return 8
	case int16, uint16:
		return 16
	case int32, uint32, float32, int, uint:
		return 32
	case complex128:
		return 182
	default:
		return 64
	}
}
