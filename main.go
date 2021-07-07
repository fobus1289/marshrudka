package marshrudka

import (
	"log"
	"os"
	"reflect"
	"strings"
)

type services map[reflect.Type]reflect.Value

type Drive struct {
	routers
	services
	Logger   *log.Logger
	handlers []interface{}
}

func NewDrive(logger *log.Logger) *Drive {

	if logger == nil {
		logger = log.New(os.Stdout, "Info: ", log.Ltime|log.Lshortfile)
	}

	return &Drive{
		Logger:   logger,
		services: map[reflect.Type]reflect.Value{},
	}
}

func (d *Drive) Group(path string, handlers ...interface{}) *group {

	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	return &group{
		Path:    path,
		actions: handlers,
		Drive:   d,
	}
}
