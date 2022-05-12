package main

import (
	"github.com/fobus1289/marshrudka/router"
	"github.com/fobus1289/marshrudka/router/request"
	"github.com/fobus1289/marshrudka/router/response"
	"log"
	"net/http"
)

func main() {

	var server = router.NewServer()

	server.POST("user/photo", func(formFile request.IFormFile) interface{} {

		var filepath string

		err := formFile.Get("photo").
			GetFirst().
			RandomFileName().
			Store("static/user", &filepath).Error()

		if err != nil {
			return response.Response().Abort(http.StatusBadRequest).
				Json(
					map[string]string{
						"error": err.Error(),
					},
				)
		}

		return response.Response().Ok(http.StatusOK).Json(
			map[string]string{
				"message": "file save success",
				"path":    filepath,
			},
		)
		
	})

	log.Fatalln(server.Run(":8080"))
}
