package marshrudka

import "reflect"

type action struct {
	Params []reflect.Type
	Call   func(values []reflect.Value) []reflect.Value
	Ret    bool
}

type actions []*action

func (a *actions) Add(action *action) {
	*a = append(*a, action)
}
