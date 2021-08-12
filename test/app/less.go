package main

import (
	"github.com/fobus1289/marshrudka/router"
	"log"
	"math/rand"
)

func main() {

	var drive = router.NewRouter()

	drive.GET("/",
		//action 1
		func() int {
			return rand.Int()
		},
		//action 2
		func(randInt int) interface{} {
			log.Println(randInt + randInt)
			return randInt + randInt
		},
		//action 3
		func(randInt int) interface{} {
			// if return type is Throw stop here request
			if randInt > 5000 {
				return router.Response(400).Throw().Text("bad")
			}

			return randInt
		},
		//last action send for client
		func(randInt int) int {
			return randInt
		},
	)
	//можно делать так
	drive.GET("/a",
		func() (bool, int, string) {
			return true, 122, "hello"
		},
		func(hello string) string {

			// if ret value has is change

			if hello == "hello" {
				return hello + " world"
			}

			return hello
		},
		func(i122 int) int {
			if i122 > 500 {
				return 122
			}
			return i122
		},
		func(btrue bool) bool {

			if !btrue {
				return true
			}

			return btrue
		},
		func(hello string, i122 int, btrue bool) interface{} {
			return struct {
				Id      int
				Message string
				Ok      bool
			}{
				Id:      i122,
				Message: hello,
				Ok:      btrue,
			}
		},
	)

	var userGrop = drive.Group("user", func() int {
		//эта значение получает активный маршрут
		//if post | get | put
		return rand.Int()
	})
	{
		userGrop.GET("/", func(randInt int) {
			log.Println(randInt)
		})

		userGrop.POST("/", func(randInt int) {
			log.Println(randInt)
		})

		userGrop.PUT("/", func(randInt int) {
			log.Println(randInt)
		})

	}

	drive.Run(":8083")
	//	http.ListenAndServe(":8083", drive)
}
