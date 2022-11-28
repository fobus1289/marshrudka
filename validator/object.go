package validator

import (
	"fmt"
	"strings"
)

type Object struct {
	Value    IValidator
	key      string
	messages map[string]any
	nested   bool
}

func ObjectValidator(key string, t IValidator, nested bool) *Object {
	s := t

	if f := fmt.Sprintf("%v", s); f == "<nil>" {
		s = nil
	}

	return &Object{
		Value:  s,
		key:    key,
		nested: nested,
	}
}

func (o *Object) Required(message string) *Object {
	return o.AddMessage("required", message)
}

func (o *Object) AddMessage(name, message string) *Object {

	if o.messages == nil {
		o.messages = make(map[string]any)
	}

	if o.Value == nil {
		if _, ok := o.messages[name]; !ok {
			o.messages[name] = o.Format(message)
		}
	}

	return o
}

func (o *Object) Format(message string) string {

	if strings.Contains(message, "%v") {
		return fmt.Sprintf(message, o.Value)
	}

	return message
}

func (o *Object) IsValid() bool {
	o.Required("required")
	return len(o.messages) == 0 && !o.nested
}

func (o *Object) Key() string {
	return o.key
}

func (o *Object) Message() map[string]any {

	if o.Value == nil {
		return o.messages
	}

	if !o.nested {
		return o.messages
	}

	child := o.Value.Validate()

	if o.messages == nil {
		o.messages = make(map[string]any)
	}

	for chk, v := range child {
		o.messages[chk] = v
	}

	return o.messages
}
