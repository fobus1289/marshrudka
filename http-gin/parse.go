package http_gin

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
	"strconv"
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

func (s *Serv) parseFunc(child bool, handler ...interface{}, ) gin.HandlersChain {

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

		handlersChain = append(handlersChain, s.getFunc(handlerValue, outValue, params, j, l, child))
	}

	return handlersChain
}

func (s *Serv) getFunc(handlerValue reflect.Value, outValues map[reflect.Type]reflect.Value, params []reflect.Type, index, l int, isChild bool) func(context *gin.Context) {
	i := index

	return func(context *gin.Context) {

		if val, ok := context.Get("outValuesReflect"); ok {
			if outValuesReflect, ok := val.(map[reflect.Type]reflect.Value); ok && len(outValuesReflect) > 0 {
				for r, v := range outValuesReflect {
					outValues[r] = v
				}
			}
		}

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

		ret(handlerValue.Call(values), outValues, context, l == i, isChild)
	}

}

func ret(retValues []reflect.Value, outValues map[reflect.Type]reflect.Value, c *gin.Context, last, isChild bool) {

	if len(retValues) < 1 {
		return
	}

	if isChild && last {

		retVal := retValues[0].Interface()

		if retVal == nil {
			return
		}

		value := reflect.ValueOf(retVal)

		if isThrow(value, c) || isResponse(value, c) {
			return
		}

		c.Data(200, "text/plain;utf-8", getPrimitiveResult(value))

		return
	}

	for _, value := range retValues {
		retVal := value.Interface()

		if retVal == nil {
			continue
		}

		val := reflect.ValueOf(retVal)

		if isThrow(val, c) {
			return
		}

		outValues[val.Type()] = val
	}

	c.Set("outValuesReflect", outValues)
}

func isThrow(val reflect.Value, c *gin.Context) bool {

	switch t := val.Interface().(type) {
	case Throw:
		c.AbortWithStatusJSON(t.StatusCode, t.Data)
		return true
	case *Throw:
		c.AbortWithStatusJSON(t.StatusCode, t.Data)
		return true
	}

	return false
}

func isResponse(val reflect.Value, c *gin.Context) bool {
	switch t := val.Interface().(type) {
	case Response:
		c.AbortWithStatusJSON(t.StatusCode, t.Data)
		return true
	case *Response:
		c.AbortWithStatusJSON(t.StatusCode, t.Data)
		return true
	}
	return false
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
