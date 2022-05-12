package request

import (
	"strings"
)

func (r *request) Query(key string) String {
	return String(r.r.URL.Query().Get(key))
}

func (r *request) HasQuery(key string) bool {
	return r.r.URL.Query().Get(key) != ""
}

func (r *request) TryGetQuery(key string, in interface{}) bool {
	var value = r.r.URL.Query().Get(key)

	if strings.TrimSpace(value) == "" {
		return false
	}

	return r.setType(value, in)
}
