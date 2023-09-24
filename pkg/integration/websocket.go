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
		log.WithError(err).Fatal("Cannot upgrade connection")
	}

	log.WithField("RemoteAddr", ws.RemoteAddr().String()).Info("Unfolded Circle Remote connected")

	// Start reading those messages
	go i.wsReader(ws)
	go i.wsWriter(ws)

	i.SendAuthenticationResponse()

}

func (i *Integration) wsReader(ws *websocket.Conn) {
	log.Debug("Start Websocket read loop")
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	defer func() {
		log.WithField("RemoteAddr", ws.RemoteAddr().String()).Info("Closing Websocket, not able to read message anymore")
		ws.Close()

		// Close Write loop also
		i.Remote.controlChannel <- ws.RemoteAddr().String()
	}()

	for {
		_, p, err := ws.ReadMessage()

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
			"RemoteAddr": ws.RemoteAddr().String(),
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

func (i *Integration) wsWriter(ws *websocket.Conn) {
	log.Debug("Start Websocket write loop")
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		log.WithField("RemoteAddr", ws.RemoteAddr().String()).Info("Closing Websocket")
		ws.Close()
		ticker.Stop()
		// Close Read loop
		i.Remote.controlChannel <- ws.RemoteAddr().String()
	}()

	for {
		select {

		case msg := <-i.Remote.controlChannel:
			// Close the writer if message was for this websocket. Closed by reader
			if ws.RemoteAddr().String() == msg {
				log.Debug("Closing write loop as read loop closed")
				return
			}

		case msg := <-i.Remote.messageChannel:

			// Remote should not be in standby as this is a response to a request
			// or if sent from sendEventMessage the sendEventMessage function makes sure the remote is not in standby
			log.WithFields(log.Fields{
				"RawMessage": string(msg),
				"RemoteAddr": ws.RemoteAddr().String()}).Debug("Send message to websocket")

			if err := ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.WithError(err).Error("Faled to set WriteDeatLine")
			}

			if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.WithError(err).Error("Failed to send message")
			}

		case <-ticker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			log.WithField("RemoteAddr", ws.RemoteAddr().String()).Debug("Send Ping Message")
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.WithField("RemoteAddr", ws.RemoteAddr().String()).Info("Could not send Ping message")
				return
			}
		}
	}
}
