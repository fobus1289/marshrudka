```
go get github.com/fobus1289/marshrudka
```

```go
//example 1
package main

import (
	"github.com/fobus1289/marshrudka/router"
	"log"
	"math/rand"
)

type User struct {
	Id    int64  `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type Users []*User

var ResponseUsers = Users{
	&User{
		Id:    1,
		Login: "login_1",
		Name:  "Jhone",
		Age:   18,
	},
	&User{
		Id:    2,
		Login: "login_2",
		Name:  "Doe",
		Age:   35,
	},
	&User{
		Id:    3,
		Login: "login_3",
		Name:  "Bob",
		Age:   27,
	},
}

/*
	Параметры могут быть пустыми
	Возвращаемый тип может быть любым
	interface{} итд, если это ссылочный тип, то он преобразуется в json
*/
func GetAllUsers() Users {
	log.Println("get all users")
	return ResponseUsers
}

/*
	если не известный параметр в функции и он ссылочный типа то он конвертироваеца в тип в параметре
	преметивы не допускаются
*/
func CreateUser(user *User) interface{} {

	for _, responseUser := range ResponseUsers {
		if responseUser.Login == user.Login {
			return router.Response(400).
				Throw().Json(
				map[string]string{
					"message": "user login is duplicate",
				},
			)
		}
	}

	user.Id = rand.Int63() % 100

	return user
}

/*
*router.Request это служебная структура роутера
 */
func DeleteUser(request *router.Request) interface{} {
	//var id = request.Query("id")
	//var id = request.QueryGetInt("id")

	var id int64

	if !request.TryQueryGetInt("id", &id) {
		return router.Response(400).
			Throw().Json(
			map[string]string{
				"message": "id can be empty or string",
			},
		)
	}

	for i, user := range ResponseUsers {
		if user.Id == id {

			ResponseUsers = append(ResponseUsers[:i], ResponseUsers[i+1:]...)

			return true
			//return 1
			//return "user deleted" + strconv.FormatInt(id, 10)
			//return user
		}
	}

	return router.Response(400).
		Throw().Json(
		map[string]string{
			"message": "user not found",
		},
	)
}

/*
	Динамическая реализация параметров таких как
	http.ResponseWriter
	*http.Request
 	*router.Request
*/
func GetUserById(request *router.Request) interface{} {
	var id = request.QueryGetInt("id")

	for _, user := range ResponseUsers {
		if user.Id == id {
			return router.Response(200).Json(user)
		}
	}

	return router.Response(404).Throw().Json(
		map[string]string{
			"message": "user not found",
		},
	)

}

func main() {
	var drive = router.NewRouter()

	drive.GET("/", GetAllUsers)
	drive.POST("/", CreateUser)
	drive.DELETE("/", DeleteUser)
	drive.GET("/by-id", GetUserById)

	drive.Run(":8081")
}

```

```go
//example 2
package main

import (
	"errors"
	"github.com/fobus1289/marshrudka/router"
	"log"
	"os"
)

//models
type UserModel struct {
	Id     int64  `json:"id"`
	Login  string `json:"login"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type UserModels []*UserModel

//services
type IUserService interface {
	All() (UserModels, error)
	Create(userModel *UserModel) error
	Delete(id int64) error
}

type UserService struct {
}

func (u *UserService) All() (UserModels, error) {

	if false {
		return nil, errors.New("bad")
	}

	return UserModels{}, nil
}

func (u *UserService) Create(userModel *UserModel) error {

	if false {
		return errors.New("bad")
	}

	return nil
}

func (u *UserService) Delete(id int64) error {

	if false {
		return errors.New("bad")
	}

	return nil
}

//controllers
type UserController struct {
	*log.Logger
}

func (c *UserController) GetAll(userService IUserService) interface{} {

	var userModels, err = userService.All()

	if err != nil {
		return router.Response(400).Throw().Json(map[string]string{
			"message": err.Error(),
		})
	}

	return userModels
}

func (c *UserController) Create(userModel *UserModel, userService IUserService) interface{} {

	if err := userService.Create(userModel); err != nil {
		return router.Response(400).Throw().Json(map[string]string{
			"message": err.Error(),
		})
	}

	return userModel
}

func (c *UserController) Delete(request *router.Request, userService *UserService) interface{} {

	var id int64

	if !request.TryParamGetInt("id", &id) {
		return router.Response(400).Throw().Json(map[string]string{
			"message": "id can be empty or string",
		})
	}

	if err := userService.Delete(id); err != nil {
		return router.Response(400).Throw().Json(map[string]string{
			"message": err.Error(),
		})
	}

	return id
	//return true
	//return "user deleted" + strconv.FormatInt(id, 10)
}

func main() {

	var drive = router.NewRouter()

	//init services
	var iUserService *IUserService

	var userService = &UserService{}

	var loggerService = log.New(os.Stdout, "INFO:", log.LstdFlags|log.Lshortfile)

	// вот регистрация услуг, которые необходимы для дальнейшего использования в параметрах контроллера
	drive.Register(loggerService)
	drive.Register(userService)
	//можно указать с интефесом или без
	drive.Register(iUserService, userService)

	//init controllers
	var userController = &UserController{}

	//dependency injection = Dep
	//type UserController struct {
	//	*log.Logger 				   = Dep
	//  UserController *UserController = Dep
	//  *UserController 			   = Dep
	//}
	//если такой такая ссылка была зарегистрована то он будет внедрон
	//только ссылочный тип структура
	drive.Dep(userController)

	var userGroup = drive.Group("user")
	{
		userGroup.GET("/", userController.GetAll)
		userGroup.POST("/", userController.Create)
		//итог будет такой
		//localhost:8082/int  valid :id{n} number | :id{s} string | :id * all
		//localhost:8082/string invalid
		userGroup.DELETE("/:id{n}", userController.Delete)
		//userGroup.DELETE("/:id{s}", userController.Delete)
		//userGroup.DELETE("/:id", userController.Delete)
	}

	drive.Run(":8082")
}


```