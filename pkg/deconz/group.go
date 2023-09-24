package deconz

import (
	"fmt"
	"strings"

	deconzgroup "github.com/jurgen-kluft/go-conbee/groups"
	log "github.com/sirupsen/logrus"
)

func (d *DeconzDevice) newDeconzGroupDevice() {

}

func (d *DeconzDevice) setGroupState() error {

	state := strings.Replace(d.Group.Action.String(), "\n", ",", -1)
	state = strings.Replace(state, " ", "", -1)

	log.Infof("Deconz, call SetGroupState with state (%s) for Light with id %d\n", state, d.Group.ID)

	conbeehost := fmt.Sprintf("%s:%d", d.deconz.host, d.deconz.port)
	ll := deconzgroup.New(conbeehost, d.deconz.apikey)
	_, err := ll.SetGroupState(d.Light.ID, d.Group.Action)
	if err != nil {
		log.Debugln("Deconz, SetGroupState Error", err)
		return err
	}
	return nil
}
