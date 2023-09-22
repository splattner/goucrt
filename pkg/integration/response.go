package integration

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

// Send a generic Response Message to Remote Two
func (i *Integration) sendResponseMessage(res interface{}, messageType int) error {

	msg, err := json.Marshal(res)
	if err != nil {
		return err
	}

	// Unmarshal againinto Event Message for some fields
	response := ResponseMessage{}
	json.Unmarshal(msg, &response)

	log.WithFields(log.Fields{
		"Message": response.Msg,
		"Id":      response.Id,
		"Kind":    response.Kind,
		"Data":    response.MsgData}).Info("Send Response Message")

	i.Remote.messageChannel <- msg

	return nil

}

func (i *Integration) authenticationResponseMessage(code int) *AuthenticationResponse {

	res := AuthenticationResponse{
		CommonResp{
			Kind: "resp",
			Id:   0,
			Msg:  "authentication",
			Code: code,
		},
		DriverVersionData{
			Name: i.Metadata.Name.En,
			Version: Version{
				Api:    API_VERSION,
				Driver: API_VERSION,
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
