package dependency

import (
	"log"
	"reflect"
)

func parseFunc(fn any) {

}

func parseStruct(st any) {

}

func ParseServiceFunc(fn any) (reflect.Type, func() reflect.Value) {

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

	return out, func() reflect.Value {
		return value.Call(nil)[0]
	}

}
