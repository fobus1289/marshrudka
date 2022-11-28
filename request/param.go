package request

import (
	"github.com/fobus1289/marshrudka/util"
)

func (r *request) ParseParam() {
	if r.ParamMap == nil {
		r.ParamMap = map[string]string{}
	}
}

func (r *request) HasParam(key string) bool {
	return r.Param(key) != ""
}

func (r *request) Param(key string) string {
	r.ParseParam()
	return r.ParamMap[key]
}

func (r *request) ParamDefault(key string, _default string) string {
	if param := r.Param(key); param != "" {
		return param
	}
	return _default
}

func (r *request) ParamInt(key string) int64 {
	return util.ToInt(r.Param(key), 0)
}

func (r *request) ParamIntDefault(key string, _default int64) int64 {
	return util.ToInt(r.Param(key), _default)
}

func (r *request) ParamFloat(key string) float64 {
	return util.ToFloat(r.Param(key), 0)
}

func (r *request) ParamFloatDefault(key string, _default float64) float64 {
	return util.ToFloat(r.Param(key), _default)
}

func (r *request) ParamIntClamp(key string, min, max int64) int64 {
	value := util.ToInt(r.Param(key), min)

	switch {
	case value > max:
		return max
	case value < min:
		return min
	default:
		return value
	}
}

func (r *request) ParamFloatClamp(key string, min, max float64) float64 {
	value := util.ToFloat(r.Param(key), min)

	switch {
	case value > max:
		return max
	case value < min:
		return min
	default:
		return value
	}
}

func (r *request) TryGetParam(key string, out any) bool {
	var query = r.Param(key)
	return stringToAny(query, out)
}
