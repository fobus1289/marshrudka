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
