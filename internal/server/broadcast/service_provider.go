package broadcast

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/server/container"
)

const (
	ServiceName = "broadcast"
)

type ServiceProvider struct {
	Options *config.Database
}

func (service *ServiceProvider) Register(di container.Container) {
	di.Set(ServiceName, func(_ container.Container) container.Service {
		var clients = make(map[*websocket.Conn]bool)
		var upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		return &Broadcast{
			Message:  make(chan []byte),
			Clients:  clients,
			Upgrader: upgrader,
		}
	})
}
