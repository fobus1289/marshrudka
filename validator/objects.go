package validator

import (
	"fmt"
	"strings"
)

type Objects struct {
	objects  []*Object
	key      string
	messages map[string]any
	nested   bool
}

func ObjectsValidator[T any](key string, t []T, nested bool) *Objects {
	var objects = &Objects{
		key:    key,
		nested: nested,
	}

	for _, v := range t {
		switch a := any(v).(type) {
		case IValidator:
			objects.objects = append(objects.objects, ObjectValidator(key, a, nested))
		}
	}

	return objects
}

func (ob *Objects) IsValid() bool {
	ob.Required("required")

	if !ob.nested {
		return len(ob.messages) == 0
	}

	var invalidCount int

	for _, o := range ob.objects {
		if !o.IsValid() {
			invalidCount++
		}
	}

	return invalidCount == 0 && len(ob.messages) == 0
}

func (ob *Objects) Key() string {
	return ob.key
}

func (ob *Objects) Message() map[string]any {

	if ob.messages == nil {
		ob.messages = make(map[string]any)
	}

	if !ob.nested {
		return ob.messages
	}

	for i, o := range ob.objects {
		if o.IsValid() || len(o.Message()) == 0 {
			continue
		}
		ob.messages[fmt.Sprintf("%s.%d", ob.key, i)] = o.Message()
	}

	return ob.messages
}

func (ob *Objects) Lenght(min, max int, message string) *Objects {

	var f int

	var isInvalid bool
	{
		l := len(ob.objects)

		if l < min {
			f = min
			isInvalid = true
		} else if l > max {
			f = max
			isInvalid = true
		}
	}

	return ob.AddMessage(isInvalid, "lenght", message, f)
}

func (ob *Objects) Required(message string) *Objects {
	return ob.AddMessage(len(ob.objects) == 0, "required", message, len(ob.objects))
}

func (ob *Objects) Min(min int, message string) *Objects {
	return ob.AddMessage(len(ob.objects) < min, "min", message, min)
}

func (ob *Objects) Max(max int, message string) *Objects {
	return ob.AddMessage(len(ob.objects) > max, "max", message, max)
}

func (ob *Objects) AddMessage(add bool, name, message string, lenght int) *Objects {

	if ob.messages == nil {
		ob.messages = make(map[string]any)
	}

	if add {
		if _, ok := ob.messages[name]; !ok {
			ob.messages[name] = ob.Format(lenght, message)
		}
	}

	return ob
}

func (ob *Objects) Format(lenght int, message string) string {

	if strings.Contains(message, "%v") {
		return fmt.Sprintf(message, lenght)
	}

	return message
}
