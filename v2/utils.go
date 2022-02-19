package v2

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

var (
	httpRes     = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	httpReq     = reflect.TypeOf(&http.Request{})
	request     = reflect.TypeOf(&Request{})
	expectsJSON = []byte("expects to receive a JSON object")
)

type (
	paramsMap         map[reflect.Type]reflect.Value
	handlersInterface []interface{}
)

func (hsi *handlersInterface) AddRange(handlers []interface{}) *handlersInterface {
	if hsi == nil {
		*hsi = handlersInterface{}
	}
	*hsi = append(*hsi, handlers...)

	return hsi
}

func trimPrefixAndSuffix(s string) string {
	s = strings.TrimPrefix(s, "/")
	s = strings.TrimSuffix(s, "/")

	if len(s) < 1 {
		s = "/"
	}

	return s
}

func leftAndRightContact(a, b string) string {
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(a, "/"), strings.TrimPrefix(b, "/"))
}

func createRequestRegular(regular string) *regexp.Regexp {

	var compile, err = regexp.Compile(regular)

	if err != nil {
		panic(err)
	}

	return compile
}

func getRegular(urlPath string) string {
	urlPath = strings.TrimPrefix(urlPath, "/")
	urlPath = strings.TrimSuffix(urlPath, "/")
	var regular = regexp.MustCompile(`(:[a-zA-Z]+)`)
	urlPath = regular.ReplaceAllString(urlPath, `([0-9a-zA-Z]+)`)
	return fmt.Sprintf("^(/?%s/?)$", urlPath)
}

func getPattern(urlPath string) []string {
	var regular = regexp.MustCompile(`(:[a-zA-Z]+)`)

	if result := regular.FindAllString(urlPath, -1); len(result) > 0 {
		var str = strings.Replace(strings.Join(result, " "), ":", "", -1)
		return strings.Split(str, " ")
	}

	return []string{}
}

func parseFunc(s *server, parent bool, actions []interface{}) handlers {

	if len(actions) < 1 && !parent {
		panic("Must be move 1 handler")
	}

	var _handlers handlers

	for a, action := range actions {

		funcValue := reflect.ValueOf(action)

		if funcValue.Kind() != reflect.Func {
			log.Fatalln("dont supported this type:", funcValue.Kind())
		}

		funcType := funcValue.Type()

		var params []reflect.Type

		for i := 0; i < funcType.NumIn(); i++ {
			in := funcType.In(i)
			params = append(params, in)
		}

		var handler = &handler{
			Server: s,
			Last:   len(actions)-1 == a && !parent,
			Params: params,
			Call:   funcValue.Call,
		}

		_handlers = append(_handlers, handler)
	}

	return _handlers
}
