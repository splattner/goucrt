package integration

import (
	"log"

	"github.com/gorilla/websocket"
)

type remote struct {
	standby   bool
	connected bool
	websocket *websocket.Conn
}

func (r *remote) EnterStandBy() {
	log.Println("Remote entered standby mode")

	r.standby = true

}

func (r *remote) ExitStandBy() {

	log.Println("Remote exited standby mode")

	r.standby = false
}
