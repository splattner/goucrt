package integration

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

func (i *Integration) sendResponseMessage(res interface{}, messageType int) error {

	msg, err := json.Marshal(res)
	if err != nil {
		return err
	}

	// Unmarshal againinto Event Message for some fields
	response := ResponseMessage{}
	json.Unmarshal(msg, &response)

	log.WithFields(log.Fields{
		"Message":    response.Msg,
		"Id":         response.Id,
		"Kind":       response.Kind,
		"Data":       response.MsgData,
		"RawMessage": string(msg)}).Info("Send Response Message")

	// Remote should not be in standby as this is a response to a request
	if !i.Remote.connected || i.Remote.websocket == nil {
		return fmt.Errorf("No Open Websocket connection, cannot send a response")
	}

	if err := i.Remote.websocket.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return err
	}

	log.WithField("RawMessage", string(msg)).Info("Send Response Message")

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
