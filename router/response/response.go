package response

const (
	contentJson = "application/json; charset=utf-8"
	contentText = "text/plain"
	contentHtml = "text/html"
	contentXml  = "application/xml"
)

func Response() IResponse {
	return &response{
		status:  200,
		abort:   false,
		marshal: jsonSend,
	}
}

type response struct {
	status      int
	filePath    string
	body        interface{}
	abort       bool
	marshal     func(interface{}) []byte
	contentType string
}

func (r *response) Abort(status int) ISend {
	r.status = status
	r.abort = true
	return r
}

func (r *response) Ok(status int) ISend {
	r.status = status
	return r
}

func (r *response) Json(data interface{}) ISend {
	r.body = data
	r.contentType = contentJson
	r.marshal = jsonSend
	return r
}

func (r *response) Text(data interface{}) ISend {
	r.body = data
	r.contentType = contentText
	r.marshal = jsonSend
	return r
}

func (r *response) Html(data interface{}) ISend {
	r.body = data
	r.contentType = contentHtml
	r.marshal = jsonSend
	return r
}

func (r *response) Xml(data interface{}) ISend {
	r.body = data
	r.contentType = contentXml
	r.marshal = xmlSend
	return r
}

func (r *response) Prepare() IPrepare {
	return r
}

func (r *response) SetStatus(status int) {
	r.status = status
}

func (r *response) GetStatusCode() int {
	return r.status
}

func (r *response) SetContentType(contentType string) {
	r.contentType = contentType
}

func (r *response) GetContentType() string {
	return r.contentType
}

func (r *response) GetData() interface{} {
	return r.body
}

func (r *response) SetData(data interface{}) {
	r.body = data
}

func (r *response) IsAbort() bool {
	return r.abort
}

func (r *response) Marshal() []byte {
	if r.marshal == nil {
		r.marshal = jsonSend
	}
	return r.marshal(r.body)
}

func (r *response) GetMarshal() func(data interface{}) []byte {
	return r.marshal
}

func (r *response) SetMarshal(f func(data interface{}) []byte) {
	if f == nil {
		f = jsonSend
	}
	r.marshal = f
}

func (r *response) File() IFile {
	return &file{
		status:  r.status,
		headers: map[string]string{},
	}
}
