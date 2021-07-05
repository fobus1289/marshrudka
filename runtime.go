package marshrudka

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
)

var (
	httpResponse     = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	httpRequest      = reflect.TypeOf(&http.Request{})
	_throw           = reflect.TypeOf(&throw{})
	_response        = reflect.TypeOf(&response{})
	_request         = reflect.TypeOf(&Request{})
	methodNotAllowed = []byte("method not allowed")
)

func (d *drive) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {

	var isNotFound = true

	for _, r := range d.routers {

		thisRoute := r.uri.MatchString(request.URL.Path)

		if !thisRoute {
			continue
		}

		isNotFound = false

		if r.method != "ANY" && r.method != request.Method {
			responseWriter.WriteHeader(405)
			_, _ = responseWriter.Write(methodNotAllowed)
			return
		}

		var params = map[reflect.Type]reflect.Value{}

		for i, a := range r.actions {

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
					value := reflect.ValueOf(&Request{
						HttpResponseWriter: responseWriter,
						HttpRequest:        request,
					})
					params[param] = value
					values = append(values, value)
					continue
				}

				if !isPrimitive(param.Kind()) && d.services[param].Kind() == reflect.Invalid {
					value := *setOther(param, request)
					params[param] = value
					values = append(values, value)
					continue
				}

				for key, value := range d.services {
					if reflect.DeepEqual(param, key) {
						values = append(values, value)
						params[key] = value
						break
					}
				}

			}

			ret := a.Call(values)

			for _, value := range ret {

				valueOf := reflect.ValueOf(value.Interface())

				if reflect.DeepEqual(valueOf.Type(), _throw) {
					var throw = valueOf.Interface().(*throw)
					responseWriter.Header().Set("Content-Type", throw.ContentType)
					responseWriter.WriteHeader(throw.StatusCode)
					data, _ := json.Marshal(throw.Data)
					_, _ = responseWriter.Write(data)
					return
				}

				if len(r.actions)-1 == i {
					var res = valueOf.Interface()

					if reflect.DeepEqual(valueOf.Type(), _response) {
						var response = valueOf.Interface().(*response)
						responseWriter.Header().Set("Content-Type", response.ContentType)
						responseWriter.WriteHeader(response.StatusCode)
						res = response.Data
					}

					data, _ := json.Marshal(res)
					_, _ = responseWriter.Write(data)
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

func setOther(param reflect.Type, request *http.Request) *reflect.Value {

	_value := reflect.New(param)

	err := json.NewDecoder(request.Body).Decode(_value.Interface())

	if err != nil {
		log.Println(err)
	}

	val := _value.Elem()

	return &val
}
