package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/fobus1289/marshrudka/request"
)

type reflectMap map[reflect.Type]reflect.Value
type reflectMapFunc map[reflect.Type]paramFunc
type httpMethodMap map[string]bool
type paramFunc func(*handlerParam) (reflect.Value, *RuntimeError)
type paramFuncs []paramFunc
type Call func(param *handlerParam) (stop bool)

type RuntimeError struct {
	Error       error
	Status      int
	ContentType string
	Data        []byte
}

var regexps = map[string]*regexp.Regexp{
	"float":   regexp.MustCompile(`^([-|+]?\d+([.]\d+))$`),
	"float?":  regexp.MustCompile(`^([-|+]?\d+([.]\d+))?$`),
	"int":     regexp.MustCompile(`^([-|+]?\d+)$`),
	"int?":    regexp.MustCompile(`^([-|+]?\d+)?$`),
	"number":  regexp.MustCompile(`^([-|+]?\d+([.]\d+)?)$`),
	"number?": regexp.MustCompile(`^([-|+]?\d+([.]\d+)?)?$`),
	"string":  regexp.MustCompile(`(\S+)`),
	"string?": regexp.MustCompile(`(\S+)?`),
}

type IRequestData interface {
}

func (pf *paramFuncs) Add(p paramFunc) {
	*pf = append(*pf, p)
}

func (pf paramFuncs) Get(h *handlerParam) ([]reflect.Value, *RuntimeError) {
	var values []reflect.Value

	for _, fn := range pf {
		value, err := fn(h)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}

func newRouter(handlers []any, s *server, httpMethods []string, actionPath string) *router {

	if len(handlers) == 0 {
		panic("handlers can be empy")
	}

	var methods = map[string]bool{
		http.MethodGet:    true,
		http.MethodHead:   true,
		http.MethodPost:   true,
		http.MethodPut:    true,
		http.MethodPatch:  true,
		http.MethodDelete: true,
	}

	httpMethod := httpMethodMap{}

	for _, method := range httpMethods {

		method = strings.ToUpper(method)

		if !methods[method] {
			panic("http method not supported " + method)
		}

		httpMethod[method] = true
	}

	validators, urlPath := parseUrl(actionPath)

	route := &router{
		Path:              urlPath,
		Paths:             strings.Split(urlPath, "/"),
		Call:              perpareActionHandler(newHandlers(handlers, s), s),
		Services:          s.Services,
		HttpUrlValidators: validators,
	}

	for _, method := range httpMethods {
		s.Routers[method] = append(s.Routers[method], route)
		sort.Sort(s.Routers[method])
	}

	return route
}

func newHandlers(actions []any, s *server) handlers {

	handlers := handlers{}

	for _, h := range actions {
		functionMeta := functionPareser(h, s)
		handlers = append(handlers, &handler{
			Params:   functionMeta.ParamFuncs,
			Call:     functionMeta.Call,
			OutTypes: functionMeta.OutTypes,
		})
	}

	return handlers
}

func parseUrl(urlPath string) ([]func(string) (bool, bool), string) {

	urlPath = strings.ReplaceAll(urlPath, " ", "")
	{
		urlPath = strings.TrimSuffix(urlPath, "/")
		urlPath = strings.TrimPrefix(urlPath, "/")
	}

	if urlPath == "" {
		return []func(string) (bool, bool){
			func(s string) (bool, bool) {
				return (s == "/"), false
			},
		}, "/"
	}

	reg := regexp.MustCompile(`({(\w+(:((\((\S+|\w+|\W+)?)?(\))|(\[(\S+|\w+|\W+)?\])|float|float\?|int|int\?|number|number\?|string\?|string)?)?)?})`)

	patterns := reg.FindAllString(urlPath, -1)

	paramValidator := map[string]func(string) (bool, bool){}

	cUrlPath := urlPath

	for _, pat := range patterns {

		if pat == "{}" {
			panic("{} error")
		}

		patTrim := strings.TrimSuffix(strings.TrimPrefix(pat, "{"), "}")

		pats := strings.Split(patTrim, ":")

		cUrlPath = strings.Replace(cUrlPath, pat, fmt.Sprintf(":%s", pats[0]), 1)

		pat = patTrim

		switch len(pats) {

		case 1:
			paramValidator[pat] = func(s string) (bool, bool) {
				return true, true
			}
		case 2:
			k, v := pats[0], pats[1]

			var param *regexp.Regexp

			if reg := regexps[v]; reg != nil {
				param = reg
			} else {
				param = regexp.MustCompile(v)
			}

			paramValidator[k] = func(s string) (bool, bool) {
				return param.MatchString(s), true
			}

		default:
			panic("default error")
		}
	}

	urlPaths := strings.Split(cUrlPath, "/")

	var validators = make([]func(string) (bool, bool), len(urlPaths))

	for index, urlpath := range urlPaths {

		urlpath = strings.TrimPrefix(urlpath, ":")

		if validator := paramValidator[urlpath]; validator != nil {
			validators[index] = validator
		} else {
			p := urlpath
			validators[index] = func(s string) (bool, bool) {
				return p == s, false
			}
		}
	}

	return validators, cUrlPath
}

func Elem(value reflect.Value) reflect.Value {

	for value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		value = value.Elem()
	}

	return value
}

