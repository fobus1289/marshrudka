package http_gin

import (
	"log"
	"reflect"
)

type Throw struct {
	StatusCode  int
	Data        interface{}
	ContentType string
}

type Response struct {
	StatusCode  int
	Data        interface{}
	ContentType string
}

func (s *Serv) Dep(owner interface{}) {

	ownerType := reflect.ValueOf(owner)

	if ownerType.Kind() == reflect.Ptr {
		ownerType = ownerType.Elem()
	}

	for i := 0; i < ownerType.NumField(); i++ {
		fieldType := ownerType.Field(i)

		service := s.services[fieldType.Type()]

		if service.Kind() != reflect.Invalid {
			fieldType.Set(service)
		}
	}
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
