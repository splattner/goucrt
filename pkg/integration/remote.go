package integration

import (
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

type remote struct {
	standby   bool
	connected bool
	websocket *websocket.Conn
}

func (r *remote) EnterStandBy() {
	log.Info("Remote entered standby mode")

	r.standby = true

}

func (r *remote) ExitStandBy() {
	log.Info("Remote exited standby mode")

	r.standby = false
}
