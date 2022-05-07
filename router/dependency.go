package router

import (
	"log"
	"reflect"
)

func (s *server) SetService(v interface{}) bool {

	if v == nil {
		return false
	}

	var refValue = reflect.ValueOf(v)
	{
		if refValue.Kind() != reflect.Ptr {
			return false
		}

		var (
			key   = refValue.Elem()
			value = key
		)

		if !key.IsValid() || key.Type().PkgPath() == "" {
			return false
		}

		if key.Kind() == reflect.Interface {
			if value = key.Elem(); value.Kind() == reflect.Invalid {
				return false
			}
			s.services[key.Type()] = value
			s.services[value.Type()] = value
			return true
		}

	}

	s.services[refValue.Type()] = refValue

	return true
}

func (s *server) GetService(out interface{}) bool {

	if out == nil {
		return false
	}

	var key = reflect.ValueOf(out)
	{
		if key.Kind() == reflect.Invalid {
			return false
		}

		if key.Kind() == reflect.Ptr {
			if key = key.Elem(); key.Kind() == reflect.Invalid {
				return false
			}
		}
	}

	if key.Kind() == reflect.Interface {

		for k, value := range s.services {
			if k.Implements(key.Type()) && key.CanSet() {
				s.services[key.Type()] = value
				key.Set(value)
				return true
			}
		}

		return false
	}

	if value := s.services[key.Type()]; value.Kind() != reflect.Invalid {
		if key.CanSet() {
			key.Set(value)
			return true
		}
	}

	return false
}

func (s *server) SetServices(services ...interface{}) bool {
	for _, service := range services {
		if !s.SetService(service) {
			return false
		}
	}
	return true
}

func (s *server) GetServices(services ...interface{}) bool {
	for _, service := range services {
		if !s.GetService(service) {
			return false
		}
	}
	return true
}

func (s *server) FillServiceFields(service interface{}) bool {

	if service == nil {
		return false
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	var (
		serviceValue = reflect.ValueOf(service)
		serviceElm   = serviceValue
		serviceLast  = serviceElm
	)

	if serviceValue.Kind() != reflect.Ptr {
		return false
	}

	if serviceElm = serviceValue.Elem(); !serviceElm.IsValid() {
		return false
	}

	if serviceElm.Kind() == reflect.Interface {
		serviceElm = serviceElm.Elem()
		if serviceElm.Kind() == reflect.Ptr {
			serviceLast = serviceElm.Elem()
		}
	} else {
		serviceLast = serviceElm
	}

	if serviceLast.NumField() == 0 {
		return true
	}

	for i := 0; i < serviceLast.NumField(); i++ {
		field := serviceLast.Field(i)
		if service, ok := s.services[field.Type()]; ok {
			if field.CanSet() {
				field.Set(service)
			}
		}
	}

	return true
}

func (s *server) FillServicesFields(services ...interface{}) bool {
	for _, service := range services {
		if s.FillServiceFields(service) {
			return false
		}
	}
	return true
}

func (s *server) GetByType(t reflect.Type) reflect.Value {

	var inType = t

	if value := s.services[inType]; value.Kind() != reflect.Invalid {
		return value
	}

	if t.PkgPath() == "" {
		return reflect.Value{}
	}

	for k, v := range s.services {
		if k.AssignableTo(inType) {
			s.services[inType] = v
			return v
		}
	}

	return reflect.Value{}
}
