package main

import (
	"github.com/fobus1289/marshrudka/socket"
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

func main() {

	sokHttp := socket.NewWebSocket(&socket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true,
	})

	sokHttp.Register(&Qa{Id: 111})

	sokHttp.Default(func(client *socket.Client, data interface{}) {
		log.Println("default ", data)
	})

	type qqq struct {
		Id   int    `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}

	sokHttp.Event("b", func(conn *socket.Client, q *Qa, sok *SocketUser, clients socket.Clients, data []string) {
		log.Println(data)
		log.Println(conn.GetId())
		log.Println(q)
		log.Println(sok)
		log.Println(clients.Size())
	})

	sokHttp.Event("c", func(conn *socket.Client, data *qqq) {
		log.Println(data, "c")
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
