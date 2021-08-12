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
