package deconz

import (
	"fmt"
	"strings"

	deconzlight "github.com/jurgen-kluft/go-conbee/lights"
	log "github.com/sirupsen/logrus"
)

func (d *DeconzDevice) newDeconzLightDevice() {

}

func (d *DeconzDevice) setLightState() error {

	state := strings.Replace(d.Light.State.String(), "\n", ",", -1)
	state = strings.Replace(state, " ", "", -1)

	log.Infof("Deconz, call SetLightState with state (%s) for Light with id %d\n", state, d.Light.ID)

	conbeehost := fmt.Sprintf("%s:%d", d.deconz.host, d.deconz.port)
	ll := deconzlight.New(conbeehost, d.deconz.apikey)
	_, err := ll.SetLightState(d.Light.ID, &d.Light.State)
	if err != nil {
		log.Debugln("Deconz, SetLightState Error", err)
		return err
	}

	return nil
}
