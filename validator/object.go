package validator

import (
	"fmt"
)

type IObjectValidator interface {
}

type Object struct {
	Value IValidator
	Key   string
	Message
	Method   string
	Child    []map[string]map[string]any
	Nested   bool
	Optional bool
}

func ObjectValidator(key string, t IValidator, methods ...string) *Object {
	s := t

	var method string
	{
		if len(methods) != 0 {
			method = methods[0]
		}
	}

	if f := fmt.Sprintf("%v", s); f == "<nil>" {
		s = nil
	}

	return &Object{
		Value:    s,
		Key:      key,
		Nested:   false,
		Message:  Message{},
		Method:   method,
		Optional: false,
	}
}

func (o *Object) Options(optional, nested bool) *Object {
	o.Optional = optional
	o.Nested = nested
	return o
}

func (o *Object) Required(message string) *Object {

	if o.Optional && o.Value == nil {
		return o
	}

	if o.Value == nil {
		o.Add(o.Key, "required", message, o.Value)
	}

	return o
}

func (o *Object) IsValid() bool {
	return o.Len() == 0
}

func (o *Object) ErrorMessage() Message {

	if !o.Nested {
		return o.Message
	}

	if o.Value == nil {
		return o.Message
	}

	fields := o.Value.Validate(o.Method)

	if fields.Len() == 0 {
		return o.Message
	}

	if o.Key != "" {
		for k, v := range fields {

			if msg := o.Message[o.Key]; msg != nil {
				msg[k] = v
			} else {
				o.Message[o.Key] = map[string]any{
					k: v,
				}
			}
		}
		return o.Message
	}

	for k, v := range fields {
		o.Message[k] = v
	}

	return o.Message
}
