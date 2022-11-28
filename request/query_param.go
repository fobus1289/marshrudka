package request

import (
	"test/util"
)

func (r *request) HasQuery(key string) bool {
	return r.Query(key) != ""
}

func (r *request) Query(key string) string {
	return r.URL.Query().Get(key)
}

func (r *request) QueryDefault(key string, _default string) string {
	if query := r.Query(key); query != "" {
		return query
	}
	return _default
}

func (r *request) QueryInt(key string) int64 {
	return util.ToInt(r.Query(key), 0)
}

func (r *request) QueryIntDefault(key string, _default int64) int64 {
	return util.ToInt(r.Query(key), _default)
}

func (r *request) QueryFloat(key string) float64 {
	return util.ToFloat(r.Query(key), 0)
}

func (r *request) QueryFloatDefault(key string, _default float64) float64 {
	return util.ToFloat(r.Query(key), _default)
}

func (r *request) QueryIntClamp(key string, min, max int64) int64 {

	value := util.ToInt(r.Query(key), min)

	switch {
	case value > max:
		return max
	case value < min:
		return min
	default:
		return value
	}
}

func (r *request) QueryFloatClamp(key string, min, max float64) float64 {
	value := util.ToFloat(r.Query(key), min)

	switch {
	case value > max:
		return max
	case value < min:
		return min
	default:
		return value
	}
}

func (r *request) TryQueryGetQuery(key string, out any) bool {
	var query = r.Query(key)
	return stringToAny(query, out)
}
