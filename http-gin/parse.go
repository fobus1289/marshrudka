package http_gin

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
)

var (
	httpResponse     = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	httpRequest      = reflect.TypeOf(&http.Request{})
	_throw           = reflect.TypeOf(&Throw{})
	_response        = reflect.TypeOf(&Response{})
	_context         = reflect.TypeOf(&gin.Context{})
	methodNotAllowed = []byte("method not allowed")
	expectsJSON      = "expects to receive a JSON object"
)

func (s *Serv) Register(_interface interface{}, _struct ...interface{}) *Serv {

	if _struct == nil {
		_structValue := reflect.ValueOf(_interface)
		_structElemet := _structValue.Elem()
		s.services[_structValue.Type()] = _structValue
		s.services[_structElemet.Type()] = _structElemet
		return s
	}

	if len(_struct) != 1 {
		log.Fatalln("something went wrong")
	}

	if implement(_interface, _struct[0]) {
		_interfaceType := reflect.TypeOf(_interface)
		_structValue := reflect.ValueOf(_struct[0])
		s.services[_interfaceType.Elem()] = _structValue
		s.services[_structValue.Type()] = _structValue
	} else {
		log.Fatalln("something went wrong")
	}

	return s
}

func (s *Serv) parseFunc(handler ...interface{}) gin.HandlersChain {

	var handlersChain gin.HandlersChain

	l := len(handler) - 1

	var outValue = map[reflect.Type]reflect.Value{}

	for j, h := range handler {

		handlerValue := reflect.ValueOf(h)

		if handlerValue.Kind() != reflect.Func {
			log.Fatalln("this type is not supported", handlerValue.Kind())
		}

		handlerType := handlerValue.Type()

		var params []reflect.Type

		for i := 0; i < handlerType.NumIn(); i++ {
			inType := handlerType.In(i)
			params = append(params, inType)
		}

		handlersChain = append(handlersChain, s.getFunc(handlerValue, outValue, params, j, l))
	}

	return handlersChain
}

func (s *Serv) getFunc(handlerValue reflect.Value, outValues map[reflect.Type]reflect.Value, params []reflect.Type, index, l int) func(context *gin.Context) {
	i := index

	return func(context *gin.Context) {

		var values []reflect.Value

		for _, param := range params {

			outParam := outValues[param]

			if outParam.Kind() != reflect.Invalid {
				values = append(values, outParam)
				continue
			}

			if reflect.DeepEqual(httpResponse, param) {
				values = append(values, reflect.ValueOf(context.Writer))
				continue
			}

			if reflect.DeepEqual(httpRequest, param) {
				values = append(values, reflect.ValueOf(context.Request))
				continue
			}

			if reflect.DeepEqual(_context, param) {
				values = append(values, reflect.ValueOf(context))
				continue
			}

			service := s.services[param]

			if service.Kind() != reflect.Invalid {
				values = append(values, service)
				continue
			}

			value := reflect.New(param)

			if err := context.BindJSON(value.Interface()); err != nil {
				context.AbortWithStatusJSON(http.StatusBadRequest, expectsJSON)
				return
			}

			value = value.Elem()

			outValues[value.Type()] = value

			values = append(values, value)
		}

		ret(handlerValue.Call(values), outValues, context, l == i)
	}

}

func ret(retValues []reflect.Value, outValues map[reflect.Type]reflect.Value, c *gin.Context, last bool) {

	if len(retValues) < 1 {
		return
	}

	if last {

		value := reflect.ValueOf(retValues[0].Interface())

		if isThrow(value, c) {
			return
		}

		if isResponse(value, c) {
			return
		}

		c.JSON(http.StatusOK, value.Interface())

		return
	}

	for _, value := range retValues {
		val := reflect.ValueOf(value.Interface())

		if isThrow(val, c) {
			return
		}

		outValues[val.Type()] = val
	}

}

func isThrow(val reflect.Value, c *gin.Context) bool {
	if reflect.DeepEqual(val.Type(), _throw) {
		outThrow := val.Interface().(*Throw)
		c.AbortWithStatusJSON(outThrow.StatusCode, outThrow.Data)
		return true
	}
	return false
}

func isResponse(val reflect.Value, c *gin.Context) bool {
	if reflect.DeepEqual(val.Type(), _response) {
		outResponse := val.Interface().(*Response)
		c.JSON(outResponse.StatusCode, outResponse.Data)
		return true
	}
	return false
}
