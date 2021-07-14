package socket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"reflect"
	"sync"
	"time"
)

type Config struct {
	HandshakeTimeout                time.Duration
	ReadBufferSize, WriteBufferSize int
	WriteBufferPool                 websocket.BufferPool
	Subprotocols                    []string
	Error                           func(w http.ResponseWriter, r *http.Request, status int, reason error)
	CheckOrigin                     func(r *http.Request) bool
	EnableCompression               bool
}

//type _map map[string]func(client *Client, data interface{})
type _map map[string]interface{}
type Channels map[string]Clients

type WebSocket struct {
	mux            *sync.RWMutex
	upgrader       *websocket.Upgrader
	services       map[reflect.Type]reflect.Value
	clients        Clients
	channels       Channels
	_default       func(client *Client, data interface{})
	_connection    func(client *Client, r *http.Request)
	_disconnection func(client *Client)
	mapper         _map
}

func (w *WebSocket) Default(event func(client *Client, data interface{})) {
	w._default = event
}

func (w *WebSocket) Event(name string, event interface{}) {
	w.mapper[name] = event
}

func (w *WebSocket) Connection(fn func(client *Client, r *http.Request)) {
	w._connection = fn
}

func (w *WebSocket) Disconnection(fn func(client *Client)) {
	w._disconnection = fn
}

func (w *WebSocket) Register(_interface interface{}, _struct ...interface{}) *WebSocket {

	if _struct == nil {
		_structValue := reflect.ValueOf(_interface)
		_structElemet := _structValue.Elem()
		w.services[_structValue.Type()] = _structValue
		w.services[_structElemet.Type()] = _structElemet
		return w
	}

	if len(_struct) != 1 {
		log.Fatalln("something went wrong 1")
	}

	if implement(_interface, _struct[0]) {
		_interfaceType := reflect.TypeOf(_interface)
		_structValue := reflect.ValueOf(_struct[0])
		w.services[_interfaceType.Elem()] = _structValue
		w.services[_structValue.Type()] = _structValue
	} else {
		log.Fatalln("something went wrong 2")
	}

	return w
}

func (w *WebSocket) Dep(owner interface{}) {

	ownerType := reflect.ValueOf(owner)

	if ownerType.Kind() == reflect.Ptr {
		ownerType = ownerType.Elem()
	}

	for i := 0; i < ownerType.NumField(); i++ {
		fieldType := ownerType.Field(i)

		service := w.services[fieldType.Type()]

		if service.Kind() != reflect.Invalid {
			fieldType.Set(service)
		}
	}

}

func implement(_interface, _struct interface{}) bool {

	structType := reflect.TypeOf(_struct)
	{
		if structType.Kind() != reflect.Ptr {
			log.Fatalln("ffs 1")
		}
	}

	interfaceType := reflect.TypeOf(_interface)
	{
		if interfaceType.Kind() != reflect.Ptr {
			log.Fatalln("ffs 2")
		}
	}

	if interfaceType.Elem().Kind() == reflect.Struct {
		return structType.AssignableTo(interfaceType)
	}

	return structType.AssignableTo(interfaceType.Elem())
}

func NewWebSocket(config *Config) *WebSocket {

	var upgrader *websocket.Upgrader

	if config != nil {
		upgrader = &websocket.Upgrader{
			HandshakeTimeout:  config.HandshakeTimeout,
			ReadBufferSize:    config.ReadBufferSize,
			WriteBufferSize:   config.WriteBufferSize,
			WriteBufferPool:   config.WriteBufferPool,
			Subprotocols:      config.Subprotocols,
			Error:             config.Error,
			CheckOrigin:       config.CheckOrigin,
			EnableCompression: config.EnableCompression,
		}
	} else {
		upgrader = &websocket.Upgrader{}
	}

	return &WebSocket{
		mux:            &sync.RWMutex{},
		upgrader:       upgrader,
		_connection:    func(client *Client, r *http.Request) {},
		_disconnection: func(client *Client) {},
		services:       map[reflect.Type]reflect.Value{},
		channels:       Channels{},
		_default:       func(client *Client, data interface{}) {},
		mapper:         _map{},
	}
}
