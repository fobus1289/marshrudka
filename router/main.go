package router

import (
	"log"
	"net/http"
	"reflect"
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

func (d *Drive) Run(addr string) {
	log.Fatalln(http.ListenAndServe(addr, d))
}

func (d *Drive) RunAsync(addr string) {
	go func() {
		log.Fatalln(http.ListenAndServe(addr, d))
	}()
}
