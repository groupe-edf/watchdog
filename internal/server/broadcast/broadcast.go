package broadcast

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Broadcast struct {
	Message  chan []byte
	Clients  map[*websocket.Conn]bool
	Upgrader websocket.Upgrader
}

func (broadcast *Broadcast) Broadcast(message interface{}) error {
	response, err := json.Marshal(message)
	if err != nil {
		return err
	}
	broadcast.Message <- []byte(response)
	return nil
}

func (broadcast *Broadcast) Run() {
	for {
		data := <-broadcast.Message
		for client := range broadcast.Clients {
			err := client.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				client.Close()
				delete(broadcast.Clients, client)
			}
		}
	}
}

func (broadcast *Broadcast) Upgrade(w http.ResponseWriter, r *http.Request) {
	connection, err := broadcast.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	broadcast.Clients[connection] = true
}
