package router

import (
	"reflect"
	"test/response"
)

var (
	responseType = reflect.TypeOf((*response.IResponse)(nil)).Elem()
	sendType     = reflect.TypeOf((*response.ISend)(nil)).Elem()
	readyType    = reflect.TypeOf((*response.IReady)(nil)).Elem()
)

func perpareActionHandler(hs handlers) Call {

	l := len(hs)

	if l == 1 {
		return oneHandler(hs[0])
	}

	lastIndex := l - 1

	ths := hs[:lastIndex]

	lastHandler := oneHandler(hs[lastIndex])

	actions := moreHandler(ths)

	return func(param *handlerParam) (stop bool) {

		for _, action := range actions {
			if action(param) {
				return true
			}
		}

		return lastHandler(param)
	}
}

func oneHandler(hs *handler) Call {

	h := hs

	outTypes := h.OutTypes

	numOut := len(outTypes)

	if numOut > 1 {
		panic("can't send more than 1")
	}

	if numOut == 0 {
		return func(param *handlerParam) (stop bool) {
			values, stop := runtimeCheck(h, param)
			if stop {
				return true
			}

			h.Call(values)

			return true
		}
	}

	outType := outTypes[0]

	switch outType {
	case responseType, sendType, readyType:
		return aborted(h, outType, 0, true)
	}

	outAction, contentType := preparationOut(outTypes[0])

	return func(param *handlerParam) (stop bool) {
		values, stop := runtimeCheck(h, param)

		if stop {
			return true
		}

		out := h.Call(values)[0]
		w := param.ResponseWriter
		w.Header().Add("Content-Type", contentType)
		w.Write(outAction(out))
		return true
	}

}

func runtimeCheck(h *handler, param *handlerParam) ([]reflect.Value, bool) {

	values, err := h.Params.Get(param)

	if err != nil {
		w := param.ResponseWriter
		w.Header().Add("Content-Type", err.ContentType)
		w.WriteHeader(err.Status)
		w.Write(err.Data)
		return nil, true
	}

	return values, false
}

func moreHandler(hs handlers) []Call {
	var calles []Call

	for _, h := range hs {

		if len(h.OutTypes) == 0 {
			h1 := h
			calles = append(calles, func(param *handlerParam) (stop bool) {

				inValues, stop := runtimeCheck(h1, param)

				if stop {
					return true
				}

				values := h1.Call(inValues)

				sessionData := param.SessionData

				for _, value := range values {
					if session := sessionData[value.Type()]; session.Kind() != reflect.Invalid {
						continue
					}
					sessionData[value.Type()] = value
				}

				return false
			})
			continue
		}

		for i, outType := range h.OutTypes {
			switch outType {
			case responseType, sendType, readyType:
				index := i
				h1 := h
				outType1 := outType
				calles = append(calles, aborted(h1, outType1, index, false))
			default:
				h1 := h
				calles = append(calles, func(hp *handlerParam) bool {

					inValues, stop := runtimeCheck(h1, hp)

					if stop {
						return true
					}

					values := h1.Call(inValues)

					sessionData := hp.SessionData

					for _, value := range values {
						if session := sessionData[value.Type()]; session.Kind() != reflect.Invalid {
							continue
						}
						sessionData[value.Type()] = value
					}

					return false
				})
			}
		}

	}

	return calles
}

func aborted(h *handler, outType reflect.Type, index int, last bool) Call {

	return func(hp *handlerParam) bool {

		inValues, stop := runtimeCheck(h, hp)

		if stop {
			return true
		}

		values := h.Call(inValues)

		readyValue := values[index]
		{
			if readyValue.Kind() == reflect.Invalid {
				return false
			}
		}

		readyInteface := readyValue.Interface()

		if ready, ok := readyInteface.(response.IReady); ok {
			if ready.HasAbort() || last {
				w := hp.ResponseWriter
				w.Header().Set("Content-Type", ready.ContentType())
				w.WriteHeader(ready.GetStatus())
				if ready.HasBody() {
					w.Write(ready.GetBody())
				}
				return true
			}
			hp.SessionData[outType] = readyValue
		}

		return false
	}
}
