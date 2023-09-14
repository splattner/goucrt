package integration

import "log"

type remote struct {
	standby bool
}

func (r *remote) EnterStandBy() {
	log.Println("Remote entered standby mode")

	r.standby = true

}

func (r *remote) ExitStandBy() {

	log.Println("Remote exited standby mode")

	r.standby = false

}
