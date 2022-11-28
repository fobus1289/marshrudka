package main

import (
	"log"
	"net/http"

	"github.com/fobus1289/marshrudka/router"
	"github.com/fobus1289/marshrudka/validator"
)

type User struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Roles []*Role `json:"roles"`
}

func (u *User) Validate() validator.MessageMapResult {
	return nil
	return validator.Build(
		validator.NumberValidator("id", u.Id).Min(1, "id can be 0"),
		validator.StringValidator("name", u.Name).Min(4, "name lenght len < 4"),
		validator.ObjectsValidator("roles", u.Roles, true),
	)
}

type Role struct {
	Id   int64
	Name string
}

func (r *Role) Validate() validator.MessageMapResult {
	return validator.Build(
		validator.NumberValidator("id", r.Id).Min(1, "id can be 0"),
		validator.StringValidator("name", r.Name).Min(4, "name lenght len < 4"),
	)
}

type UserService interface {
	Get() string
}

type userService struct {
	Id int
}

func (us *userService) Get() string {
	return "hello user service"
}

func asd() UserService {
	return nil
}

func main() {

	server := router.NewServer()

	server.DeserializeError(func(err error) *router.RuntimeError {
		log.Println(err)
		return &router.RuntimeError{
			Status: 400,
		}
	})

	//single
	//scope
	userSerice := userService{Id: 11}

	server.AddScoped(func() UserService {
		v := UserService(&userSerice)
		return v
	}())

	server.Use(func() map[string]any {
		return map[string]any{
			"id":   111111,
			"name": "asdas",
		}
	})

	server.UseService()

	server.GET("/",
		func(u *userService) int {
			u.Id = 2222
			log.Println(u)
			return 2143421
		},
		func(u UserService, user *User, m map[string]any) int {
			log.Println(u, "2")
			log.Println(user)
			log.Println(userSerice)
			log.Println(m)
			return 121
		},
	)

	http.ListenAndServe(":8081", server)
}
