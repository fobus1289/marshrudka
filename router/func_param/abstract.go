package func_param

import "reflect"

type IFuncParam interface {
	Find(t reflect.Type) IIS
	Push(t reflect.Type) bool
	Map(fn func(iis IIS))
	Count() int
	All() []IIS
}

type IIS interface {
	GetType() reflect.Type
	IsValid() bool
	CanAlloc() bool
	NewType() reflect.Value
	Has(t reflect.Type) bool
	Get(t reflect.Type) reflect.Type
	GetPointer() reflect.Type
	GetElement() reflect.Type
	FindMap(rm map[reflect.Type]reflect.Value) reflect.Value
}
