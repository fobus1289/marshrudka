package router

import (
	"context"
	"errors"
	"fmt"
	request2 "github.com/fobus1289/marshrudka/router/request"
	"log"
	"net/http"
	"path"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

var (
	httpRes          = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	iFormFile        = reflect.TypeOf((*request2.IFormFile)(nil)).Elem()
	iModel           = reflect.TypeOf(request2.IModel(nil))
	iService         = reflect.TypeOf((*IService)(nil)).Elem()
	httpReq          = reflect.TypeOf(&http.Request{})
	request          = reflect.TypeOf((*request2.IRequest)(nil)).Elem()
	iParam           = reflect.TypeOf((*request2.IParam)(nil)).Elem()
	iQueryParam      = reflect.TypeOf((*request2.IQueryParam)(nil)).Elem()
	whatWentWrong    = []byte("what went wrong :(")
	whatWentWrongErr = errors.New(string(whatWentWrong))
	emptyBody        = []byte("empty body :(")
	methodNotAllowed = []byte("method not allowed")
)

type (
	reflectMap        map[reflect.Type]reflect.Value
	handlersInterface []interface{}
)

type options struct {
	server   *server
	parent   bool
	handlers []interface{}
	methods  map[string]bool
	path     string
}

func (hsi handlersInterface) AddRange(handlers []interface{}) handlersInterface {
	return append(hsi, handlers...)
}

func interfaceJoin(a, b []interface{}) handlersInterface {
	return append(a, b...)
}

func pathJoin(s ...string) string {
	return path.Join(s...)
}

func initRouter(o *options) *router {
	var _handlers = parseFunc(o.server, o.parent, o.handlers)
	var route = &router{
		Path:       o.path,
		Match:      createRequestRegular(getRegular(o.path)),
		WhereMatch: nil,
		Params:     getPattern(o.path),
		Methods:    o.methods,
		Handlers:   _handlers,
	}
	route.HandlerFunc = route.ServeHTTP
	_handlers.SetRouter(route)
	o.server.routers = append(o.server.routers, route)
	return route
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

func rebuildRouters(s *server) {
	var routes routers

	for _, r := range s.routers {
		if !strings.Contains(r.Path, ":") && !strings.Contains(r.Path, "*") {
			routes = append(routes, r)
		}
	}

	for _, r := range s.routers {
		if strings.Contains(r.Path, ":") && !strings.Contains(r.Path, "*") {
			routes = append(routes, r)
		}
	}

	var routesStar = routers{}

	for _, r := range s.routers {
		if strings.Contains(r.Path, "*") {
			routesStar = append(routesStar, r)
		}
	}

	sort.Slice(routesStar, func(i, j int) bool {
		return len(routesStar[i].Path) > len(routesStar[j].Path)
	})

	s.routers = append(routes, routesStar...)
}

func showUrl(s *server, addr string) {
	showAddr := strings.TrimSuffix(addr, "/")

	if strings.HasPrefix(addr, ":") {
		showAddr = "http://localhost" + showAddr
	}

	showAddr += "/"

	bigLen := 0

	for _, r := range s.routers {
		if len(showAddr+r.Path) > bigLen {
			bigLen = len(showAddr + r.Path)
		}
	}

	for _, r := range s.routers {
		var methods []string

		for s, _ := range r.Methods {
			methods = append(methods, s)
		}

		path := showAddr + strings.TrimPrefix(r.Path, "/")

		pathLen := len(path) - 1

		if pathLen < bigLen {
			fmt.Println("path-> ", path, strings.Repeat(" ", bigLen-pathLen), " methods ", strings.Join(methods, ","))
		} else {
			fmt.Println("path-> ", path, " methods ", strings.Join(methods, ","))
		}
	}
}

func getRegular(urlPath string) string {
	urlPath = path.Clean(urlPath)
	{
		if urlPath == "." || urlPath == "/" {
			return "^(/)$"
		}
		urlPath = strings.TrimPrefix(urlPath, "/")
	}

	if index := strings.Index(urlPath, "*"); index != -1 {
		urlPath = urlPath[0 : index+1]
	}
	urlPath = strings.Replace(urlPath, "*", "(.*)?", -1)
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

	for i, action := range actions {
		switch a := action.(type) {
		case http.Handler:
			actions[i] = a.ServeHTTP
		}
	}

	if len(actions) < 1 && !parent {
		panic("Must be move 1 handler")
	}

	var (
		_handlers handlers
	)

	for a, action := range actions {

		funcValue := reflect.ValueOf(action)

		if funcValue.Kind() != reflect.Func {
			log.Fatalln("dont supported this type:", funcValue.Kind())
		}

		funcType := funcValue.Type()

		var (
			params     []reflect.Type
			parseParam = map[reflect.Type]func(http.ResponseWriter, *http.Request, *handler) reflect.Value{}
		)

		for i := 0; i < funcType.NumIn(); i++ {
			var inType = funcType.In(i)
			params = append(params, inType)
			if prm := getParseFunc(inType, s); prm != nil {
				parseParam[inType] = prm
			}
		}

		var handler = &handler{
			Server:     s,
			Params:     params,
			ParseParam: parseParam,
			Call:       funcValue.Call,
			Last:       len(actions)-1 == a && !parent,
		}

		_handlers = append(_handlers, handler)
	}

	return _handlers
}

func getParseFunc(key reflect.Type, s *server) func(http.ResponseWriter, *http.Request, *handler) reflect.Value {

	switch key {
	case httpRes:
		return func(w http.ResponseWriter, r *http.Request, h *handler) reflect.Value {
			return reflect.ValueOf(w)
		}
	case httpReq:
		return func(w http.ResponseWriter, r *http.Request, h *handler) reflect.Value {
			return reflect.ValueOf(r)
		}
	case request:
		return func(w http.ResponseWriter, r *http.Request, h *handler) reflect.Value {
			ctx := context.WithValue(r.Context(), "params", &request2.Params{
				Keys:  h.Router.Params,
				Match: h.Router.Match,
			})
			return reflect.ValueOf(request2.NewRequest(w, r.WithContext(ctx)))
		}
	case iFormFile:
		return func(w http.ResponseWriter, r *http.Request, h *handler) reflect.Value {
			return reflect.ValueOf(request2.NewFormFile(w, r))
		}
	}

	if value := s.GetByType(key); value.IsValid() {
		return func(w http.ResponseWriter, r *http.Request, h *handler) reflect.Value {
			return s.GetByType(key)
		}
	}

	if isBodyObject(key) {
		return func(w http.ResponseWriter, r *http.Request, h *handler) reflect.Value {
			return read(key, r)
		}
	}

	return nil
}

func isBodyObject(key reflect.Type) bool {

	if key.Kind() == reflect.Ptr {
		key = key.Elem()
	}

	switch key.Kind() {
	case
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
		reflect.Array,
		reflect.Chan,
		reflect.Func,
		reflect.String,
		reflect.UnsafePointer:
		return false
	}

	return true
}
