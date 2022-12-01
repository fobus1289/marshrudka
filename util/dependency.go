package util

import (
	"log"
	"reflect"
)

func ParseFunc(fn any) (reflect.Type, reflect.Value) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	value := reflect.ValueOf(fn)
	{
		if value.Kind() != reflect.Func {
			panic("value.Kind() != reflect.Func")
		}
	}

	valueType := value.Type()
	{
		outNumb := valueType.NumOut()

		if outNumb == 0 || outNumb > 1 {
			panic("outNumb == 0 || outNumb > 1")
		}
	}

	out := valueType.Out(0)
	{
		tount := out.Elem()
		{
			if out.Kind() == reflect.Ptr {
				tount = out.Elem()
			}
		}

		switch tount.Kind() {
		case reflect.Interface:
			if out.PkgPath() == "" {
				panic("can be return any")
			}
		case reflect.Struct:
		default:
			panic("type not supported")
		}
	}

	return out, value
}
