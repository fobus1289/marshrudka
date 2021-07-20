package router

import (
	"net/http"
	"strings"
)

func (d *Drive) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ref = reflectMap{}

	if !d.handlers.each(w, r, ref) {
		return
	}

	d.routes.each(w, r, ref)
}

func (d *Drive) GET(name string, handler ...interface{}) {

	_handlers := d.parseFunc(false, handler...)

	match, names := parsePath(name)

	route := &route{
		path:       name,
		match:      match,
		paramNames: names,
		methods: map[string]bool{
			http.MethodGet: true,
		},
		handlers: _handlers,
	}
	_handlers.setRoutes(route)
	d.routes = append(d.routes, route)
}

func (d *Drive) POST(name string, handler ...interface{}) {

	_handlers := d.parseFunc(false, handler...)

	match, names := parsePath(name)

	route := &route{
		path:       name,
		match:      match,
		paramNames: names,
		methods: map[string]bool{
			http.MethodPost: true,
		},
		handlers: _handlers,
	}

	_handlers.setRoutes(route)

	d.routes = append(d.routes, route)
}

func (d *Drive) PUT(name string, handler ...interface{}) {

	_handlers := d.parseFunc(false, handler...)

	match, names := parsePath(name)

	route := &route{
		path:       name,
		match:      match,
		paramNames: names,
		methods: map[string]bool{
			http.MethodPut: true,
		},
		handlers: _handlers,
	}
	_handlers.setRoutes(route)
	d.routes = append(d.routes, route)

}

func (d *Drive) PATCH(name string, handler ...interface{}) {

	_handlers := d.parseFunc(false, handler...)

	match, names := parsePath(name)
	route := &route{
		path:       name,
		match:      match,
		paramNames: names,
		methods: map[string]bool{
			http.MethodPatch: true,
		},
		handlers: _handlers,
	}

	_handlers.setRoutes(route)

	d.routes = append(d.routes, route)
}

func (d *Drive) DELETE(name string, handler ...interface{}) {

	_handlers := d.parseFunc(false, handler...)

	match, names := parsePath(name)

	route := &route{
		path:       name,
		match:      match,
		paramNames: names,
		methods: map[string]bool{
			http.MethodDelete: true,
		},
		handlers: _handlers,
	}
	_handlers.setRoutes(route)
	d.routes = append(d.routes, route)

}

func (d *Drive) ANY(name string, handler ...interface{}) {

	_handlers := d.parseFunc(false, handler...)

	match, names := parsePath(name)
	route := &route{
		path:       name,
		match:      match,
		paramNames: names,
		methods: map[string]bool{
			"ANY": true,
		},
		handlers: _handlers,
	}
	_handlers.setRoutes(route)

	d.routes = append(d.routes, route)
}

func (d *Drive) MATCH(name string, methods []string, handler ...interface{}) {

	_handlers := d.parseFunc(false, handler...)

	match, names := parsePath(name)

	var methodsMap = map[string]bool{}

	for _, method := range methods {
		methodsMap[strings.ToUpper(method)] = true
	}
	route := &route{
		path:       name,
		match:      match,
		paramNames: names,
		methods:    methodsMap,
		handlers:   _handlers,
	}
	_handlers.setRoutes(route)
	d.routes = append(d.routes, route)
}
