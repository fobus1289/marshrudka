package http_gin

func (s *Serv) GET(name string, handler ...interface{}) {
	s.g.GET(name, s.parseFunc(handler...)...)
}

func (s *Serv) POST(name string, handler ...interface{}) {
	s.g.POST(name, s.parseFunc(handler...)...)
}

func (s *Serv) PUT(name string, handler ...interface{}) {
	s.g.PUT(name, s.parseFunc(handler...)...)
}

func (s *Serv) PATCH(name string, handler ...interface{}) {
	s.g.PATCH(name, s.parseFunc(handler...)...)
}

func (s *Serv) ANY(name string, handler ...interface{}) {
	s.g.Any(name, s.parseFunc(handler...)...)
}

func (s *Serv) DELETE(name string, handler ...interface{}) {
	s.g.DELETE(name, s.parseFunc(handler...)...)
}
