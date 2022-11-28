package validator

type AnyArray interface {
	~[]int | ~[]int8 | ~[]int16 | ~[]int32 | ~[]int64 |
		~[]uint | ~[]uint8 | ~[]uint16 | ~[]uint32 | ~[]uint64 | ~[]uintptr |
		~[]float32 | ~[]float64 |
		~[]string | ~string
}

type Array[T AnyArray] struct {
	*AnyType[T, int]
}

func ArrayValidator[T AnyArray](key string, t T) *Array[T] {
	array := &Array[T]{
		&AnyType[T, int]{
			value: t,
			key:   key,
		},
	}
	return array
}

func (ar *Array[T]) Min(min int, message string) *Array[T] {
	ar.AddMessage(len(ar.value) < min, "min", message, min)
	return ar
}

func (ar *Array[T]) Max(max int, message string) *Array[T] {
	ar.AddMessage(len(ar.value) > max, "max", message, max)
	return ar
}

func (ar *Array[T]) Required(message string) *Array[T] {
	l := len(ar.value)
	ar.AddMessage(l == 0, "required", message, l)
	return ar
}

func (ar *Array[T]) Lenght(min, max int, message string) *Array[T] {
	var f int

	var isInvalid bool
	{
		l := len(ar.value)

		if l < min {
			f = min
			isInvalid = true
		} else if l > max {
			f = max
			isInvalid = true
		}
	}

	ar.AddMessage(isInvalid, "lenght", message, f)

	return ar
}
