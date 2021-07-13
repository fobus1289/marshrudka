package socket

func (w *WebSocket) Clients() Clients {
	return w.clients
}

func (c Clients) Size() int {
	return len(c)
}

func (c Clients) ForEach(fn func(client *Client)) {
	for _, client := range c {
		fn(client)
	}
}

func (c Clients) Find(fn func(client *Client) bool) *Client {
	for _, client := range c {
		if fn(client) {
			return client
		}
	}
	return nil
}

func (c Clients) Filter(fn func(client *Client) bool) Clients {
	var clients Clients

	for _, client := range c {
		if fn(client) {
			clients = append(clients, client)
		}
	}

	return clients
}
