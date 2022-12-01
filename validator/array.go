package validator

type Array[T AnyArray] struct {
	Key   string
	Value []T
	Message
	Optional bool
}

func ArrayValidator[T AnyArray](key string, value []T) *Array[T] {
	return &Array[T]{
		Key:      key,
		Value:    value,
		Message:  Message{},
		Optional: false,
	}
}

func (arr *Array[T]) Options(optional bool) *Array[T] {
	arr.Optional = optional
	return arr
}

func (arr *Array[T]) Min(min int, message string) *Array[T] {

	if arr.Optional && arr.Value == nil {
		return arr
	}

	if len(arr.Value) < min {
		arr.Add(arr.Key, "min", message, arr.Value)
	}
	return arr
}

func (arr *Array[T]) Max(max int, message string) *Array[T] {

	if arr.Optional && arr.Value == nil {
		return arr
	}

	if len(arr.Value) > max {
		arr.Add(arr.Key, "max", message, arr.Value)
	}
	return arr
}

func (arr *Array[T]) Null(message string) *Array[T] {

	if arr.Optional && arr.Value == nil {
		return arr
	}

	if arr.Value == nil {
		arr.Add(arr.Key, "null", message, arr.Value)
	}

	return arr
}

func (arr *Array[T]) Empty(message string) *Array[T] {

	if arr.Optional && arr.Value == nil {
		return arr
	}

	if len(arr.Value) == 0 {
		arr.Add(arr.Key, "empty", message, arr.Value)
	}

	return arr
}

func (arr *Array[T]) Lenght(min, max int, message string) *Array[T] {

	if arr.Optional && arr.Value == nil {
		return arr
	}

	l := len(arr.Value)

	if l < min || l > max {
		arr.Add(arr.Key, "lenght", message, arr.Value)
	}

	return arr
}

func (arr *Array[T]) Omit(message string, a T) *Array[T] {

	if arr.Optional && arr.Value == nil {
		return arr
	}

	for _, v := range arr.Value {
		if v == a {
			arr.Add(arr.Key, "omit", message, v)
			return arr
		}
	}

	return arr
}

func (arr *Array[T]) Only(message string, a T) *Array[T] {

	if arr.Optional && arr.Value == nil {
		return arr
	}

	for _, v := range arr.Value {
		if v == a {
			return arr
		}
	}

	arr.Add(arr.Key, "only", message, arr.Value)

	return arr
}

func (arr *Array[T]) ErrorMessage() Message {
	return arr.Message
}