func CopyValue(t reflect.Type, v reflect.Value) reflect.Value {

	newValue := reflect.New(t)

	newValue.Elem().Set(v.Elem())

	return newValue
}

func preparationOut(outType reflect.Type, s *server) (func(reflect.Value) []byte, string) {

	var ptrCount int

	if outType.Kind() == reflect.Ptr && outType.Implements(jwtUserType) {
		return func(v reflect.Value) []byte {
			jwtUser := v.Interface().(request.IJwtUser)
			token, _ := s.Jwt.Encode(jwtUser)
			data, _ := json.Marshal(jwtUser.Out(token))
			return data
		}, "application/json; charset=utf-8"
	}

	for outType.Kind() == reflect.Ptr {
		outType = outType.Elem()
		ptrCount++
	}

	switch outType.Kind() {
	case reflect.Invalid,
		reflect.Uintptr,
		reflect.Chan,
		reflect.Func,
		reflect.Pointer,
		reflect.UnsafePointer:
		panic("invalid wirte value")
	}

	var emptyData []byte

	var outFunc func(value reflect.Value) []byte
	var contentType string

	switch outType.Kind() {
	case reflect.Bool:
		if ptrCount == 0 {
			outFunc = func(value reflect.Value) []byte {
				return []byte([]byte(strconv.FormatBool(value.Bool())))
			}
		} else {
			outFunc = func(value reflect.Value) []byte {
				for i := 0; i < ptrCount; i++ {
					value = value.Elem()
				}
				return []byte([]byte(strconv.FormatBool(value.Bool())))
			}
		}

	case reflect.String:
		if ptrCount == 0 {
			outFunc = func(value reflect.Value) []byte {
				return []byte(value.String())
			}
		} else {
			outFunc = func(value reflect.Value) []byte {
				for i := 0; i < ptrCount; i++ {
					value = value.Elem()
				}
				return []byte(value.String())
			}
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if ptrCount == 0 {
			outFunc = func(value reflect.Value) []byte {
				return []byte(strconv.FormatInt(value.Int(), 10))
			}
		} else {
			outFunc = func(value reflect.Value) []byte {
				for i := 0; i < ptrCount; i++ {
					value = value.Elem()
				}
				return []byte(strconv.FormatInt(value.Int(), 10))
			}
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if ptrCount == 0 {
			outFunc = func(value reflect.Value) []byte {
				return []byte(strconv.FormatUint(value.Uint(), 10))
			}
		} else {
			outFunc = func(value reflect.Value) []byte {
				for i := 0; i < ptrCount; i++ {
					value = value.Elem()
				}
				return []byte(strconv.FormatUint(value.Uint(), 10))
			}
		}

	case reflect.Float32, reflect.Float64:
		if ptrCount == 0 {
			outFunc = func(value reflect.Value) []byte {
				return []byte(strconv.FormatFloat(value.Float(), 'f', -1, 64))
			}
		} else {
			outFunc = func(value reflect.Value) []byte {
				for i := 0; i < ptrCount; i++ {
					value = value.Elem()
				}
				return []byte(strconv.FormatFloat(value.Float(), 'f', -1, 64))
			}
		}

	case reflect.Complex64, reflect.Complex128:
		if ptrCount == 0 {
			outFunc = func(value reflect.Value) []byte {
				return []byte(strconv.FormatComplex(value.Complex(), 'f', -1, 128))
			}
		} else {
			outFunc = func(value reflect.Value) []byte {
				for i := 0; i < ptrCount; i++ {
					value = value.Elem()
				}
				return []byte(strconv.FormatComplex(value.Complex(), 'f', -1, 128))
			}
		}

	case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
		outFunc = func(value reflect.Value) []byte {
			outValue := value.Interface()
			{
				if outValue == nil {
					return emptyData
				}
			}
			if data, err := json.Marshal(outValue); err == nil {
				return data
			}
			return emptyData
		}

	}

	switch outType.Kind() {

	case
		reflect.Bool,
		reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		contentType = "text/plain; charset=utf-8"
	case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
		contentType = "application/json; charset=utf-8"
	}

	return outFunc, contentType
}

func BitSize(in any) int {
	switch in.(type) {
	case int8, uint8:
		return 8
	case int16, uint16:
		return 16
	case int32, uint32, float32, int, uint:
		return 32
	case int64, uint64, float64:
		return 64
	}
	return -1
}
