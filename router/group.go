package router

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

func (g *group) Group(_path string, handlers ...interface{}) IGroup {
	return &group{
		Path:              pathJoin(g.Path, _path),
		HandlersInterface: interfaceJoin(g.HandlersInterface, handlers),
		Server:            g.Server,
	}
}

func (g *group) GET(_path string, handlers ...interface{}) IMatch {
	return g.Server.GET(pathJoin(g.Path, _path), interfaceJoin(g.HandlersInterface, handlers)...)
}

func (g *group) POST(_path string, handlers ...interface{}) IMatch {
	return g.Server.POST(pathJoin(g.Path, _path), interfaceJoin(g.HandlersInterface, handlers)...)
}

func (g *group) PUT(_path string, handlers ...interface{}) IMatch {
	return g.Server.PUT(pathJoin(g.Path, _path), interfaceJoin(g.HandlersInterface, handlers)...)
}

func (g *group) PATCH(_path string, handlers ...interface{}) IMatch {
	return g.Server.PATCH(pathJoin(g.Path, _path), interfaceJoin(g.HandlersInterface, handlers)...)
}

func (g *group) DELETE(_path string, handlers ...interface{}) IMatch {
	return g.Server.DELETE(pathJoin(g.Path, _path), interfaceJoin(g.HandlersInterface, handlers)...)
}

func (g *group) ANY(_path string, handlers ...interface{}) IMatch {
	return g.Server.ANY(pathJoin(g.Path, _path), interfaceJoin(g.HandlersInterface, handlers)...)
}

func (g *group) MATCH(_path string, methods []string, handlers ...interface{}) IMatch {
	return g.Server.MATCH(pathJoin(g.Path, _path), methods, interfaceJoin(g.HandlersInterface, handlers)...)
}

func (g *group) FileServer(method, path, dir string, beforeHandlers ...interface{}) {
	g.Server.FileServer(method, path, dir, interfaceJoin(g.HandlersInterface, beforeHandlers)...)
}
