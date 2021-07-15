package main

import (
	http_gin "github.com/fobus1289/marshrudka/http-gin"
	"github.com/fobus1289/marshrudka/socket"
	"github.com/gin-contrib/cors"
	"log"
	"net/http"
	"reflect"
	"strconv"
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

type SocketUser struct {
	Id   int64
	Name string
}

type Qa struct {
	Id int
}

type Ba struct {
	Id int
}

type Da struct {
	Id int
}

type AAA interface {
	Name()
}

func (d *Da) Name() {
	println("name da")
}

func main() {

	dr := http_gin.NewDrive()

	dr.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	dr.Register(&Da{12})
	dr.Register(&Ba{23232})
	dr.Register((*AAA)(nil), &Da{1222})

	type qaa struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
		*Da
	}

	group := dr.Group("asd")
	{
		group.GET("/", func() {
			log.Println("hello")
		})

		group.POST("/", func() {
			log.Println("hello")
		})

		child := group.Group("a")
		{
			child.GET("qq", func() interface{} {
				return &http_gin.Response{
					StatusCode: 200,
					Data: map[string]interface{}{
						"id":   1,
						"name": "fobus",
						"age":  99,
					},
				}
			})

			child1 := child.Group("as", func() { println("as start") })
			{
				child1.GET("xua", func() {
					println("xua start")
				})
			}
		}
	}

	dr.POST("/",
		func(r *http.Request, d AAA) interface{} {
			log.Println()
			d.Name()
			return 1222
		},
		func(r *http.Request, i int, d *Da, ba *Ba, q map[string]interface{}) interface{} {
			log.Println(i)
			return &http_gin.Response{
				StatusCode: 200,
				Data:       q,
			}
		},
	)

	dr.Run("localhost:8080")

	return

	sokHttp := socket.NewWebSocket(&socket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true,
	})

	type Use struct {
		*Qa
		*Ba
		AAA
	}

	sokHttp.Register(&Qa{Id: 1})
	//sokHttp.Register((*AAA)(nil), &Da{Id: 1})
	sokHttp.Register(&Ba{Id: 2})
	sokHttp.Register(&Da{Id: 3})

	use := &Use{}

	sokHttp.Dep(use)

	log.Println(use.Qa.Id)
	log.Println(use.Ba.Id)
	log.Println(use.AAA)

	return
	sokHttp.Default(func(client *socket.Client, data interface{}) {
		log.Println("default ", data)
	})

	type qqq struct {
		Id   int    `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}

	sokHttp.Event("b", func(conn *socket.Client, q *Qa, sok *SocketUser, clients socket.Clients, data []string) interface{} {
		return map[string]int{
			"a": 111,
			"b": 222,
			"c": 333,
		}
	})

	sokHttp.Event("c", func(conn *socket.Client, data *qqq) interface{} {
		log.Println(data, "c")

		return data
	})

	var id int64

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		id++
		client, err := sokHttp.NewClient(writer, request, nil)
		if err != nil {
			return
		}

		client.SetOwner(&SocketUser{
			Id:   id,
			Name: "asdsadsads" + strconv.Itoa(int(id)),
		})

		client.SetId(id)
	})

	log.Fatalln(http.ListenAndServe(":8000", nil))

	return
	userservice := &UserService{
		Id:   22,
		Data: []byte("asdasdsadsa"),
	}

	//userservice2 := userservice
	//userservice2.Id = 1212
	//log.Println(userservice2)
	//log.Println(userservice)

	value := reflect.ValueOf(userservice)
	vale2 := value.Elem().Interface().(UserService)
	vale2.Id = 1111
	log.Println(vale2)
	log.Println(value.Interface())

	//log.Println(reflect.Copy(value1, value))

	//drive := marshrudka.NewDrive(nil)
	//drive.Use(cross)
	//
	//drive.MATCH("/:id", []string{"get", "post", "put"},
	//	func(request *marshrudka.Request) string {
	//		return request.Param("id")
	//	},
	//)
	//
	//log.Fatalln(http.ListenAndServe(":8080", drive))
}

type UserService struct {
	Id   int
	Data []byte
}
