package router

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/fobus1289/marshrudka/request"
	"github.com/fobus1289/marshrudka/validator"
)

var (
	requestType      = reflect.TypeOf((*request.IRequest)(nil)).Elem()
	deserializeType  = reflect.TypeOf((*request.IDeserialize)(nil)).Elem()
	formFileType     = reflect.TypeOf((*request.IFormFile)(nil)).Elem()
	paramType        = reflect.TypeOf((*request.IParam)(nil)).Elem()
	headerType       = reflect.TypeOf((*request.IHeader)(nil)).Elem()
	queryParamType   = reflect.TypeOf((*request.IQueryParam)(nil)).Elem()
	httpResponseType = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	requestData      = reflect.TypeOf((*IRequestData)(nil)).Elem()
	httpRequestType  = reflect.TypeOf(&http.Request{})
	validatorType    = reflect.TypeOf((*validator.IValidator)(nil)).Elem()
	jwtUserType      = reflect.TypeOf((*request.IJwtUser)(nil)).Elem()
	emptyValue       = reflect.Value{}
)

func functionPareser(fn any, s *server) struct {
	Call       func([]reflect.Value) []reflect.Value
	ParamFuncs paramFuncs
	OutTypes   []reflect.Type
} {

	value := reflect.ValueOf(fn)

	if value.Kind() != reflect.Func {
		panic("error type not Func")
	}

	defer func() {
		if err := recover(); err != nil {
			log.Fatalln(err)
		}
	}()

	var funs paramFuncs

	valueType := value.Type()

	for i := 0; i < valueType.NumIn(); i++ {
		param := valueType.In(i)

		p := param

		for p.Kind() == reflect.Ptr {
			p = p.Elem()
		}

		if param.PkgPath() == "" && param.Kind() == reflect.Interface {
			panic("param.PkgPath() ==''")
		}

		switch param {

		case requestType, deserializeType, formFileType, paramType, queryParamType, headerType:
			funs.Add(requestToReflectValue())
			continue
		case httpResponseType:
			funs.Add(httpResponseToReflectValue())
			continue
		case httpRequestType:
			funs.Add(httpRequestToReflectValue())
			continue
		default:
			service := s.Services[param]

			if service != nil {
				funs = append(funs, instanceServiceToReflectValue(param, service))
				continue
			}

			var has bool

			for key, value := range s.Services {
				if param.Kind() != reflect.Struct && key.AssignableTo(param) {
					funs = append(funs, instanceServiceToReflectValue(param, value))
					has = true
					break
				}
			}

			if has {
				continue
			}
		}

		if val := deserialize(param, s); val != nil {
			funs = append(funs, val)
			continue
		}

	}

	var outTypes []reflect.Type

	for i := 0; i < valueType.NumOut(); i++ {

		out := valueType.Out(i)

		if out.Kind() == reflect.Interface {
			if out.PkgPath() == "" || out.NumMethod() == 0 {
				log.Fatalln("any or interface not supported")
			}
		}

		outTypes = append(outTypes, out)
	}

	return struct {
		Call       func([]reflect.Value) []reflect.Value
		ParamFuncs paramFuncs
		OutTypes   []reflect.Type
	}{
		Call:       value.Call,
		ParamFuncs: funs,
		OutTypes:   outTypes,
	}

}

func deserialize(t reflect.Type, s *server) paramFunc {

	param := t

	if s.Jwt != nil && param.Implements(jwtUserType) {
		jwt := s.Jwt

		paramType := param
		{
			if paramType.Kind() == reflect.Interface {
				panic("paramType.Kind() == reflect.Interface")
			}

			for paramType.Kind() == reflect.Ptr {
				paramType = paramType.Elem()
			}

		}

		return func(hp *handlerParam) (reflect.Value, *RuntimeError) {
			session := hp.SessionData[paramType]

			if session.Kind() == reflect.Invalid {

				authorization := hp.Request.Header.Get("Authorization")
				authorization = strings.TrimPrefix(authorization, "Bearer ")

				newValue := reflect.New(paramType)

				jwtUser, ok := newValue.Interface().(request.IJwtUser)

				if !ok {
					err := errors.New("invalid jwt user type")
					return emptyValue, &RuntimeError{
						Error:       err,
						Status:      http.StatusUnauthorized,
						ContentType: "application/json; charset=utf-8",
						Data:        []byte(err.Error()),
					}
				}

				if err := jwt.DecodeWithExpired(authorization, jwtUser); err != nil {
					return emptyValue, &RuntimeError{
						Error:       err,
						Status:      http.StatusUnauthorized,
						ContentType: "application/json; charset=utf-8",
						Data:        []byte(err.Error()),
					}
				}

				session = newValue

				hp.SessionData[paramType] = session
			}

			return session, nil
		}
	}

	if !param.Implements(validatorType) {
		return func(hp *handlerParam) (reflect.Value, *RuntimeError) {
			return hp.SessionData[param], nil
		}
	}

	deserializeErrorFunc := s.DeserializeErrorFunc

	return func(hp *handlerParam) (reflect.Value, *RuntimeError) {

		value := hp.SessionData[requestData]
		{
			if value.Kind() == reflect.Invalid {

				newValue := reflect.New(param)

				if err := json.NewDecoder(hp.Request.Body).Decode(newValue.Interface()); err != nil {
					return emptyValue, deserializeErrorFunc(err)
				}

				value = newValue.Elem()

				if valid, ok := value.Interface().(validator.IValidator); ok {
					mapResult := valid.Validate(hp.Request.Method)
					if len(mapResult) != 0 {
						data, _ := json.Marshal(mapResult)
						return emptyValue, &RuntimeError{
							Status:      http.StatusBadRequest,
							ContentType: "application/json; charset=utf-8",
							Data:        data,
						}
					}
				}

				hp.SessionData[value.Type()] = value
			}
		}

		return value, nil
	}
}

func requestToReflectValue() paramFunc {
	return func(hp *handlerParam) (reflect.Value, *RuntimeError) {

		value := hp.SessionData[requestType]
		{
			if value.Kind() == reflect.Invalid {
				value = reflect.ValueOf(request.NewRequest(hp.Request, hp.HttpParamMap))
				hp.SessionData[requestType] = value
			}
		}

		return value, nil
	}
}

func httpRequestToReflectValue() paramFunc {
	return func(hp *handlerParam) (reflect.Value, *RuntimeError) {
		return reflect.ValueOf(hp.Request), nil
	}
}

func httpResponseToReflectValue() paramFunc {
	return func(hp *handlerParam) (reflect.Value, *RuntimeError) {
		return reflect.ValueOf(hp.Request), nil
	}
}

func instanceServiceToReflectValue(t reflect.Type, service paramFunc) paramFunc {

	return func(hp *handlerParam) (reflect.Value, *RuntimeError) {
		value, err := service(hp)

		if err != nil {
			return emptyValue, err
		}

		return value, nil
	}
}
