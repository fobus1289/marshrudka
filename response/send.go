package response

import "github.com/fobus1289/marshrudka/util"

const (
	HTML = "text/html; charset=utf-8"
	JSON = "application/json; charset=utf-8"
	TEXT = "text/plain; charset=utf-8"
	XML  = "application/xml; charset=utf-8"
)

func (r *response) Json(data any) ISend {
	r.Content = JSON
	util.SerializeJson(data, &r.Body)
	return r
}

func (r *response) Xml(data any) ISend {
	r.Content = XML
	util.SerializeXml(data, &r.Body)
	return r
}

func (r *response) Text(data string) ISend {
	r.Content = TEXT
	r.Body = []byte(data)
	return r
}

func (r *response) Html(data string) ISend {
	r.Content = HTML
	r.Body = []byte(data)
	return r
}
