package util

type M[T ~string] map[string]T

const (
	key = "~~key~~"
)

func Any[T string](key string, value T) M[T] {
	return M[T]{
		key: value,
	}
}

func (m M[T]) Min(min int, message string) M[T] {

	if len(m[key]) < min {
		m["min"] = T(message)
	}

	return m
}

func (m M[T]) Max(max int, message string) M[T] {

	if len(m[key]) > max {
		m["max"] = T(message)
	}

	return m
}
