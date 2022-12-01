package main

import (
	"log"
	"net/http"
	"time"

	"github.com/fobus1289/marshrudka/request"
	"github.com/fobus1289/marshrudka/router"
	"github.com/fobus1289/marshrudka/validator"
)

type User struct {
	Id    int           `json:"id"`
	Name  string        `json:"name"`
	Roles []*Role       `json:"roles"`
	Role  *Role         `json:"role"`
	Exp   time.Duration `json:"exp"`
	Iat   time.Duration `json:"iat"`
}

func (u *User) Build(expired, iat time.Duration) request.IJwtUser {
	u.Exp = expired
	u.Iat = iat
	return u
}

func (u *User) Expired() time.Duration {
	return u.Exp
}

func (u *User) Out(token string) any {
	return map[string]any{
		"user":  u,
		"token": token,
	}
}

func (u *User) Validate(method string) validator.Message {

	return validator.Build(
		validator.NumberValidator("id", &u.Id).Min(1, "id < 1").Omit("omit", 0),
		validator.StringValidator("name", &u.Name).Email("email"),
		validator.ObjectValidator("role", u.Role, method).
			Options(false, true).
			Required("role req"),
		validator.ObjectsValidator("roles", u.Roles, method).
			Options(false, true).
			Min(1, "role < 1"),
	)
}

type Role struct {
	Id   int64
	Name string
}

func (r *Role) Validate(method string) validator.Message {

	return validator.Build(
		validator.NumberValidator("id", &r.Id).
			Min(1, "id lenght len > 1"),
		validator.StringValidator("name", &r.Name).
			Min(1, "name lenght len > 1").
			Equals("asda", "equas"),
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

type name struct {
	Message validator.Message
}

func main() {

	server := router.NewServer()

	server.UseAuth(&request.Jwt{
		Secret:  []byte("123456"),
		Expired: 15,
	})

	server.DeserializeError(func(err error) *router.RuntimeError {
		log.Println(err)
		return &router.RuntimeError{
			Status: 400,
		}
	})

	server.RuntimeError(func(err error) *router.RuntimeError {
		log.Println(err)
		return &router.RuntimeError{
			Status: 500,
		}
	})

	server.AddScoped(&userService{Id: 11})

	server.Use(func() map[string]any {
		return map[string]any{
			"id":   111111,
			"name": "asdas",
		}
	})

	server.UseService()

	server.GET("/",
		func(user *User) *User {
			log.Println(user)
			return user
		},
	)

	server.GET("{id}", func(u UserService) {
		log.Println(u, "{id}")
	})

	http.ListenAndServe(":8081", server)
}
