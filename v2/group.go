package v2

import (
	"fmt"
	"strings"
)

type group struct {
	Path              string
	HandlersInterface handlersInterface
	handlers          handlers
	Server            *server
	child             *group
}

func (g *group) Use(handlers ...interface{}) {
	g.HandlersInterface.AddRange(handlers)
}

func (g *group) Group(path string, handlers ...interface{}) IRouter {

	var _handlersInterface = handlersInterface{}

	return &group{
		Path:              fmt.Sprintf("%s/%s", strings.TrimSuffix(g.Path, "/"), strings.TrimPrefix(path, "/")),
		HandlersInterface: *_handlersInterface.AddRange(g.HandlersInterface).AddRange(handlers),
		Server:            g.Server,
	}
}

func (g *group) GET(path string, handlers ...interface{}) {
	g.Server.GET(path, handlers...)
}

func (g group) POST(path string, handlers ...interface{}) {
	g.Server.POST(path, handlers...)
}

func (g group) PUT(path string, handlers ...interface{}) {
	g.Server.PUT(path, handlers...)
}

func (g group) PATCH(path string, handlers ...interface{}) {
	g.Server.PATCH(path, handlers...)
}

func (g group) DELETE(path string, handlers ...interface{}) {
	g.Server.DELETE(path, handlers...)
}

func (g group) ANY(path string, handlers ...interface{}) {
	g.Server.ANY(path, handlers...)
}

func (g group) MATCH(path string, methods []string, handlers ...interface{}) {
	g.Server.MATCH(path, methods, handlers...)
}
