package response

type IResponse interface {
	Send(status int) ISend
	Abort(status int) ISend
}

type ISend interface {
	Json(data any) ISend
	Text(data string) ISend
	Html(data string) ISend
	Xml(data any) ISend
}

type IReady interface {
	HasAbort() bool
	HasBody() bool
	ContentType() string
	GetBody() []byte
	GetStatus() int
}

type response struct {
	IsAbort bool
	Body    []byte
	Content string
	Status  int
}

func (r *response) Send(status int) ISend {
	r.Status = status
	return r
}

func (r *response) Abort(status int) ISend {
	r.Status = status
	r.IsAbort = true
	return r
}
