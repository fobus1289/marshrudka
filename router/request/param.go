package request

import (
	"strings"
)

func (r *request) Param(key string) string {
	r.parseParam()
	return r.paramsMap[key]
}

func (r *request) HasParam(key string) bool {
	r.parseParam()
	return r.paramsMap[key] != ""
}

func (r *request) TryGetParam(key string, in interface{}) bool {
	r.parseParam()

	var value = r.paramsMap[key]

	if strings.TrimSpace(value) == "" {
		return false
	}

	return r.setType(value, in)
}

func (r *request) parseParam() {
	r.Do(func() {
		var params = r.params
		if params == nil {
			return
		}

		var (
			keys  = params.Keys
			match = params.Match
		)

		if len(keys) < 1 {
			return
		}

		var httpParams = match.FindStringSubmatch(r.r.URL.Path)

		if len(httpParams) <= 2 {
			return
		}

		httpParams = httpParams[2:]

		if len(httpParams) != len(keys) {
			return
		}

		for i, param := range keys {
			r.paramsMap[param] = httpParams[i]
		}
	})
}
