package integration

import (
	log "github.com/sirupsen/logrus"
)

type remote struct {
	standby bool
	// Channel to send new messages over websocket.
	messageChannel chan []byte
}

func (r *remote) EnterStandBy() {
	log.Info("Remote entered standby mode")

	r.standby = true

}

func (r *remote) ExitStandBy() {
	log.Info("Remote exited standby mode")

	r.standby = false
}
