package response

func (r *response) HasAbort() bool {
	return r.IsAbort
}

func (r *response) HasBody() bool {
	return r.Body != nil
}

func (r *response) ContentType() string {
	return r.Content
}

func (r *response) GetBody() []byte {
	return r.Body
}

func (r *response) GetStatus() int {
	return r.Status
}
