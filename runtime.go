package marshrudka

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var (
	httpResponse     = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	httpRequest      = reflect.TypeOf(&http.Request{})
	_throw           = reflect.TypeOf(&Throw{})
	_response        = reflect.TypeOf(&response{})
	_request         = reflect.TypeOf(&Request{})
	methodNotAllowed = []byte("method not allowed")
	expectsJSON      = []byte("expects to receive a JSON object")
)

func (d *Drive) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {

	var isNotFound = true

	for _, r := range d.routers {
		route := r
		matches := route.uri.FindStringSubmatch(request.URL.Path)

		if len(matches) < 1 {
			continue
		}

		isNotFound = false

		if request.Method == http.MethodOptions {
			return
		}

		if !route.method["ANY"] && !route.method[request.Method] {
			var hasRoute bool

			for _, r2 := range d.routers {
				_matches := r2.uri.FindStringSubmatch(request.URL.Path)
				if len(_matches) > 0 && r2.method[request.Method] {
					route = r2
					matches = _matches
					hasRoute = true
				}
			}

			if !hasRoute {
				responseWriter.WriteHeader(405)
				_, _ = responseWriter.Write(methodNotAllowed)
				return
			}

		}

		responseWriter.Header().Set("Cache-Control", "no-cache")
		responseWriter.Header().Set("Accept-Encoding", "gzip, deflate, br")
		responseWriter.Header().Set("Content-Type", TEXT_HTML)

		var params = map[reflect.Type]reflect.Value{}

		for i, a := range route.actions {

			var values []reflect.Value

			for _, param := range a.Params {

				mapParam := params[param]

				if mapParam.Kind() != reflect.Invalid {
					values = append(values, mapParam)
					continue
				}

				if reflect.DeepEqual(param, httpRequest) {
					value := reflect.ValueOf(request)
					params[param] = value
					values = append(values, value)
					continue
				}

				if reflect.DeepEqual(param, httpResponse) {
					value := reflect.ValueOf(responseWriter)
					params[param] = value
					values = append(values, value)
					continue
				}

				if reflect.DeepEqual(param, _request) {

					req := &Request{
						HttpResponseWriter: responseWriter,
						HttpRequest:        request,
						Params:             map[string]string{},
					}

					if len(r.params) > 0 {
						var paramIndex int

						for _, match := range matches {
							if strings.HasPrefix(match, "/") || strings.TrimSpace(match) == "" {
								continue
							}
							req.Params[r.params[paramIndex]] = match
							paramIndex++
						}
					}

					value := reflect.ValueOf(req)
					params[param] = value

					values = append(values, value)

					continue
				}

				for key, value := range d.services {
					if reflect.DeepEqual(param, key) {
						if param.Kind() == reflect.Ptr {
							values = append(values, value)
							params[key] = value
						} else {
							values = append(values, value)
							params[key] = value
						}
						break
					}
				}

				if !isPrimitive(param.Kind()) && d.services[param].Kind() == reflect.Invalid {
					value := setOther(param, request, responseWriter)
					if value == nil {
						return
					}
					params[param] = *value
					values = append(values, *value)
					continue
				}

			}

			ret := a.Call(values)

			for _, value := range ret {

				valueOf := reflect.ValueOf(value.Interface())

				if reflect.DeepEqual(valueOf.Type(), _throw) {
					var throw = valueOf.Interface().(*Throw)
					responseWriter.Header().Set("Content-Type", throw.ContentType)
					responseWriter.WriteHeader(throw.StatusCode)

					_, _ = responseWriter.Write(getPrimitiveResult(reflect.ValueOf(throw.Data)))
					return
				}

				if len(r.actions)-1 == i {

					if reflect.DeepEqual(valueOf.Type(), _response) {
						var response = valueOf.Interface().(*response)
						responseWriter.Header().Set("Content-Type", response.ContentType)
						responseWriter.WriteHeader(response.StatusCode)
						_, _ = responseWriter.Write(getPrimitiveResult(reflect.ValueOf(response.Data)))
						return
					}

					_, _ = responseWriter.Write(getPrimitiveResult(valueOf))
				}

				params[valueOf.Type()] = valueOf
			}

		}

		break
	}

	if isNotFound {
		http.NotFound(responseWriter, request)
	}

}

func setOther(param reflect.Type, request *http.Request, w http.ResponseWriter) *reflect.Value {

	_value := reflect.New(param)

	err := json.NewDecoder(request.Body).Decode(_value.Interface())

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(expectsJSON)
		log.Println(err)
		return nil
	}

	val := _value.Elem()

	return &val
}

func getParamValue(kind reflect.Kind, value string) reflect.Value {
	switch kind {
	case reflect.Bool:
		out, _ := strconv.ParseBool(value)
		return reflect.ValueOf(out)
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		out, _ := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(out)
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		out, _ := strconv.ParseUint(value, 10, 64)
		return reflect.ValueOf(out)
	case reflect.Float32,
		reflect.Float64:
		out, _ := strconv.ParseFloat(value, 10)
		return reflect.ValueOf(out)
	}

	return reflect.ValueOf(value)
}

func getPrimitiveResult(value reflect.Value) []byte {

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

func Equals(a, b reflect.Type) bool {

	if a.Kind() == reflect.Ptr {
		a = a.Elem()
	}

	if b.Kind() == reflect.Ptr {
		b = b.Elem()
	}

	return reflect.DeepEqual(a, b)
}
