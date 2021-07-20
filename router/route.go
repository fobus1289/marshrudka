package router

import (
	"net/http"
	"regexp"
)

type route struct {
	path       string
	match      *regexp.Regexp
	paramNames []string
	methods    map[string]bool
	handlers   handlers
}

type routes []*route

func (r routes) each(w http.ResponseWriter, req *http.Request, refMap reflectMap) {

	var found bool

	for i, route := range r {

		str := route.match.FindString(req.URL.Path)

		if str == "" {
			continue
		}

		found = true

		var _route = route

		if !_route.methods["ANY"] && !_route.methods[req.Method] {
			var hasOther bool

			for j, rou := range r {

				if i == j {
					continue
				}

				str = rou.match.FindString(req.URL.Path)

				if str == "" {
					continue
				}

				if rou.methods[req.Method] {
					_route = rou
					hasOther = true
				}

			}

			if !hasOther {
				w.WriteHeader(http.StatusMethodNotAllowed)
				_, _ = w.Write(methodNotAllowed)
				return
			}
		}

		_route.handlers.each(w, req, refMap)
		return
	}

	if !found {
		http.NotFound(w, req)
	}

}
