package main

import (
	"github.com/fobus1289/marshrudka/router"
	"log"
	"net/http"
)

func main() {

	var server = router.NewServer()

	//multipart/form-data
	//application/json
	//application/xml
	//error empty body or invalid parameter
	//status code 400
	server.BodyParseError(func() interface{} {
		return "no body"
	})

	//other errors runtime nil pointer panic ...
	//status code 500
	server.RuntimeError(func(err error) interface{} {
		return err.Error()
	})

	// server.GET | server.POST | server.PUT | server.PATCH | server.DELETE
	// return value 1 types primitive types | reference types | interface{}

	server.MATCH("other-MATCH-route", []string{"GET", "POST"}, func() {

	})

	//GET POST PUT PATCH DELETE
	server.ANY("other-ANY-route", func() {

	})

	server.GET("/", func() string {
		return "GET"
	})

	server.POST("/", func() any {
		return struct {
			Method string
		}{
			Method: "POST",
		}
	})

	server.PUT("/", func() map[string]string {
		return map[string]string{
			"Method": "PUT",
		}
	})

	server.PATCH("/", func() any {
		type route struct {
			Method string `json:"method"`
		}

		return route{
			Method: "PATCH",
		}
	})

	server.DELETE("/", func() bool {
		return true
	})

	//example file
	server.GET("*", http.FileServer(http.Dir("static/dist"))).
		Where("*", "static.*|asset.*").StripPrefix("static")

	log.Fatalln(server.Run(":8080"))
}
