package validator

type IObjectsValidator interface {
}

type Objects struct {
	Objects []*Object
	Key     string
	Method  string
	Message
	Nested   bool
	Optional bool
}

func ObjectsValidator[T any](key string, t []T, methods ...string) *Objects {

	var method string
	{
		if len(methods) != 0 {
			method = methods[0]
		}
	}

	if len(t) == 0 || t == nil {
		return &Objects{
			Key:      key,
			Objects:  nil,
			Nested:   false,
			Method:   method,
			Message:  Message{},
			Optional: false,
		}
	}

	var objects []*Object

	for _, v := range t {
		switch a := any(v).(type) {
		case IValidator:
			objects = append(objects,
				ObjectValidator(
					"",
					a,
					method,
				),
			)
		}
	}

	return &Objects{
		Key:      key,
		Nested:   false,
		Objects:  objects,
		Message:  Message{},
		Method:   method,
		Optional: false,
	}
}

func (ob *Objects) Options(optional, nested bool) *Objects {
	ob.Optional = optional
	ob.Nested = nested
	for _, o := range ob.Objects {
		o.Nested = nested
		o.Optional = optional
	}
	return ob
}

func (ob *Objects) Null(message string) *Objects {

	if ob.Optional && ob.Objects == nil {
		return ob
	}

	if ob.Objects == nil {
		ob.Add(ob.Key, "null", message, ob.Objects)
	}

	return ob
}

func (ob *Objects) Lenght(min, max int, message string) *Objects {

	if ob.Optional && ob.Objects == nil {
		return ob
	}

	l := len(ob.Objects)

	if l < min || l > max {
		ob.Add(ob.Key, "lenght", message, l)
	}

	return ob
}

func (ob *Objects) Min(min int, message string) *Objects {

	if ob.Optional && ob.Objects == nil {
		return ob
	}

	l := len(ob.Objects)

	if l < min {
		ob.Add(ob.Key, "min", message, l)
	}

	return ob
}

func (ob *Objects) Max(max int, message string) *Objects {

	if ob.Optional && ob.Objects == nil {
		return ob
	}

	l := len(ob.Objects)

	if l > max {
		ob.Add(ob.Key, "max", message, l)
	}

	return ob
}

func (ob *Objects) ErrorMessage() Message {

	if !ob.Nested {
		return ob.Message
	}

	if len(ob.Objects) == 0 {
		return ob.Message
	}

	messages := []any{}

	for _, object := range ob.Objects {
		if object.Value == nil {
			continue
		}

		fields := object.ErrorMessage()

		if len(fields) == 0 {
			continue
		}

		messages = append(messages, fields)
	}

	if len(messages) == 0 {
		return ob.Message
	}

	if message := ob.Message[ob.Key]; message != nil {
		message["items"] = messages
	} else {
		ob.Message[ob.Key] = map[string]any{
			"items": messages,
		}
	}

	return ob.Message
}
