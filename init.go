package marshrudka

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

func (d *Drive) Register(_interface interface{}, _struct ...interface{}) *Drive {

	if _struct == nil {
		_structValue := reflect.ValueOf(_interface)
		_structElemet := _structValue.Elem()
		d.services[_structValue.Type()] = _structValue
		d.services[_structElemet.Type()] = _structElemet
		return d
	}

	if len(_struct) != 1 {
		log.Fatalln("something went wrong")
	}

	if implement(_interface, _struct[0]) {
		_interfaceType := reflect.TypeOf(_interface)
		_structValue := reflect.ValueOf(_struct[0])
		d.services[_interfaceType.Elem()] = _structValue
		d.services[_structValue.Type()] = _structValue
	} else {
		log.Fatalln("something went wrong")
	}

	return d
}

func implement(_interface, _struct interface{}) bool {

	structType := reflect.TypeOf(_struct)
	{
		if structType.Kind() != reflect.Ptr {
			log.Fatalln("ffs 1")
		}
	}

	interfaceType := reflect.TypeOf(_interface)
	{
		if interfaceType.Kind() != reflect.Ptr {
			log.Fatalln("ffs 2")
		}
	}

	if interfaceType.Elem().Kind() == reflect.Struct {
		return structType.AssignableTo(interfaceType)
	}

	return structType.AssignableTo(interfaceType.Elem())
}

func (d *Drive) Use(handlers ...interface{}) *Drive {
	d.handlers = handlers
	return d
}

func (d *Drive) checkHas() actions {
	var driverActions actions

	if len(d.handlers) > 0 {
		driverActions = parseFunc("", []string{}, d.handlers...).actions
	}

	return driverActions
}

func (d *Drive) ANY(path string, handlers ...interface{}) *router {

	_router := parseFunc(path, []string{"ANY"}, handlers...)

	_router.actions = append(d.checkHas(), _router.actions...)

	d.routers.Add(_router)
	return _router
}

func (d *Drive) GET(path string, handlers ...interface{}) *router {
	_router := parseFunc(path, []string{http.MethodGet}, handlers...)
	_router.actions = append(d.checkHas(), _router.actions...)
	d.routers.Add(_router)
	return _router
}

func (d *Drive) POST(path string, handlers ...interface{}) *router {
	_router := parseFunc(path, []string{http.MethodPost}, handlers...)
	_router.actions = append(d.checkHas(), _router.actions...)
	d.routers.Add(_router)
	return _router
}

func (d *Drive) PUT(path string, handlers ...interface{}) *router {
	_router := parseFunc(path, []string{http.MethodPut}, handlers...)
	_router.actions = append(d.checkHas(), _router.actions...)
	d.routers.Add(_router)
	return _router
}

func (d *Drive) PATCH(path string, handlers ...interface{}) *router {
	_router := parseFunc(path, []string{http.MethodPatch}, handlers...)
	_router.actions = append(d.checkHas(), _router.actions...)
	d.routers.Add(_router)
	return _router
}

func (d *Drive) DELETE(path string, handlers ...interface{}) *router {
	_router := parseFunc(path, []string{http.MethodDelete}, handlers...)
	_router.actions = append(d.checkHas(), _router.actions...)
	d.routers.Add(_router)
	return _router
}

func (d *Drive) MATCH(path string, methods []string, handlers ...interface{}) *router {

	for i, method := range methods {
		methods[i] = strings.ToUpper(method)
	}

	_router := parseFunc(path, methods, handlers...)
	_router.actions = append(d.checkHas(), _router.actions...)
	d.routers.Add(_router)
	return _router
}

func (g *group) MATCH(path string, methods []string, handlers ...interface{}) {

	if strings.HasPrefix(path, "/") {
		path = g.Path + path[1:]
	} else {
		path = g.Path + path
	}

	g.Drive.MATCH(path, methods, append(g.actions, handlers...)...)
}

func (g *group) ANY(path string, handlers ...interface{}) *router {

	if strings.HasPrefix(path, "/") {
		path = g.Path + path[1:]
	} else {
		path = g.Path + path
	}

	return g.Drive.ANY(path, append(g.actions, handlers...)...)
}

func (g *group) GET(path string, handlers ...interface{}) *router {

	if strings.HasPrefix(path, "/") {
		path = g.Path + path[1:]
	} else {
		path = g.Path + path
	}

	return g.Drive.GET(path, append(g.actions, handlers...)...)
}

func (g *group) POST(path string, handlers ...interface{}) *router {

	if strings.HasPrefix(path, "/") {
		path = g.Path + path[1:]
	} else {
		path = g.Path + path
	}

	return g.Drive.POST(path, append(g.actions, handlers...)...)
}

