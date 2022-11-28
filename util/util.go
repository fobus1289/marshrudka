package util

import (
	"encoding/json"
	"encoding/xml"
	"reflect"
	"strconv"
)

func ToFloat(value string, _default float64) float64 {
	if val, err := strconv.ParseFloat(value, 64); err == nil {
		return val
	} else {
		return _default
	}
}

func ToInt(value string, _default int64) int64 {
	if val, err := strconv.ParseInt(value, 10, 64); err == nil {
		return val
	} else {
		return _default
	}
}

func InjectPrimitiveValue(value string, in any) bool {
	var inValue = reflect.ValueOf(in)

	if !inValue.IsValid() || inValue.Kind() != reflect.Ptr {
		return false
	}

	var inTmp = inValue.Elem()

	switch t := inTmp.Interface().(type) {
	case string:
		inTmp.SetString(value)
	case bool:
		var _bool, err = strconv.ParseBool(value)
		if err != nil {
			return false
		}
		inTmp.SetBool(_bool)
	case int8, int16, int32, int, int64:
		var _int, err = strconv.ParseInt(value, 10, BitSize(t))
		if err != nil {
			return false
		}
		inTmp.SetInt(_int)
	case uint8, uint16, uint32, uint, uint64:
		var _uint, err = strconv.ParseUint(value, 10, BitSize(t))
		if err != nil {
			return false
		}
		inTmp.SetUint(_uint)
	case float32, float64:
		var _float, err = strconv.ParseFloat(value, BitSize(t))
		if err != nil {
			return false
		}
		inTmp.SetFloat(_float)
	default:
		return false
	}

	return true
}

func BitSize(in any) int {
	switch in.(type) {
	case int8, uint8:
		return 8
	case int16, uint16:
		return 16
	case int32, uint32, float32, int, uint:
		return 32
	case int64, uint64, float64:
		return 64
	}
	return -1
}

func IsPtr(in any) bool {

	if in == nil {
		return false
	}

	return reflect.ValueOf(in).Kind() == reflect.Ptr
}

func IsValid(in any) bool {
	if in == nil {
		return false
	}
	return reflect.ValueOf(in).Kind() != reflect.Invalid
}

func IsEquals(in any, kinds ...reflect.Kind) bool {

	if !IsValid(in) {
		return false
	}

	valueKind := Elem(in).Kind()

	for _, kind := range kinds {
		if valueKind == kind {
			return true
		}
	}

	return false
}

func IsPtrOrInterface(in any) bool {
	if in == nil {
		return false
	}

	kind := reflect.ValueOf(in).Kind()

	return kind == reflect.Ptr || kind == reflect.Interface
}

func Elem(in any) reflect.Value {

	value := reflect.ValueOf(in)

	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	return value
}

func SerializeJson(in any, body *[]byte) {
	if !IsEquals(in, reflect.Struct, reflect.Array, reflect.Slice, reflect.Map) {
		return
	}

	if data, err := json.Marshal(in); err == nil {
		*body = data
	}
}

func SerializeXml(in any, body *[]byte) {
	if !IsEquals(in, reflect.Struct, reflect.Array, reflect.Slice, reflect.Map) {
		return
	}

	if data, err := xml.Marshal(in); err == nil {
		*body = data
	}
}
