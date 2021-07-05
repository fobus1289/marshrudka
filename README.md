```
go get github.com/fobus1289/marshrudka
```

```go
package main

import (
	"github.com/fobus1289/marshrudka"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func UserHandler(request *marshrudka.Request) *User {
	return &User{
		Id:   1,
		Name: "user1",
		Age:  18,
	}
}

func UserHandler_v2(request *marshrudka.Request) interface{} {

	if request.Query("id") == "" {
		return marshrudka.Response(400).Throw(
			marshrudka.TEXT_PLAIN,
			"bad request",
		)
	}

	id, _ := strconv.ParseInt(request.Query("id"), 10, 32)

	return &User{
		Id:   int(id),
		Name: "user1",
		Age:  18,
	}
}

func main() {

	drive := marshrudka.NewDrive(nil)

	drive.GET("user", UserHandler)
	drive.GET("user/v2", UserHandler_v2)

	log.Fatalln(http.ListenAndServe(":8080", drive))
}

```

```go
package main

import (
	"errors"
	"math/rand"
	"github.com/fobus1289/marshrudka"
	"log"
	"net/http"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Age      int    `json:"age"`
}

type IAuthService interface {
	SignIn(username, password string) (*User, error)
	SignUp(user *User) (*User, error)
}

type AuthService struct {
}

func (a *AuthService) SignIn(username, password string) (*User, error) {

	//logic ....
	if false {
		return nil, errors.New("username or password is incorrect")
	}

	return &User{
		Id:       rand.Int() % 100,
		Username: username,
		Password: password,
	}, nil
}

func (a *AuthService) SignUp(user *User) (*User, error) {

	//logic ....
	if false {
		return nil, errors.New("username or password is incorrect")
	}

	return user, nil
}

func main() {

	drive := marshrudka.NewDrive(nil)

	var IAuth *IAuthService
	var authService = &AuthService{}

	drive.Register(IAuth, authService)
	//or
	drive.Register(authService)

	//other services can be similarly registered
	drive.Register(here)

	drive.GET("/",
		func(auth IAuthService, r *http.Request) interface{} {
			username := r.URL.Query().Get("username")
			password := r.URL.Query().Get("password")

			user, err := auth.SignIn(username, password)

			if err != nil {
				return marshrudka.
					Response(400).
					Throw(marshrudka.TEXT_PLAIN, err.Error())
			}

			return marshrudka.Response(200).Json(user)
		},
	)

	log.Fatalln(http.ListenAndServe(":8080", drive))
}


```