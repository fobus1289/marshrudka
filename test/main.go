package main

import (
	"fmt"
	v2 "github.com/fobus1289/marshrudka/v2"
	"log"
	"regexp"
	"strings"
)

func main() {
	type user struct {
		Id   int
		Name string
		Age  int
	}
	var server = v2.NewServer()

	server.GET(":id/:name/:age", func(request *v2.Request, object interface{}) interface{} {
		return v2.Response().Ok(201).Json([]string{"1", "@", "!"})
	})

	server.GET("/", func() interface{} {
		return v2.Response().Ok(201).Xml(&user{
			Id:   1,
			Name: "fobus",
			Age:  18,
		})
	})

	server.POST("/", func() string {
		return "POST"
	})

	server.PUT("/", func() string {
		return "PUT"
	})

	server.PATCH("/", func() string {
		return "PATCH"
	})

	server.DELETE("/", func() string {
		return "DELETE"
	})

	log.Fatalln(server.Run(":8080"))
}

func createRequestRegular(regular string) (*regexp.Regexp, error) {

	var compile, err = regexp.Compile(regular)

	if err != nil {
		panic(err)
	}

	return compile, nil
}

func getRegular(urlPath string) string {
	urlPath = strings.TrimPrefix(urlPath, "/")
	urlPath = strings.TrimSuffix(urlPath, "/")
	var regular = regexp.MustCompile(`(:[a-zA-Z]+)`)
	urlPath = regular.ReplaceAllString(urlPath, `([0-9a-zA-Z]+)`)
	return fmt.Sprintf("^(/?%s/?)$", urlPath)
}

func getPattern(urlPath string) []string {
	var regular = regexp.MustCompile(`(:[a-zA-Z]+)`)

	if result := regular.FindAllString(urlPath, -1); len(result) > 0 {
		var str = strings.Replace(strings.Join(result, " "), ":", "", -1)
		return strings.Split(str, " ")
	}

	return []string{}
}
