package main

//go:generate go mod init example
import (
	"github.com/fobus1289/marshrudka/router"
	"github.com/fobus1289/marshrudka/router/request"
	"log"
	"math/rand"
)

//model
type User struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Salary float64 `json:"salary"`
}

//controllers
type UserController struct {
	UserService *UserService
}

//action
func (u UserController) Get(param request.IParam) User {
	id := param.Param("id").CastInt(0)

	return u.UserService.GetUserById(int(id))
}

//action
func (u UserController) Create(user User) User {
	return u.UserService.CreateUser(user)
}

//services
type UserService struct {
}

func (u UserService) GetUserById(id int) User {
	return User{
		Id:     id,
		Name:   "user 1",
		Age:    18,
		Salary: 100.00,
	}
}

func (u UserService) CreateUser(user User) User {
	user.Id = rand.Int()
	return user
}

func main() {

	var server = router.NewServer()
	//first one
	userService := new(UserService)
	userController := new(UserController)

	server.SetService(userService)
	server.FillServiceFields(userController)

	userGroup := server.Group("user")
	{
		userGroup.GET("/", userController.Get)
		userGroup.GET("/", userController.Create)
	}

	log.Fatalln(server.Run(":8080"))
}
