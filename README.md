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
