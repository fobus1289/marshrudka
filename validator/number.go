package validator

import "golang.org/x/exp/constraints"

type AnyNumber interface {
	constraints.Integer |
		constraints.Float
}

func NumberValidator[T AnyNumber](key string, t T) *Number[T] {
	return &Number[T]{
		&AnyType[T, T]{
			value: t,
			key:   key,
		},
	}
}

type Number[T constraints.Ordered] struct {
	*AnyType[T, T]
}

func (n *Number[T]) Min(min T, message string) *Number[T] {
	n.AddMessage(n.value < min, "min", message, min)
	return n
}

func (n *Number[T]) Max(max T, message string) *Number[T] {
	n.AddMessage(n.value > max, "max", message, max)
	return n
}

func (n *Number[T]) Required(message string) *Number[T] {
	var _default T
	n.AddMessage(n.value == _default, "required", message, _default)
	return n
}

func (n *Number[T]) Lenght(min, max T, message string) *Number[T] {
	var f T

	var isInvalid bool
	{
		if n.value < min {
			f = min
			isInvalid = true
		} else if n.value > max {
			f = max
			isInvalid = true
		}
	}

	n.AddMessage(isInvalid, "lenght", message, f)

	return n
}
