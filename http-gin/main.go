package http_gin

import (
	"github.com/gin-gonic/gin"
	"log"
	"reflect"
)

type Serv struct {
	g        *gin.Engine
	services map[reflect.Type]reflect.Value
}

func NewDrive() *Serv {
	return &Serv{g: gin.New(), services: map[reflect.Type]reflect.Value{}}
}

func (s *Serv) Use(handler ...interface{}) {
	s.g.Use(s.parseFunc(handler...)...)
}

func (s *Serv) Run(addr string) {
	log.Fatalln(s.g.Run(addr))
}

func (s *Serv) Group(name string, handler ...interface{}) *Group {
	return &Group{
		routerGroup: s.g.Group(name, s.parseFunc(handler...)...),
		serv:        s,
	}
}

type Group struct {
	routerGroup *gin.RouterGroup
	serv        *Serv
	group       *Group
}

func (g *Group) Group(name string, handler ...interface{}) *Group {

	childGroup := &Group{
		routerGroup: g.routerGroup.Group(name, g.serv.parseFunc(handler...)...),
		serv:        g.serv,
	}

	g.group = childGroup

	return childGroup
}

func (g *Group) GET(name string, handler ...interface{}) *Group {
	g.routerGroup.GET(name, g.serv.parseFunc(handler...)...)
	return g
}

func (g *Group) POST(name string, handler ...interface{}) *Group {
	g.routerGroup.POST(name, g.serv.parseFunc(handler...)...)
	return g
}

func (g *Group) PUT(name string, handler ...interface{}) *Group {
	g.routerGroup.PUT(name, g.serv.parseFunc(handler...)...)
	return g
}

func (g *Group) PATCH(name string, handler ...interface{}) *Group {
	g.routerGroup.PATCH(name, g.serv.parseFunc(handler...)...)
	return g
}

func (g *Group) ANY(name string, handler ...interface{}) *Group {
	g.routerGroup.Any(name, g.serv.parseFunc(handler...)...)
	return g
}

func (g *Group) DELETE(name string, handler ...interface{}) *Group {
	g.routerGroup.DELETE(name, g.serv.parseFunc(handler...)...)
	return g
}
