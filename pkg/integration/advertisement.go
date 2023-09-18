package integration

import (
	log "github.com/sirupsen/logrus"

	"github.com/grandcat/zeroconf"
)

// TODO: not working?
func (i *Integration) startAdvertising() {
	log.Info("Start advertising UC Integration with mdns")

	txt := []string{
		"name=" + i.Metadata.Name.En,
		"developer=" + i.Metadata.Developer.Name,
		"ver=" + i.Metadata.Version,
		"ws_path=/ws",
	}

	server, err := zeroconf.Register(i.Metadata.DriverId, "_uc-integration._tcp", "local.", i.config["listenport"].(int), txt, nil)
	if err != nil {
		panic(err)
	}

	i.mdns = server
}

func (i *Integration) stopAdvertising() {
	log.Info("Stop advertising UC Integration")

	i.mdns.Shutdown()
}
