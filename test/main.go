package main

import (
	"embed"
	"github.com/fobus1289/marshrudka/router"
	"github.com/fobus1289/marshrudka/router/request"
	"github.com/fobus1289/marshrudka/router/response"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"sync"
)

type user struct {
	Id     int    `form:"id" json:"id" validate:"max=1,min=0"`
	Name   string `form:"name" json:"name" validate:"min=4,nonnil"`
	Age    int    `form:"age" json:"age"`
	Test   int    `form:"test" json:"test"`
	Status bool   `form:"status" json:"status"`
	NewVal int    `form:"comp" json:"new_val"`
}

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

func main() {

	//log.Println(os.Stat("static"))
	//
	//return
	var u = &user{
		Id:     1,
		Name:   "hello user",
		Age:    18,
		Test:   123,
		Status: true,
		NewVal: 222,
	}

	var (
		in = request.IModel(u)
		//out request.IModel
		//outStruct *user
	)

	var server = router.NewServer()
	server.SetService(&in)

	//server.GetService(&out)
	//log.Println(out.Validate())

	var userGroup = server.Group("user")
	{
		userGroup.POST(":id/:name", func(req request.IRequest, is []int, out *user) interface{} {
			log.Println(req.Param("name"))
			log.Println(req.HasParam("id"))
			log.Println(req.HasQuery("id"))
			log.Println(out)
			return response.Response().Ok(200).Json(is)
		}).WhereIn(map[string]string{
			"id":   `\d+`,
			"name": `(\w+)`,
		})
	}

	server.GET("/",
		func(req request.IRequest) interface{} {
			var users = make([]*user, 0, 25)
			log.Println(u)
			for i := 0; i < 25; i++ {
				users = append(users,
					&user{
						Id:     rand.Int(),
						Name:   strconv.FormatInt(rand.Int63(), 10),
						Age:    rand.Int(),
						Test:   rand.Int(),
						Status: rand.Int()%2 == 0,
						NewVal: rand.Int(),
					},
				)
			}
			return int8(122)
		},
		func(req request.IRequest) interface{} {
			return request.IModel(&user{})
		},
		func(req request.IRequest) interface{} {
			return nil
		},
	)

	//go func() {
	//	for {
	//		time.Sleep(1000 * time.Millisecond)
	//		server.HasClient()
	//	}
	//}()
	server.GET("*/dddd", func() {
		log.Println("dddd")
	}).Where("*", "ddqqq")

	server.GET("*/bb", func() {

	}).Where("*", "bd")

	server.GET("*/a", func() {

	}).Where("*", "das")

	server.GET("*/ccc", func() {
		log.Println("fas")
	}).Where("*", "fas")

	server.GET("*", http.FileServer(http.Dir("static/dist"))).
		Where("*", "static.*|asset.*").StripPrefix("static")

	log.Fatalln(server.Run(":8080"))
}
