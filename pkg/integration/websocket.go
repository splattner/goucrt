package integration

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
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

	log.Println("Unfolded Circle Remote two connected")

	i.websocket = ws

	// Start reading those messages
	i.reader()
}

func (i *Integration) reader() {
	for {
		_, p, err := i.websocket.ReadMessage()
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

func (i *Integration) sendEventMessage(res *interface{}, messageType int) error {
	log.Println("Send Event Message")

	msg, _ := json.Marshal(res)
	log.Println(string(msg))

	if i.Remote.standby {
		log.Println("Remote is in standby mode, not sending event")
		return nil
	}

	return i.websocket.WriteMessage(messageType, msg)

}

func (i *Integration) sendResponseMessage(res *interface{}, messageType int) error {
	log.Println("Send Response Message")

	msg, _ := json.Marshal(res)
	log.Println(string(msg))

	return i.websocket.WriteMessage(messageType, msg)

}
