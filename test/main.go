package main

import (
	"github.com/fobus1289/marshrudka"
	"log"
	"net/http"
)

func cross(r *http.Request, w http.ResponseWriter) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		set := w.Header().Set
		set("Access-Control-Allow-Headers", "*")
		set("Access-Control-Allow-Methods", "*")
		set("Access-Control-Allow-Origin", "*")
	}
}

func main() {
	drive := marshrudka.NewDrive(nil)
	drive.Use(cross)

	drive.MATCH("/:id", []string{"get", "post", "put"},
		func(request *marshrudka.Request) string {
			return request.Param("id")
		},
	)

	log.Fatalln(http.ListenAndServe(":8080", drive))
}
