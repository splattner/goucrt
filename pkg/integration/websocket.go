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

func (i *integration) wsEndpoint(w http.ResponseWriter, r *http.Request) {

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

func (i *integration) reader() {
	for {
		messageType, p, err := i.websocket.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		req := RequestMessage{}
		var res interface{}

		if json.Unmarshal(p, &req) != nil {
			log.Println("Cannot unmarshall")
		}

		// Event Message
		if req.Kind == "event" {
			res = i.handleEvent(&req, p)
		}

		// Request Message
		if req.Kind == "req" {
			res = i.handleRequest(&req, p)

			if err := i.sendResponseMessage(&res, messageType); err != nil {
				log.Println(err)
			}
		}

	}
}

func (i *integration) sendEventMessage(res *interface{}, messageType int) error {
	log.Println("Send Event Message")

	msg, _ := json.Marshal(res)
	log.Println(string(msg))

	if i.Remote.standby {
		log.Println("Remote is in standby mode, not sending event")
		return nil
	}

	return i.websocket.WriteMessage(messageType, msg)

}

func (i *integration) sendResponseMessage(res *interface{}, messageType int) error {
	log.Println("Send Response Message")

	msg, _ := json.Marshal(res)
	log.Println(string(msg))

	return i.websocket.WriteMessage(messageType, msg)

}
