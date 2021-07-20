package router

import (
	"log"
	"reflect"
	"regexp"
	"strings"
)

func (d *Drive) parseFunc(parent bool, actions ...interface{}) handlers {

	var _handlers handlers

	for a, action := range actions {

		funcValue := reflect.ValueOf(action)

		if funcValue.Kind() != reflect.Func {
			log.Fatalln("dont supported this type:", funcValue.Kind())
		}

		funcType := funcValue.Type()

		var params []reflect.Type

		for i := 0; i < funcType.NumIn(); i++ {
			in := funcType.In(i)
			params = append(params, in)
		}

		var handler = handler{
			Drive:  d,
			last:   len(actions)-1 == a && !parent,
			params: params,
			call:   funcValue.Call,
		}

		_handlers = append(_handlers, &handler)
	}

	return _handlers
}

func parsePath(path string) (*regexp.Regexp, []string) {

	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}
	if strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}

	path, names := pattern(path+"/", []string{})

	path = multiplication(path)

	path = "/?" + path

	return regexp.MustCompile(`^(` + path + `)$`), names
}

func pattern(path string, names []string) (string, []string) {

	dotIndex := strings.Index(path, ":")

	if dotIndex == -1 {
		return path, names
	}

	hasOne := path[dotIndex:]
	slashIndex := strings.Index(hasOne, "/")

	if slashIndex == -1 {
		return path, names
	}

	hasTwo := hasOne[:slashIndex]

	hasPattern := `(\w+)`

	patternIndex := strings.Index(hasTwo, "{")

	name := hasTwo[1:]

	if patternIndex != -1 {
		pat := hasTwo[patternIndex:]

		if pat[1] == 115 {
			hasPattern = `([A-z]+)`
		}

		if pat[1] == 110 {
			hasPattern = `([0-9]+)`
		}
		name = name[:patternIndex-1]
	}

	names = append(names, name)
	path = strings.Replace(path, hasTwo, hasPattern, 1)

	return pattern(path, names)
}

func multiplication(path string) string {
	multIndex := strings.Index(path, "*")

	if multIndex == -1 {
		return path + "?"
	}

	return path[:multIndex] + `?(\S+)?`
}
