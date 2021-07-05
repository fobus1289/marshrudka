package main

import (
	"marshrudka"
	"net/http"
)

func main() {
	drive := marshrudka.NewDrive(nil)

	drive.ANY("",
		func(request *marshrudka.Request) struct {
			Id   int
			Name string
		} {

			return struct {
				Id   int
				Name string
			}{1, "fobus"}
		},
	)

	http.ListenAndServe(":8080", drive)
}
