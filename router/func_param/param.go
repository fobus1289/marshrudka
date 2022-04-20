package func_param

import (
	"net/http"
	"reflect"
)

var (
	invalidIs = &is{
		PointerType: reflect.TypeOf(nil),
		ElementType: reflect.TypeOf(nil),
		Type:        Invalid,
	}
)

type Kind int8

const (
	Invalid Kind = -1
	Number  Kind = iota
	String
	Boolean
	NativeInterface
	Array
	CustomInterface
	Pointer
	Struct
	Map
)

type is struct {
	PointerType reflect.Type
	ElementType reflect.Type
	Type        Kind
}

type iss []IIS

func NewIsContainer() IFuncParam {
	return &iss{}
}

func (i iss) Count() int {
	return len(i)
}

func (i iss) All() []IIS {
	return i
}

func (i iss) Map(fn func(iis IIS)) {
	for _, _is := range i {
		fn(_is)
	}
}

func (i iss) Find(t reflect.Type) IIS {

	for _, _is := range i {
		switch t {
		case _is.GetElement(), _is.GetPointer():
			return _is
		}
	}

	return invalidIs
}

func (i *iss) Push(t reflect.Type) bool {
	if n := newIs(t); n.IsValid() {
		*i = append(*i, n)
		return true
	}
	return false
}

func newIs(in interface{}) *is {

	var inValue = reflect.ValueOf(in)
	{
		if inValue.Kind() == reflect.Invalid {
			return invalidIs
		}
	}

	var _is = &is{}
	{
		var pointerType = inValue.Type()

		switch pointerType.Kind() {
		case reflect.Ptr:
			_is.Type = Pointer
		case reflect.String:
			_is.Type = String
		case reflect.Bool:
			_is.Type = Boolean
		case reflect.Slice:
			_is.Type = Array
		case reflect.Map:
			_is.Type = Map
		case reflect.Struct:
			_is.Type = Struct
		case reflect.Interface:
			_is.Type = CustomInterface
		case
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Uintptr:
			_is.Type = Number
		default:
			return invalidIs
		}

		if pointerType.Kind() == reflect.Ptr {
			_is.ElementType = inValue.Elem().Type()
		}

		if pointerType.Kind() == reflect.Interface && pointerType.PkgPath() == "" {
			_is.Type = NativeInterface
		}

	}

	return _is
}

func (i *is) GetType() reflect.Type {
	return i.PointerType
}

func (i *is) IsValid() bool {
	return i.Type != -1
}

func (i *is) CanAlloc() bool {
	switch i.Type {
	case Invalid, Number, String, Boolean:
		return false
	}
	return true
}

func (i *is) NewType() reflect.Value {
	if !i.IsValid() || !i.CanAlloc() {
		return reflect.Value{}
	}
	return reflect.New(i.PointerType)
}

func (i *is) Has(t reflect.Type) bool {
	switch t {
	case i.ElementType, i.PointerType:
		return true
	}
	return false
}

func (i *is) Get(t reflect.Type) reflect.Type {
	switch t {
	case i.ElementType:
		return i.ElementType
	case i.PointerType:
		return i.PointerType
	default:
		return invalidIs.PointerType
	}
}

func (i *is) GetPointer() reflect.Type {
	return i.PointerType
}

func (i *is) GetElement() reflect.Type {
	return i.ElementType
}

func (i *is) FindMap(rm map[reflect.Type]reflect.Value) reflect.Value {
	if len(rm) < 1 {
		return reflect.Value{}
	}
	for k, v := range rm {
		switch k {
		case i.ElementType:
			return v
		case i.PointerType:
			return v
		}
	}
	return reflect.Value{}
}

func (i iss) Fill(w http.ResponseWriter, r *http.Request, rm map[reflect.Type]reflect.Value) {
	//var params []reflect.Value

	//for _, _iss := range i {
	//
	//	//if prm := rm[param]; prm.Kind() != reflect.Invalid {
	//	//	params = append(params, prm)
	//	//	continue
	//	//}
	//	//
	//	//switch _iss.PointerType {
	//	//case httpRes:
	//	//	var res = reflect.ValueOf(w)
	//	//	rm[res.Type()] = res
	//	//	params = append(params, res)
	//	//	continue
	//	//case httpReq:
	//	//	var req = reflect.ValueOf(r)
	//	//	rm[req.Type()] = req
	//	//	params = append(params, req)
	//	//	continue
	//	//case request:
	//	//	var formFileValue = reflect.ValueOf(request2.NewRequest(w, r, h.Router.Params, h.Router.Match))
	//	//	rm[formFileValue.Type()] = formFileValue
	//	//	params = append(params, formFileValue)
	//	//	continue
	//	//case iFormFile:
	//	//	var formFileValue = reflect.ValueOf(request2.NewFormFile(w, r))
	//	//	rm[formFileValue.Type()] = formFileValue
	//	//	params = append(params, formFileValue)
	//	//	continue
	//	//}
	//	//
	//	//if service := h.Server.GetByType(param); service.Kind() != reflect.Invalid {
	//	//	rm[service.Type()] = service
	//	//	params = append(params, service)
	//	//	continue
	//	//}
	//	//
	//	//if object := read(param, r); object.Kind() != reflect.Invalid {
	//	//	rm[object.Type()] = object
	//	//	params = append(params, object)
	//	//}
	//
	//}
}
