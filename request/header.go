package request

func (r *request) GetHeader(key string) string {
	return r.Request.Header.Get(key)
}

func (r *request) Authorization() string {
	return r.GetHeader("Authorization")
}
