package socket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	Conn   *websocket.Conn
	socket *WebSocket
	id     int64
	owner  interface{}
}

type Clients []*Client

var (
	typeClients = reflect.TypeOf(Clients{})
	typeClient  = reflect.TypeOf(&Client{})
)

func (c *Clients) Delete(client *Client) bool {

	tmp := *c

	for i, c2 := range tmp {
		if c2 == client {
			*c = append(tmp[:i], tmp[i+1:]...)
			return true
		}
	}

	return false
}

func (c *Client) SetOwner(owner interface{}) {
	c.owner = owner
}

func (c *Client) GetOwner() interface{} {
	return c.owner
}

func (c *Client) SetId(id int64) {
	c.id = id
}

func (c *Client) GetId() int64 {
	return c.id
}

func (w *WebSocket) NewClient(res http.ResponseWriter, req *http.Request, header http.Header) (*Client, error, ) {

	conn, err := w.upgrader.Upgrade(res, req, header)

	if err != nil {
		return nil, err
	}

	var client = &Client{
		Conn:   conn,
		socket: w,
		id:     time.Now().UnixNano(),
	}

	w.clients = append(w.clients, client)

	go client.Read()

	return client, nil
}

func (c *Client) WriteBytes(data []byte) {
	if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
		c.socket.clients.Delete(c)
	}
}

func (c *Client) WriteString(data string) {
	if err := c.Conn.WriteMessage(websocket.TextMessage, []byte(data)); err != nil {
		c.socket.clients.Delete(c)
	}
}

func (c *Client) Read() {

	defer func() {
		c.socket.clients.Delete(c)
		_ = c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		if data, ok := c.socket.hasEvent(c, message); !ok {
			c.socket._default(c, data)
		}

	}

}

func (w *WebSocket) hasEvent(c *Client, msg []byte) (string, bool) {

	str := string(msg)

	index := strings.Index(str, ":")

	if index == -1 {
		return str, false
	}

	event := str[:index]
	fu := w.mapper[event]

	if fu == nil {
		return str, false
	}

	data := str[index+1:]

	w.parseFunc(fu, c, data)
	return data, true
}

func (w *WebSocket) parseFunc(fu interface{}, c *Client, msg string) {

	fuValue := reflect.ValueOf(fu)

	if fuValue.Kind() != reflect.Func {
		log.Println("dont supported:", fuValue.Kind())
		return
	}

	fuType := fuValue.Type()

	var params []reflect.Value

	owner := c.owner

	for i := 0; i < fuType.NumIn(); i++ {
		inType := fuType.In(i)

		if owner != nil {
			ownerValue := reflect.ValueOf(owner)
			if reflect.DeepEqual(inType, ownerValue.Type()) {
				params = append(params, ownerValue)
				continue
			}
		}

		if reflect.DeepEqual(inType, typeClient) {
			params = append(params, reflect.ValueOf(c))
			continue
		}

		if reflect.DeepEqual(inType, typeClients) {
			params = append(params, reflect.ValueOf(w.clients))
			continue
		}

		service := w.services[inType]

		if service.Kind() != reflect.Invalid {
			params = append(params, service)
			continue
		}

		value := getParamValue(inType, msg)
		params = append(params, value)
	}

	fuValue.Call(params)
}

func getParamValue(value reflect.Type, data string) reflect.Value {

	var kind = value.Kind()

	if value.Kind() == reflect.Ptr {
		kind = value.Elem().Kind()
	}

	switch kind {
	case reflect.Bool:
		out, _ := strconv.ParseBool(data)
		return reflect.ValueOf(out)
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		out, _ := strconv.ParseInt(data, 10, 64)
		return reflect.ValueOf(out)
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		out, _ := strconv.ParseUint(data, 10, 64)
		return reflect.ValueOf(out)
	case reflect.Float32,
		reflect.Float64:
		out, _ := strconv.ParseFloat(data, 10)
		return reflect.ValueOf(out)
	case reflect.String:
		return reflect.ValueOf(data)
	case reflect.Struct, reflect.Slice, reflect.Interface, reflect.Map:
		return setOther(value, data)
	}

	return reflect.ValueOf([]byte(data))
}

func setOther(_type reflect.Type, data string) reflect.Value {

	if _type.Kind() == reflect.Slice {
		outerObj := reflect.New(reflect.SliceOf(_type.Elem()))
		outerObj.Elem().Set(reflect.MakeSlice(reflect.SliceOf(_type.Elem()), 0, 0))
		_ = json.Unmarshal([]byte(data), outerObj.Interface())
		return outerObj.Elem()
	}

	_value := reflect.New(_type)

	_ = json.Unmarshal([]byte(data), _value.Interface())

	return _value.Elem()
}