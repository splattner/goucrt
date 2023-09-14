package integration

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/splattner/goucrt/pkg/integration/entities"
)

func (i *integration) handleEvent(req *RequestMessage, p []byte) interface{} {
	log.Println("Handle Event Message")

	var res interface{}

	switch req.Msg {
	case "enter_standby":
		i.Remote.EnterStandBy()

	case "exit_standby":
		i.Remote.ExitStandBy()

	case "connect":
		log.Println("connect event")

		connectEvent := ConnectEvent{}
		json.Unmarshal(p, &connectEvent)

		i.handleCconnectEvent(&connectEvent)

	case "disconnect":
		log.Println("connect event")

		connectEvent := ConnectEvent{}
		json.Unmarshal(p, &connectEvent)

		i.handleCconnectEvent(&connectEvent)

	case "abort_driver_setup":
		log.Println("abort_driver_setup event")

	default:
		log.Println("mesage not know: " + req.Msg)
	}

	return res

}

func (i *integration) sendEntityRemoved(e entities.Entity) {

	var res interface{}

	msg_data := EntityRemovedEventData{
		DeviceId:   e.DeviceId,
		EntityType: e.EntityType.Type,
		EntityId:   e.Id,
	}

	res = EntityRemovedEvent{
		CommonEvent{
			Kind: "event",
			Msg:  "entity_removed",
			Cat:  "ENTITY",
		},
		msg_data,
	}

	i.sendEventMessage(&res, websocket.TextMessage)
}

func (i *integration) sendEntityAvailable(e entities.Entity) {

	var res interface{}

	res = EntityAvailableEvent{
		CommonEvent{
			Kind: "event",
			Msg:  "entity_available",
			Cat:  "ENTITY",
		},
		e,
	}

	i.sendEventMessage(&res, websocket.TextMessage)
}

func (i *integration) sendDeviceStateEvent() {

	var res interface{}

	res = DeviceStateEventMessage{
		CommonEvent{Kind: "event", Msg: "device_state", Cat: "DEVICE", Ts: "//Todo"},
		DeviceState{DeviceId: DeviceId{DeviceId: i.DeviceId}, State: i.DeviceState},
	}

	i.sendEventMessage(&res, websocket.TextMessage)
}

// Emitted for all driver setup flow state changes.
func (i *integration) sendDriverSetupChangeEvent(eventType DriverSetupEventType, state DriverSetupState, err DriverSetupError, required_user_action interface{}) {
	var res interface{}

	res = DriverSetupChangeEvent{
		CommonEvent{Kind: "event", Msg: "driver_setup_change", Cat: "DEVICE", Ts: "//Todo"},
		DriverSetupChangeData{EventType: eventType, State: state, Error: err, RequiredUserAction: required_user_action},
	}

	i.sendEventMessage(&res, websocket.TextMessage)

}

func (i *integration) handleCconnectEvent(e *ConnectEvent) {

	// Cat should be "DEVICE"

	switch e.Msg {
	case "connect":
		log.Println("connect")

	case "disconnect":
		log.Println("disconnect")

	default:
		log.Println("Unknown connect message")

	}

}
