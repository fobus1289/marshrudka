package router

import (
	"log"
	"net/http"
)

type Drive struct {
	services reflectMap
	routes   routes
	handlers handlers
}

func NewRouter() *Drive {
	return &Drive{
		services: reflectMap{},
		routes:   nil,
		handlers: nil,
	}
}

func (d *Drive) Run(addr string) {
	log.Fatalln(http.ListenAndServe(addr, d))
}

func (d *Drive) RunAsync(addr string) {
	go func() {
		log.Fatalln(http.ListenAndServe(addr, d))
	}()
}
