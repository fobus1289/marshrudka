package router

import (
	"reflect"
)

func (s *server) AddScoped(scoped any) IServer {
	s.Scopeds = append(s.Scopeds, scoped)
	return s
}

func (s *server) AddSingleton(singleton any) IServer {
	s.Singletons = append(s.Singletons, singleton)
	return s
}

func (s *server) UseService() {

	for _, singleton := range s.Singletons {

		value := reflect.ValueOf(singleton)
		{
			if value.Kind() == reflect.Invalid {
				panic("UseService invalid type")
			}
		}

		if value.Kind() == reflect.Func {

			valueType := value.Type()
			{
				if valueType.NumOut() != 1 || valueType.NumIn() > 0 {
					panic("valueType.NumOut() != 1 || valueType.NumIn() > 0")
				}
			}

			outType := valueType.Out(0)

			s.Services[outType] = func(hp *handlerParam) (reflect.Value, *RuntimeError) {
				persistent := s.PersistentService[outType]
				if persistent.Kind() == reflect.Invalid {
					persistent = value.Call(nil)[0]
				}
				return persistent, nil
			}

			continue
		}
		s.Services[value.Type()] = func(hp *handlerParam) (reflect.Value, *RuntimeError) {
			return value, nil
		}

	}

	for _, scoped := range s.Scopeds {

		value := reflect.ValueOf(scoped)
		{
			if value.Kind() == reflect.Invalid {
				panic("UseService invalid type")
			}
		}

		if value.Kind() != reflect.Func {

			valueType := value.Type()

			for value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
				value = value.Elem()
			}

			if value.Kind() == reflect.Invalid {
				panic("valueType2.Kind() == reflect.Invalid")
			}

			s.Services[valueType] = func(hp *handlerParam) (reflect.Value, *RuntimeError) {

				session := hp.SessionData[valueType]

				if session.Kind() == reflect.Invalid {
					newValue := reflect.New(value.Type())
					newValue.Elem().Set(value)
					session = newValue
					hp.SessionData[valueType] = session
				}

				if valueType.Kind() == reflect.Ptr {
					return session, nil
				}

				return session.Elem(), nil
			}

			continue
		}

		if value.Kind() == reflect.Func {

			valueType := value.Type()
			{
				if valueType.NumOut() != 1 || valueType.NumIn() > 0 {
					panic("valueType.NumOut() != 1 || valueType.NumIn() > 0")
				}
			}

			outType := valueType.Out(0)

			s.Services[outType] = func(hp *handlerParam) (reflect.Value, *RuntimeError) {
				session := hp.SessionData[outType]
				if session.Kind() == reflect.Invalid {
					session = value.Call(nil)[0]
				}
				return session, nil
			}
			continue
		}

	}

}
