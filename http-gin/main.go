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
	g := gin.New()
	g.HandleMethodNotAllowed = true
	return &Serv{g: g, services: map[reflect.Type]reflect.Value{}}
}

func (s *Serv) Use(handler ...interface{}) {
	s.g.Use(s.parseFunc(false, handler...)...)
}

func (s *Serv) Run(addr string) {
	log.Fatalln(s.g.Run(addr))
}

func (s *Serv) Group(name string, handler ...interface{}) *Group {
	group := &Group{
		RouterGroup: s.g.Group(name, s.parseFunc(false, handler...)...),
		serv:        s,
	}

	return group
}

type Group struct {
	RouterGroup *gin.RouterGroup
	serv        *Serv
	group       *Group
}

func (g *Group) Group(name string, handler ...interface{}) *Group {

	childGroup := &Group{
		RouterGroup: g.RouterGroup.Group(name, g.serv.parseFunc(false, handler...)...),
		serv:        g.serv,
	}

	g.group = childGroup

	return childGroup
}

func (g *Group) GET(name string, handler ...interface{}) *Group {
	g.RouterGroup.GET(name, g.serv.parseFunc(true, handler...)...)
	return g
}

func (g *Group) POST(name string, handler ...interface{}) *Group {
	g.RouterGroup.POST(name, g.serv.parseFunc(true, handler...)...)
	return g
}

func (g *Group) PUT(name string, handler ...interface{}) *Group {
	g.RouterGroup.PUT(name, g.serv.parseFunc(true, handler...)...)
	return g
}

func (g *Group) PATCH(name string, handler ...interface{}) *Group {
	g.RouterGroup.PATCH(name, g.serv.parseFunc(true, handler...)...)
	return g
}

func (g *Group) ANY(name string, handler ...interface{}) *Group {
	g.RouterGroup.Any(name, g.serv.parseFunc(true, handler...)...)
	return g
}

func (g *Group) DELETE(name string, handler ...interface{}) *Group {
	g.RouterGroup.DELETE(name, g.serv.parseFunc(true, handler...)...)
	return g
}
