package pool

import (
	"github.com/gorilla/websocket"
)

type Pool struct {
	Broadcast chan []byte
	Clients   map[*websocket.Conn]bool
}

func (pool *Pool) Run() {
	for {
		data := <-pool.Broadcast
		for client := range pool.Clients {
			err := client.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				client.Close()
				delete(pool.Clients, client)
			}
		}
	}
}
