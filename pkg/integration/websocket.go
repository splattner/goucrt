package integration

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (i *Integration) wsEndpoint(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
	}

	log.WithField("RemoteAddr", ws.RemoteAddr().String()).Info("Unfolded Circle Remote connected")

	if i.Remote.websocket != nil {
		// TODO: do we need to support this? Just overwrite the old websocket?
		log.WithField("RemoteAddr", r.RemoteAddr).Info("There is already a websocket connection open")
	}

	i.Remote.websocket = ws
	i.Remote.connected = true

	// Start reading those messages
	go i.wsReader()
	go i.wsWriter()

	i.SendAuthenticationResponse()

}

func (i *Integration) wsReader() {
	log.Debug("Start Websocket read loop")
	i.Remote.websocket.SetReadLimit(maxMessageSize)
	i.Remote.websocket.SetReadDeadline(time.Now().Add(pongWait))
	i.Remote.websocket.SetPongHandler(func(string) error {
		i.Remote.websocket.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	defer func() {
		log.WithField("RemoteAddr", i.Remote.websocket.RemoteAddr().String()).Info("Closing Websocket, not able to read message anymore")
		i.Remote.websocket.Close()

		i.Remote.websocket = nil
		i.Remote.connected = false
	}()

	for {
		_, p, err := i.Remote.websocket.ReadMessage()

		if err != nil {
			log.Error(err)
			return
		}

		req := RequestMessage{}

		if json.Unmarshal(p, &req) != nil {
			log.Error("Cannot unmarshall " + string(p))
			continue
		}

		log.WithFields(log.Fields{
			"RemoteAddr": i.Remote.websocket.RemoteAddr().String(),
			"Message":    req.Msg,
			"Kind":       req.Kind,
			"Id":         req.Id,
		}).Info("Message received")

		// Event Message
		if req.Kind == "event" {
			i.handleEvent(&req, p)
		}

		// Request Message
		if req.Kind == "req" {
			i.handleRequest(&req, p)
		}

	}
}

func (i *Integration) wsWriter() {
	log.Debug("Start Websocket write loop")
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		log.WithField("RemoteAddr", i.Remote.websocket.RemoteAddr().String()).Info("Closing Websocket, no response in time to Ping message")
		ticker.Stop()
		i.Remote.websocket.Close()

		i.Remote.websocket = nil
		i.Remote.connected = false
	}()

	for {
		select {

		case msg := <-i.Remote.messageChannel:

			// Remote should not be in standby as this is a response to a request
			// or if sent from sendEventMessage the sendEventMessage function makes sure the remote is not in standby
			if i.Remote.connected && i.Remote.websocket != nil {
				log.WithFields(log.Fields{
					"RawMessage": string(msg),
					"RemoteAddr": i.Remote.websocket.RemoteAddr().String()}).Debug("Send message to websocket")

				if err := i.Remote.websocket.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
					log.WithError(err).Error("Faled to set WriteDeatLine")
				}

				if err := i.Remote.websocket.WriteMessage(websocket.TextMessage, msg); err != nil {
					log.WithError(err).Error("Failed to send message")
				}

			} else {
				log.Info("Remote not connected")
			}

		case <-ticker.C:
			if i.Remote.websocket != nil && i.Remote.connected {
				i.Remote.websocket.SetWriteDeadline(time.Now().Add(writeWait))
				log.WithField("RemoteAddr", i.Remote.websocket.RemoteAddr().String()).Debug("Send Ping Message")
				if err := i.Remote.websocket.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.WithField("RemoteAddr", i.Remote.websocket.RemoteAddr().String()).Info("Could not send Ping message to")
					return
				}
			}
		}
	}
}
