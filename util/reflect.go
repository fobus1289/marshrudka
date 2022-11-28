package util

import (
	"errors"
	"reflect"
)

var (
	singleReflectValue = reflect.Value{}
)

type IReflection interface {
}

type reflection struct {
	Value reflect.Value
}

func NewReflection(value any) (IReflection, error) {
	val := reflect.ValueOf(value)
	{
		if val.Kind() == reflect.Invalid {
			return nil, errors.New("reflect.Invalid")
		}
	}

	return &reflection{
		Value: val,
	}, nil
}

func (r *reflection) Is(kind reflect.Kind) bool {
	return r.Value.Kind() == kind
}

func (r *reflection) IsPtr() bool {
	return r.Is(reflect.Ptr)
}

func (r *reflection) IsInterface() bool {
	return r.Is(reflect.Interface)
}

func (r *reflection) IsEquals(kinds ...reflect.Kind) bool {

	valueKind := r.Value.Kind()

	for _, kind := range kinds {
		if valueKind == kind {
			return true
		}
	}

	return false
}

func (r *reflection) Omit(kinds ...reflect.Kind) bool {
	return r.IsEquals(kinds...) == false
}

func (r *reflection) Only(kinds ...reflect.Kind) bool {
	return r.IsEquals(kinds...)
}

func (r *reflection) Elem() reflect.Value {

	if !r.IsPtr() || !r.IsInterface() {
		return singleReflectValue
	}

	value := r.Value

	for r.IsPtr() || r.IsInterface() {
		value = value.Elem()
	}

	return value
}

func (r *reflection) Set(v any) bool {

	value := r.Elem()
	{
		if value.Kind() == reflect.Invalid || !value.CanSet() || !value.CanAddr() {
			return false
		}
	}

	return true
}
