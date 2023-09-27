package deconz

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Stop the listen Loop
func (d *Deconz) Stop() {

	d.controlChannel <- "stop"

}

// Connect to DeCONZ Websocket and start listening for events
func (d *Deconz) StartandListenLoop() {

	log.Info("Deconz, Starting Deconz Websocket Loop")

	ticker := time.NewTicker(pingPeriod)

	socketUrl := fmt.Sprintf("ws://%s:%d", d.host, d.websocketport)
	log.WithField("SocketURL", socketUrl).Debug("Deconz,Trying to connect to Deconz Websocket")
	ws, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Deconz, Error connecting to Websocket Server:", err)
	}
	log.Debugln("Deconz, Connected to Deconz websocket")

	defer func() {
		log.WithField("RemoteAddr", ws.RemoteAddr().String()).Info("Closing Websocket")
		ws.Close()
		ticker.Stop()
	}()

	go d.websocketReceiveHandler(ws)

	// Our main loop for the client
	// We send our relevant packets here
	log.Debugln("Deconz, Starting Deconz Websocket client main loop")
	for {
		select {
		case <-d.controlChannel:
			log.Debug("Closing write loop as read loop closed")
			return

		case <-ticker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			log.WithField("RemoteAddr", ws.RemoteAddr().String()).Debug("Deconz, Send Ping Message")
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.WithField("RemoteAddr", ws.RemoteAddr().String()).Info("Could not send Ping message")
				return
			}
		}
	}
}

// Read from Websocket and process events
func (d *Deconz) websocketReceiveHandler(ws *websocket.Conn) {

	log.Info("Deconz, Starting Deconz Websocket receive handler")

	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	defer func() {
		log.WithField("RemoteAddr", ws.RemoteAddr().String()).Info("Closing Websocket, not able to read message anymore")
		ws.Close()
		// Notify Write looü
		d.controlChannel <- ws.RemoteAddr().String()
	}()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.WithError(err).Debug("Deconz, Error in Deconz Websocket Message receive")
			return
		}

		log.WithField("message", string(msg)).Trace("Deconz, Received Deconz Websocket Message")

		var message DeconzWebSocketMessage
		err = json.Unmarshal(msg, &message)

		if err != nil {
			log.WithError(err).Debug("Unmarshal to DeconzWebSocketMessage failed")
			return
		}

		// Handling light Resources
		if message.Type == "event" && message.Resource == "lights" && message.Event == "changed" {
			if message.State.On != nil ||
				message.State.Hue != nil ||
				message.State.Effect != "" ||
				message.State.Bri != nil ||
				message.State.Sat != nil ||
				message.State.CT != nil ||
				message.State.Reachable != nil ||
				message.State.ColorMode != "" ||
				message.State.ColorLoopSpeed != nil {
				// only if some state acually changed

				for _, l := range d.allDeconzDevices {
					if l.Type == LightDeconzDeviceType {
						if fmt.Sprint(l.Light.ID) == message.ID {
							log.WithFields(log.Fields{
								"ID":   l.Light.ID,
								"Name": l.Light.Name}).Debug("Deconz Websocket changed event for light")
							l.updateState(&message.State)
							l.stateChangeHandler(&message.State)
							break
						}

					}

				}
			}
		}

		// Handling group Resources
		if message.Type == "event" && message.Resource == "groups" && message.Event == "changed" {

			for _, l := range d.allDeconzDevices {
				if l.Type == GroupDeconzDeviceType {
					if fmt.Sprint(l.Group.ID) == message.ID {
						log.WithFields(log.Fields{
							"ID":   l.Group.ID,
							"Name": l.Group.Name}).Debug("Deconz Websocket changed event for group")
						l.updateState(&message.State)
						l.stateChangeHandler(&message.State)
						break
					}

				}

			}

		}

		// Handling sensor Resources
		if message.Type == "event" && message.Resource == "sensors" && message.Event == "changed" {

			for _, l := range d.allDeconzDevices {
				if l.Type == SensorDeconzDeviceType {
					if fmt.Sprint(l.Sensor.ID) == message.ID {
						// Send to all devices which handles this sensor
						log.WithFields(log.Fields{
							"ID":   l.Sensor.ID,
							"Name": l.Sensor.Name}).Debug("Deconz, Websocket changed event for sensor")
						l.updateState(&message.State)
						l.stateChangeHandler(&message.State)
					}
				}
			}
		}
	}
}