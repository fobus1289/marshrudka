package main

import (
	"embed"
	"github.com/fobus1289/marshrudka/router"
	"github.com/fobus1289/marshrudka/router/request"
	"log"
	"net/http"
	"reflect"
	"sync"
)

func (u *user) Validate() bool {
	return false
}

type Xml struct {
	XMLName string      `xml:"users"`
	Data    interface{} `json:"user"`
}

func Set(v interface{}) bool {
	var refValue = reflect.ValueOf(v)

	if refValue.Kind() != reflect.Ptr {
		return false
	}

	var (
		key   = refValue.Elem()
		value = key
	)

	if !key.IsValid() {
		return false
	}

	if key.Kind() == reflect.Interface {
		if value = key.Elem(); value.Kind() == reflect.Invalid {
			return false
		}

		log.Println(key.Type())
		log.Println(value.Type())
		return true
	}

	log.Println(refValue.Type())
	log.Println(value.Type())

	return true
}

type Sht struct {
	Counter int
	*sync.Mutex
}

type name interface {
}

type asd struct {
	Id   int    `json:"id" form:"id" validate:""`
	Name string `form:"name" json:"name" validate:"min=4,nonnil"`
}

func (a *asd) Validate() bool {
	return false
}

//go:embed static/dist/*
var frontend embed.FS

type qw struct {
	A *user
}

type user struct {
	Id     int    `form:"id" json:"id" validate:"max=1,min=0"`
	Name   string `form:"name" json:"name" validate:"min=4,nonnil"`
	Age    int    `form:"age" json:"age"`
	Test   int    `form:"test" json:"test"`
	Status bool   `form:"status" json:"status"`
	NewVal int    `form:"comp" json:"new_val"`
	User   request.IModel
}

var ch = make(chan func(), 50)

func a() {
	for f1 := range ch {
		func(c chan func()) {
			select {
			case f := <-ch:
				f()
			}
		}(ch)
		_ = f1
	}
}

func main() {

	//var strs = []string{"1", "2", "3", "4"}
	//log.Println(strs[1:])
	//return
	var server = router.NewServer()

	server.BodyParseError(func() interface{} {
		return "no body"
	})

	server.RuntimeError(func(err error) interface{} {
		return err.Error()
	})

	server.GET("/", func(b [10]asd, param request.IRequest) {
		log.Println(b)
	})

	type Us struct {
		Id   int     `json:"id,omitempty"`
		Name string  `json:"name,omitempty"`
		Age  float32 `json:"age,omitempty"`
	}

	server.POST("/:id", func(files request.IFormFile) {
		//var paths []string
		for _, file := range files.Files().Get("file").Files() {
			log.Println(file.Info().ContentType())
		}
	})

	server.PUT("/", func() {
		log.Println("PUT")
	})

	server.PATCH("/", func() {
		log.Println("PATCH")
	})

	server.DELETE("/", func() {
		log.Println("DELETE")
	})

	server.GET("*", http.FileServer(http.Dir("static/dist"))).
		Where("*", "static.*|asset.*").StripPrefix("static")

	log.Fatalln(server.Run(":8080"))
}
