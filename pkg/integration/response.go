package integration

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func (i *Integration) sendResponseMessage(res interface{}, messageType int) error {
	log.Println("Send Response Message")

	msg, _ := json.Marshal(res)
	log.Println(string(msg))

	if !i.Remote.connected || i.Remote.websocket == nil {
		return fmt.Errorf("No Open Websocket connection, cannot send a response")
	}

	return i.Remote.websocket.WriteMessage(messageType, msg)

}

func (i *Integration) authenticationResponseMessage(code int) *AuthenticationResponse {

	res := AuthenticationResponse{
		CommonResp{
			Kind: "resp",
			Id:   0,
			Msg:  "authentication",
			Code: code,
		},
		AuthenticationResponseData{
			Name: i.Metadata.Name.En,
			Version: Version{
				Api:    i.Metadata.Version,
				Driver: i.Metadata.Version,
			},
		},
	}

	return &res
}

// Send a AuthenticationResponse Message
func (i *Integration) SendAuthenticationResponse() {
	msg := i.authenticationResponseMessage(200)

	i.sendResponseMessage(msg, websocket.TextMessage)
}
