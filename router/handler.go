package router

import (
	"net/http"
	"reflect"
)

type handler struct {
	*Drive
	*route
	last   bool
	params []reflect.Type
	call   func([]reflect.Value) []reflect.Value
}

type handlers []*handler

func (h handlers) setRoutes(r *route) {
	for _, h2 := range h {
		h2.route = r
	}
}

func (h handlers) each(w http.ResponseWriter, r *http.Request, refMap reflectMap) (next bool) {

	for _, handler := range h {

		var retValues []reflect.Value

		for _, param := range handler.params {

			if reflect.DeepEqual(param, _request) {
				request := &Request{
					Request:  r,
					Response: w,
					params:   map[string]string{},
				}

				if handler.route == nil {
					reqVl := reflect.ValueOf(request)
					retValues = append(retValues, reqVl)
					refMap[param] = reqVl
					continue
				}

				route := handler.route

				path := route.match.FindString(r.URL.Path)
				matches := route.match.FindStringSubmatch(r.URL.Path)

				for _, pr := range route.paramNames {
					for _, match := range matches {
						if path == match {
							continue
						}
						request.params[pr] = match
						break
					}
				}

				reqVl := reflect.ValueOf(request)
				retValues = append(retValues, reqVl)
				refMap[param] = reqVl
				continue
			}

			ref := refMap[param]

			if ref.Kind() != invalid {
				retValues = append(retValues, ref)
				continue
			}

			service := handler.services[param]

			if service.Kind() != invalid {
				retValues = append(retValues, service)
				refMap[param] = service
				continue
			}

			if reflect.DeepEqual(param, _httpReq) {
				req := reflect.ValueOf(r)
				retValues = append(retValues, req)
				refMap[req.Type()] = req
				continue
			}

			if reflect.DeepEqual(param, _httpRes) {
				res := reflect.ValueOf(w)
				retValues = append(retValues, res)
				refMap[res.Type()] = res
				continue
			}

			jsonVal := setOther(param, r, w)

			if jsonVal == nil {
				return false
			}

			jsonValu := *jsonVal
			refMap[jsonValu.Type()] = jsonValu
			retValues = append(retValues, jsonValu)
		}

		if !handler.ret(retValues, refMap, w, r) {
			return false
		}

	}

	return true
}

func (h *handler) ret(retValues []reflect.Value, refMap reflectMap, w http.ResponseWriter, r *http.Request) bool {
	ret := h.call(retValues)

	retLen := len(ret) < 1

	if retLen && h.last {
		return false
	} else if retLen {
		return true
	}

	if !h.last {
		for _, value := range ret {
			val := value.Interface()

			if val == nil {
				continue
			}

			_value := reflect.ValueOf(val)

			if isThrow(_value, w) {
				return false
			}

			refMap[_value.Type()] = _value
		}
		return true
	}

	val := ret[0].Interface()

	if val == nil {
		return false
	}

	_value := reflect.ValueOf(val)

	if isThrow(_value, w) || isResponse(_value, w) || isFileResponse(_value, w, r) {
		return false
	}

	data := valueBytes(_value)
	_, _ = w.Write(data)
	return false
}
