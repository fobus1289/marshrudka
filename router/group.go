package router

import (
	"fmt"
	"strings"
)

type group struct {
	path    string
	handler []interface{}
	group   *group
	drive   *Drive
}

func (d *Drive) Group(name string, handler ...interface{}) *group {
	return &group{
		path:    name,
		handler: handler,
		drive:   d,
	}
}

func (g *group) Group(name string, handler ...interface{}) *group {

	newGroup := &group{
		path:    fmt.Sprintf("%s/%s", strings.TrimSuffix(g.path, "/"), strings.TrimPrefix(name, "/")),
		handler: append(g.handler, handler...),
		drive:   g.drive,
	}

	g.group = newGroup

	return newGroup
}

func (g *group) GET(name string, handler ...interface{}) {
	g.drive.GET(fmt.Sprintf("%s/%s", strings.TrimSuffix(g.path, "/"), strings.TrimPrefix(name, "/")), append(g.handler, handler...)...)
}

func (g *group) POST(name string, handler ...interface{}) {
	g.drive.POST(fmt.Sprintf("%s/%s", strings.TrimSuffix(g.path, "/"), strings.TrimPrefix(name, "/")), append(g.handler, handler...)...)
}

func (g *group) PUT(name string, handler ...interface{}) {
	g.drive.PUT(fmt.Sprintf("%s/%s", strings.TrimSuffix(g.path, "/"), strings.TrimPrefix(name, "/")), append(g.handler, handler...)...)
}

func (g *group) PATCH(name string, handler ...interface{}) {
	g.drive.PATCH(fmt.Sprintf("%s/%s", strings.TrimSuffix(g.path, "/"), strings.TrimPrefix(name, "/")), append(g.handler, handler...)...)
}

func (g *group) DELETE(name string, handler ...interface{}) {
	g.drive.DELETE(fmt.Sprintf("%s/%s", strings.TrimSuffix(g.path, "/"), strings.TrimPrefix(name, "/")), append(g.handler, handler...)...)
}

func (g *group) ANY(name string, handler ...interface{}) {
	g.drive.ANY(fmt.Sprintf("%s/%s", strings.TrimSuffix(g.path, "/"), strings.TrimPrefix(name, "/")), append(g.handler, handler...)...)
}

func (g *group) MATCH(name string, methods []string, handler ...interface{}) {
	g.drive.MATCH(fmt.Sprintf("%s/%s", strings.TrimSuffix(g.path, "/"), strings.TrimPrefix(name, "/")), methods, append(g.handler, handler...)...)
}
