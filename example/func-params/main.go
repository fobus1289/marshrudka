package main

import (
	"fmt"
	"github.com/fobus1289/marshrudka/router"
	"log"
)

//go:generate go mod init example
func main() {

	var server = router.NewServer()

	//you can override
	//unmarshal body if body empty or invalid parameter error 400
	server.BodyParseError(func() interface{} {
		return "error message type any"
	})

	//void parameter
	//void return
	server.GET("/", func() {
		println("hello")
	})

	//header multipart/form-data
	//header application/json
	//header application/xml
	//auto unmarshal body if body empty error 400
	server.POST("/", func(user map[string]interface{}) {
		fmt.Printf("%v", user)
	})

	type User struct {
		Id   int    `json:"Id"`
		Name string `json:"Name"`
	}

	//header multipart/form-data
	//header application/json
	//header application/xml
	//auto unmarshal body if body empty error 400
	server.PUT("/", func(user User) {
		fmt.Printf("%v", user)
	})

	//header multipart/form-data
	//header application/json
	//header application/xml
	//auto unmarshal body if body empty error 400
	server.PATCH("/", func(user interface{}) {
		fmt.Printf("%v", user)
	})

	log.Fatalln(server.Run(":8080"))
}
