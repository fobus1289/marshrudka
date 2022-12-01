package validator

import (
	"golang.org/x/exp/constraints"
)

type AnyNumber interface {
	constraints.Integer |
		constraints.Float
}

func NumberValidator[T AnyNumber](key string, value *T) *Number[T] {
	return &Number[T]{
		key:      key,
		value:    value,
		Optional: false,
		Message:  Message{},
	}
}

type Number[T constraints.Ordered] struct {
	key      string
	value    *T
	Optional bool
	Message
}

func (n *Number[T]) Options(optional bool) *Number[T] {
	n.Optional = optional
	return n
}

func (n *Number[T]) Min(min T, message string) *Number[T] {

	var value T
	{
		if n.value != nil {
			value = *n.value
		}
	}

	if n.Optional && n.value == nil {
		return n
	}

	if value < min {
		n.Add(n.key, "min", message, n.value)
	}
	return n
}

func (n *Number[T]) Max(max T, message string) *Number[T] {
	var value T
	{
		if n.value != nil {
			value = *n.value
		}
	}

	if n.Optional && n.value == nil {
		return n
	}

	if value > max {
		n.Add(n.key, "max", message, n.value)
	}

	return n
}

func (n *Number[T]) Required(message string) *Number[T] {
	var value T
	{
		if n.value != nil {
			value = *n.value
		}
	}

	if n.Optional && n.value == nil {
		return n
	}

	var _default T

	if value == _default {
		n.Add(n.key, "required", message, n.value)
	}

	return n
}

func (n *Number[T]) Lenght(min, max T, message string) *Number[T] {

	var value T
	{
		if n.value != nil {
			value = *n.value
		}
	}

	if n.Optional && n.value == nil {
		return n
	}

	if value < min || value > max {
		n.Add(n.key, "lenght", message, n.value)
	}

	return n
}

func (n *Number[T]) Omit(message string, values ...T) *Number[T] {

	if len(values) == 0 {
		return n
	}

	var value T
	{
		if n.value != nil {
			value = *n.value
		}
	}

	if n.Optional && n.value == nil {
		return n
	}

	for _, v := range values {
		if value == v {
			n.Add(n.key, "omit", message, n.value)
			return n
		}
	}

	return n
}

func (n *Number[T]) Only(message string, values ...T) *Number[T] {

	if len(values) == 0 {
		return n
	}

	var value T
	{
		if n.value != nil {
			value = *n.value
		}
	}

	if n.Optional && n.value == nil {
		return n
	}

	for _, v := range values {
		if value == v {
			return n
		}
	}

	n.Add(n.key, "only", message, n.value)

	return n
}

func (n *Number[T]) ErrorMessage() Message {
	return n.Message
}