func (g *group) PUT(path string, handlers ...interface{}) *router {

	if strings.HasPrefix(path, "/") {
		path = g.Path + path[1:]
	} else {
		path = g.Path + path
	}

	return g.Drive.PUT(path, append(g.actions, handlers...)...)
}

func (g *group) PATCH(path string, handlers ...interface{}) *router {

	if strings.HasPrefix(path, "/") {
		path = g.Path + path[1:]
	} else {
		path = g.Path + path
	}

	return g.Drive.PATCH(path, append(g.actions, handlers...)...)
}

func (g *group) DELETE(path string, handlers ...interface{}) *router {

	if strings.HasPrefix(path, "/") {
		path = g.Path + path[1:]
	} else {
		path = g.Path + path
	}

	return g.Drive.DELETE(path, append(g.actions, handlers...)...)
}

func parseFunc(path string, methods []string, handlers ...interface{}) *router {
	_actions := actions{}

	for index, handler := range handlers {
		_func := reflect.ValueOf(handler)

		if _func.Kind() != reflect.Func {
			log.Fatalln("type not supported", _func.Kind())
		}

		_action := &action{}

		_funcType := _func.Type()

		if (len(handlers)-1) == index && _funcType.NumOut() > 1 {
			log.Fatalln("error end function cannot return data greater than 1:", _funcType.NumOut())
		}

		for i := 0; i < _funcType.NumIn(); i++ {

			paramType := _funcType.In(i)

			if isPrimitive(paramType.Kind()) && index == 0 {
				log.Fatalln("error the first function cannot accept primitive data types:", paramType.Kind())
			}

			_action.Params = append(_action.Params, paramType)
		}

		_action.Ret = _funcType.NumOut() == 0

		_action.Call = _func.Call

		_actions.Add(_action)
	}

	uri, params := parseUrl(path)

	_map := map[string]bool{}

	for _, method := range methods {
		_map[method] = true
	}

	return &router{
		path:    path,
		params:  params,
		method:  _map,
		uri:     uri,
		actions: _actions,
	}

}

func parseUrl(path string) (*regexp.Regexp, []string) {

	if path == "" {
		path = "/"
	}

	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}

	if strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}

	if strings.Index(path, ":") != -1 {

		reg, names := pattern(path+"/", []string{})

		path = strings.TrimSuffix(reg, "/")

		rep := strings.Replace(path, `(\w+?)`, "", -1)
		rep = strings.Replace(rep, `//`, "/", -1)

		path = fmt.Sprintf("^((/?%s/?)|(/?%s?))$", path, rep)

		return regexp.MustCompile(path), names
	}

	return regexp.MustCompile(fmt.Sprintf("^/?%s/?$", path)), []string{}
}

func pattern(path string, names []string) (string, []string) {

	index := strings.Index(path, ":")

	if index == -1 {
		return path, names
	}

	hasOne := path[index:]
	index2 := strings.Index(hasOne, "/")

	if index2 == -1 {
		return path, names
	}

	hasTwo := hasOne[:index2]
	names = append(names, hasTwo[1:])
	path = strings.Replace(path, hasTwo, `(\w+?)`, 1)

	return pattern(path, names)
}

func isPrimitive(kind reflect.Kind) bool {

	switch kind {
	case reflect.Bool,
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
		reflect.String:
		return true
	default:
		return false
	}
}

func newRegexp(url string) *regexp.Regexp {

	l := len(url)
	var _Regexp *regexp.Regexp

	if l == 1 {
		if strings.HasPrefix(url, "/") || url == "" {
			url = "/"
			_regexp, _ := regexp.Compile(`^/$`)
			_Regexp = _regexp

		} else {
			tmpUrl := fmt.Sprintf(`^(/?)%s(/?)$`, url)
			_regexp, _ := regexp.Compile(tmpUrl)
			_Regexp = _regexp
			url = fmt.Sprintf(`/%s/`, url)
		}
	}

	if l != 1 {

		tmpUrl := strings.TrimPrefix(url, "/")
		tmpUrl = strings.TrimSuffix(tmpUrl, "/")
		tmpUrl = fmt.Sprintf(`^(/?)%s(/?)$`, tmpUrl)
		_regexp, _ := regexp.Compile(tmpUrl)
		_Regexp = _regexp

		if !strings.HasPrefix(url, "/") {
			url = fmt.Sprintf(`/%s`, url)
		}

		if !strings.HasSuffix(url, "/") {
			url = fmt.Sprintf(`%s/`, url)
		}

	}

	return _Regexp
}
