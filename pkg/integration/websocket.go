package integration

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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
		log.Println(err)
	}

	log.Println("Unfolded Circle Remote with addr " + r.RemoteAddr + " connected")

	if i.Remote.websocket != nil {
		// TODO: do we need to support more?
		log.Println("There is already a websocket connection open, cannot open an othner one")
	}

	i.Remote.websocket = ws
	i.Remote.connected = true

	i.SendAuthenticationResponse()

	// Start reading those messages
	go i.wsReader()
	go i.wsWriter()
}

func (i *Integration) wsReader() {
	i.Remote.websocket.SetReadLimit(maxMessageSize)
	i.Remote.websocket.SetReadDeadline(time.Now().Add(pongWait))
	i.Remote.websocket.SetPongHandler(func(string) error {
		i.Remote.websocket.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	defer func() {
		log.Println("Closing Websocket, not able to read message anymore")
		i.Remote.websocket.Close()

		i.Remote.websocket = nil
		i.Remote.connected = false
	}()

	for {
		_, p, err := i.Remote.websocket.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		req := RequestMessage{}

		if json.Unmarshal(p, &req) != nil {
			log.Println("Cannot unmarshall " + string(p))
		}

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
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		log.Println("Closing Websocket, no response in time to Ping message")
		ticker.Stop()
		i.Remote.websocket.Close()

		i.Remote.websocket = nil
		i.Remote.connected = false
	}()

	for {
		select {
		case <-ticker.C:
			if i.Remote.websocket != nil && i.Remote.connected {
				i.Remote.websocket.SetWriteDeadline(time.Now().Add(writeWait))
				if err := i.Remote.websocket.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Println("Could not send Ping message")
					return
				}
			}
		}
	}
}

func (i *Integration) sendEventMessage(res *interface{}, messageType int) error {
	log.Println("Send Event Message")

	msg, _ := json.Marshal(res)
	log.Println(string(msg))

	if i.Remote.standby || !i.Remote.connected {
		log.Println("Remote is in standby mode or not (yet) connected, not sending event")
		return nil
	}

	return i.Remote.websocket.WriteMessage(messageType, msg)

}
