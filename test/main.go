package main

import (
	"github.com/fobus1289/marshrudka/router"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func CORSEnabledFunction(w http.ResponseWriter, r *http.Request) interface{} {
	// Set CORS headers for the preflight request
	///*
	//	//ctx:= r.WithContext(context.Background())
	//	hj, ok := w.(http.Hijacker)
	//
	//	if ok {
	//
	//		conn, _, _ := hj.Hijack()
	//
	//		conn.Close()
	//
	//		return router.Throw{
	//			StatusCode: 204,
	//		}
	//	}
	//*/
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return router.Throw{
			StatusCode: 204,
		}
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	return nil
}

func main() {

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"*"},
		AllowedHeaders: []string{"*"},
	})
	//c := cors.New(cors.Options{
	//	AllowedOrigins: []string{"http://foo.com", "http://foo.com:8080"},
	//	AllowCredentials: true,
	//	// Enable Debugging for testing, consider disabling in production
	//	Debug: true,
	//})
	route := router.NewRouter()
	_ = c

	route.Use(CORSEnabledFunction)

	//route.Use(func(w http.ResponseWriter, r *http.Request) {
	//	log.Println(r.Method)
	//	w.Header().Set("Accept", "*/*")
	//	w.Header().Set("Cache-Control", "no-cache")
	//	w.Header().Set("Accept-Encoding", "gzip, deflate, br")
	//	w.Header().Set("Access-Control-Allow-Origin", "*")
	//	w.Header().Set("Access-Control-Allow-Methods", "POST")
	//	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//	w.Header().Set("Access-Control-Max-Age", "3600")
	//	//w.WriteHeader(http.StatusNoContent)
	//})

	//route.Use(CORSEnabledFunction)

	route.MATCH("/", []string{"get", "post", "put"}, func() {

	})

	adminGroup := route.Group("/admin")
	{
		adminGroup.GET("/user", func() {
			println("admin/user")
		})

		adminCompanyGroup := adminGroup.Group("/company")
		{
			adminCompanyGroup.GET("/user/", func() {
				println("admin/company/user get")
			})
			adminCompanyGroup.POST("/user/", func() {
				println("admin/company/user post")
			})
		}

	}

	route.ANY("/any/:id{n}", func(request *router.Request) interface{} {

		//log.Println(request.Request.ParseMultipartForm(0))
		//log.Println(request.Request.MultipartForm)

		file, err := request.FormFile("file")

		if err != nil {
			return router.Response(400).Throw().Text("bad request")
		}

		err = file.SetName("rd").Store("static")

		log.Println(err)
		return router.Response(200).File("static/test.txt", "ok.txt").Download()
	})

	route.Run(":8080")
}
