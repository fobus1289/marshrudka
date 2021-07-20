package main

import (
	"github.com/fobus1289/marshrudka/socket"
	"log"
	"net/http"
)

func sok() {
	webSocket := socket.NewWebSocket(&socket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		WriteBufferPool: nil,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true,
	})

	webSocket.Default(Default)
	webSocket.Connection(Connection)
	webSocket.Disconnection(Disconnection)

	webSocket.Event("move", move)
	webSocket.Event("me", me)
	webSocket.Event("live", live)

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		webSocket.NewClient(writer, request, nil)
	})

	http.ListenAndServe(":8080", nil)
}

func Default(client *socket.Client, data interface{}) {
	//	println(client.GetId(), " move")
}

func Connection(currentClient *socket.Client, r *http.Request) {

	player := &Player{
		Id:       currentClient.GetId(),
		Position: &Vector{},
	}

	currentClient.BroadcastClients("conn", player)

	clients := currentClient.Clients().Filter(func(client *socket.Client) bool {
		return currentClient != client
	})

	if len(clients) > 0 {
		var players []interface{}
		clients.ForEach(func(client *socket.Client) {
			players = append(players, client.GetOwner())
		})
		currentClient.Emit("users", map[string]interface{}{
			"players": players,
		})
	}

	currentClient.SetOwner(player)

	currentClient.Emit("me", player)
}

func Disconnection(client *socket.Client) {

}

type Player struct {
	Id       int64   `json:"id"`
	Position *Vector `json:"position"`
}

type Vector struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

func me(client *socket.Client) {
	client.BroadcastClients("me", client.GetId())
}

func live(client *socket.Client) {
	client.BroadcastClients("live", client.GetId())
}

func move(client *socket.Client, player *Player) {
	log.Println(player)
	client.BroadcastClients("move", player)
}
