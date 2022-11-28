package validator

import (
	"fmt"
	"strings"
)

type AnyType[T, E any] struct {
	value    T
	key      string
	messages map[string]any
}

func (a *AnyType[T, E]) Min(min T, message string) *AnyType[T, E] {
	return a
}

func (a *AnyType[T, E]) Max(max T, message string) *AnyType[T, E] {
	return a
}

func (a *AnyType[T, E]) Lenght(min, max T, message string) *AnyType[T, E] {
	return a
}

func (a *AnyType[T, E]) Key() string {
	return a.key
}

func (a *AnyType[T, E]) IsValid() bool {
	return len(a.messages) == 0
}

func (a *AnyType[T, E]) Message() map[string]any {
	return a.messages
}

func (a *AnyType[T, E]) Format(message string, format E) string {

	if strings.Contains(message, "%v") {
		return fmt.Sprintf(message, format)
	}

	return message
}

func (a *AnyType[T, E]) AddMessage(add bool, name, message string, format E) *AnyType[T, E] {

	if add && a.messages == nil {
		a.messages = make(map[string]any)
	}

	if add {
		a.messages[name] = a.Format(message, format)
	}

	return a
}
