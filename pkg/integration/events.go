package integration

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func (i *Integration) handleEvent(req *RequestMessage, p []byte) interface{} {
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

		i.handleConnectEvent(&connectEvent)

	case "disconnect":
		log.Println("connect event")

		connectEvent := ConnectEvent{}
		json.Unmarshal(p, &connectEvent)

		i.handleConnectEvent(&connectEvent)

	case "abort_driver_setup":
		log.Println("abort_driver_setup event")

	default:
		log.Println("mesage not know: " + req.Msg)
	}

	return res

}

func (i *Integration) sendEntityRemoved(e interface{}) {

	var res interface{}

	msg_data := EntityRemovedEventData{
		DeviceId:   GetDeviceId(e),
		EntityType: GetEntityType(e).Type,
		EntityId:   GetEntityId(e),
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

func (i *Integration) sendEntityAvailable(e interface{}) {

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

func (i *Integration) sendDeviceStateEvent() {

	var res interface{}

	now := time.Now()
	res = DeviceStateEventMessage{
		CommonEvent{Kind: "event", Msg: "device_state", Cat: "DEVICE", Ts: now.Format(time.UnixDate)},
		DeviceState{DeviceId: DeviceId{DeviceId: i.DeviceId}, State: string(i.deviceState)},
	}

	i.sendEventMessage(&res, websocket.TextMessage)
}

// Emitted for all driver setup flow state changes.
func (i *Integration) sendDriverSetupChangeEvent(eventType DriverSetupEventType, state DriverSetupState, err DriverSetupError, required_user_action *RequiredUserAction) {
	var res interface{}

	now := time.Now()

	if required_user_action == nil {
		res = DriverSetupChangeEvent{
			CommonEvent{Kind: "event", Msg: "driver_setup_change", Cat: "DEVICE", Ts: now.Format(time.UnixDate)},
			DriverSetupChangeData{EventType: eventType, State: state, Error: err},
		}
	} else {
		res = DriverSetupChangeEvent{
			CommonEvent{Kind: "event", Msg: "driver_setup_change", Cat: "DEVICE", Ts: now.Format(time.UnixDate)},
			DriverSetupChangeData{EventType: eventType, State: state, Error: err, RequiredUserAction: *required_user_action},
		}
	}

	i.sendEventMessage(&res, websocket.TextMessage)

}

func (i *Integration) handleConnectEvent(e *ConnectEvent) {

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
