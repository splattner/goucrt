package integration

import (
	log "github.com/sirupsen/logrus"

	"github.com/grandcat/zeroconf"
)

// Start Advertising the integration with mDNS
func (i *Integration) startAdvertising() {
	log.Info("Start advertising UC Integration with mDNS")

	txt := []string{
		"name=" + i.Metadata.Name.En,
		"developer=" + i.Metadata.Developer.Name,
		"ver=" + i.Metadata.Version,
		"ws_path=" + i.Config.WebsocketPath,
	}

	server, err := zeroconf.Register(i.Metadata.DriverId, "_uc-integration._tcp", "local.", i.Config.ListenPort, txt, nil)
	if err != nil {
		panic(err)
	}

	i.mdns = server
}

// Stop mDNS advertisement
func (i *Integration) stopAdvertising() {
	if i.mdns != nil {
		log.Info("Stop advertising UC Integration")
		i.mdns.Shutdown()
	}
}
