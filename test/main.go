package main

import (
	"log"
	"net/http"
	"reflect"
)

func cross(r *http.Request, w http.ResponseWriter) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		set := w.Header().Set
		set("Access-Control-Allow-Headers", "*")
		set("Access-Control-Allow-Methods", "*")
		set("Access-Control-Allow-Origin", "*")
	}
}

func main() {

	userservice := &UserService{
		Id:   22,
		Data: []byte("asdasdsadsa"),
	}

	//userservice2 := userservice
	//userservice2.Id = 1212
	//log.Println(userservice2)
	//log.Println(userservice)

	value := reflect.ValueOf(userservice)
	vale2 := value.Elem().Interface().(UserService)
	vale2.Id = 1111
	log.Println(vale2)
	log.Println(value.Interface())

	//log.Println(reflect.Copy(value1, value))

	//drive := marshrudka.NewDrive(nil)
	//drive.Use(cross)
	//
	//drive.MATCH("/:id", []string{"get", "post", "put"},
	//	func(request *marshrudka.Request) string {
	//		return request.Param("id")
	//	},
	//)
	//
	//log.Fatalln(http.ListenAndServe(":8080", drive))
}

type UserService struct {
	Id   int
	Data []byte
}
