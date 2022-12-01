package validator

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)

type AnyArray interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

func format(message string, format any) string {

	if strings.Contains(message, "%v") {
		if format == nil {
			format = "<nil>"
		}
		return fmt.Sprintf(message, format)
	}

	return message
}

type Message map[string]map[string]any

func (msg Message) Add(key, name, value string, formatV any) {

	if ms := msg[key]; ms != nil {
		ms[name] = format(value, formatV)
	} else {
		msg[key] = map[string]any{
			name: format(value, formatV),
		}
	}

}

func joinKeys(keys ...string) string {
	return strings.Join(keys, ".")
}

func (msg Message) Len() int {
	return len(msg)
}

func isEmailValid(e string) bool {
	return emailRegex.MatchString(e)
}
