package router

import "path"

type IGroup interface {
	Use(handlers ...any) IGroup
	Group(actionPath string, handlers ...any) IGroup
	IMethod
}

type group struct {
	Path     string
	Parent   *group
	Server   *server
	Handlers []any
}

func (s *server) Group(actionPath string, handlers ...any) IGroup {
	return &group{
		Path:     actionPath,
		Server:   s,
		Handlers: handlers,
	}
}

func (g *group) Use(handlers ...any) IGroup {
	g.Handlers = append(g.Handlers, handlers...)
	return g
}

func (g *group) Group(actionPath string, handlers ...any) IGroup {
	return &group{
		Path:     path.Join(g.Path, actionPath),
		Parent:   g,
		Server:   g.Server,
		Handlers: append(g.Handlers, handlers...),
	}
}

func (g *group) GET(actionPath string, handlers ...any) IRouter {
	return g.Server.GET(path.Join(g.Path, actionPath), append(g.Handlers, handlers...)...)
}

func (g *group) POST(actionPath string, handlers ...any) IRouter {
	return g.Server.POST(path.Join(g.Path, actionPath), append(g.Handlers, handlers...)...)
}

func (g *group) PUT(actionPath string, handlers ...any) IRouter {
	return g.Server.PUT(path.Join(g.Path, actionPath), append(g.Handlers, handlers...)...)
}

func (g *group) PATCH(actionPath string, handlers ...any) IRouter {
	return g.Server.PATCH(path.Join(g.Path, actionPath), append(g.Handlers, handlers...)...)
}

func (g *group) DELETE(actionPath string, handlers ...any) IRouter {
	return g.Server.DELETE(path.Join(g.Path, actionPath), append(g.Handlers, handlers...)...)
}

func (g *group) MATCH(actionPath string, methods []string, handlers ...any) IRouter {
	return g.Server.MATCH(path.Join(g.Path, actionPath), methods, append(g.Handlers, handlers...)...)
}

func (g *group) ANY(actionPath string, handlers ...any) IRouter {
	return g.Server.ANY(path.Join(g.Path, actionPath), append(g.Handlers, handlers...)...)
}
