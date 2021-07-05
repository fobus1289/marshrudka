package marshrudka

import (
	"log"
	"os"
	"reflect"
	"strings"
)

type services map[reflect.Type]reflect.Value

type drive struct {
	routers
	services
	Logger   *log.Logger
	handlers []interface{}
}

func NewDrive(logger *log.Logger) *drive {

	if logger == nil {
		logger = log.New(os.Stdout, "Info: ", log.Ltime|log.Lshortfile)
	}

	return &drive{
		Logger:   logger,
		services: map[reflect.Type]reflect.Value{},
	}
}

func (d *drive) Group(path string, handlers ...interface{}) *group {

	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	return &group{
		Path:    path,
		actions: handlers,
		drive:   d,
	}
}
