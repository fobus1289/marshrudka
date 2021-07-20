package router

func (d *Drive) Use(handler ...interface{}) *Drive {
	d.handlers = d.parseFunc(true, handler...)
	return d
}
