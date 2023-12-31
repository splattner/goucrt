package integration

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
	"k8s.io/utils/strings/slices"
)

func (i *Integration) sendEventMessage(res *interface{}, messageType int) error {

	msg, _ := json.Marshal(res)

	// Unmarshal againinto Event Message for some fields
	event := EventMessage{}
	if err := json.Unmarshal(msg, &event); err != nil {
		log.WithError(err).Error("Cannot unmarshal Event Message")
		return err
	}

	if i.Remote.standby {
		log.WithFields(log.Fields{
			"Message": event.Msg,
			"Kind":    event.Kind,
			"standby": i.Remote.standby,
		}).Info("Remote is in standby mode or not (yet) connected, not sending event / no websocket")
		return nil
	}

	log.WithFields(log.Fields{
		"Message": event.Msg,
		"Kind:":   event.Kind,
		"Data":    event.MsgData,
	}).Info("Send Event Message")

	select {
	case i.Remote.messageChannel <- msg:
		log.Debug("Message sent for processing")
	default:
		log.Debug("Message Channel not ready")

	}

	return nil

}

// Handle events received from the Remote
func (i *Integration) handleEvent(req *RequestMessage, p []byte) interface{} {

	var res interface{}

	switch req.Msg {
	case "enter_standby":
		i.Remote.EnterStandBy()

	case "exit_standby":
		i.Remote.ExitStandBy()

	case "connect":
		connectEvent := ConnectEvent{}
		if err := json.Unmarshal(p, &connectEvent); err != nil {
			log.WithError(err).Error("Cannot unmarshal ConnectEvent")
			return nil
		}

		i.handleConnectEvent(&connectEvent)

	case "disconnect":
		connectEvent := ConnectEvent{}
		if err := json.Unmarshal(p, &connectEvent); err != nil {
			log.WithError(err).Error("Cannot unmarshal ConnectEvent")
			return nil
		}
		i.handleConnectEvent(&connectEvent)

	case "abort_driver_setup":

		abortDriverSetupEvent := AbortDriverSetupEvent{}

		if err := json.Unmarshal(p, &abortDriverSetupEvent); err != nil {
			log.WithError(err).Error("Cannot unmarshal AbortDriverSetupEvent")
			return nil
		}

		i.handleAbortDriverSetupEvent(&abortDriverSetupEvent)

	default:
		log.WithField("Message", req.Msg).Debug("Mesage not know")
	}

	return res

}

func (i *Integration) sendEntityRemoved(e interface{}) {

	var res interface{}
	now := time.Now()

	msg_data := EntityRemovedEventData{
		DeviceId:   i.getDeviceId(e),
		EntityType: i.getEntityType(e).Type,
		EntityId:   i.getEntityId(e),
	}

	res = EntityRemovedEvent{
		CommonEvent{
			Kind: "event",
			Msg:  "entity_removed",
			Cat:  "ENTITY",
			Ts:   now.Format(time.RFC3339),
		},
		msg_data,
	}

	if err := i.sendEventMessage(&res, websocket.TextMessage); err != nil {
		log.WithError(err).Error("Cannot send Event Message")
	}
}

func (i *Integration) sendEntityAvailable(e interface{}) {

	var res interface{}
	now := time.Now()

	res = EntityAvailableEvent{
		CommonEvent{
			Kind: "event",
			Msg:  "entity_available",
			Cat:  "ENTITY",
			Ts:   now.Format(time.RFC3339),
		},
		e,
	}

	// Only send event when connected, otherwise we assume this is still during setup e.g. discovering of entities
	if i.deviceState == ConnectedDeviceState {
		if err := i.sendEventMessage(&res, websocket.TextMessage); err != nil {
			log.WithError(err).Error("Cannot send Event Message")
		}
	}
}

func (i *Integration) sendDeviceStateEvent() {

	var res interface{}

	now := time.Now()
	res = DeviceStateEventMessage{
		CommonEvent{Kind: "event", Msg: "device_state", Cat: "DEVICE", Ts: now.Format(time.RFC3339)},
		DeviceState{DeviceId: DeviceId{DeviceId: i.DeviceId}, State: string(i.deviceState)},
	}

	if err := i.sendEventMessage(&res, websocket.TextMessage); err != nil {
		log.WithError(err).Error("Cannot send Event Message")
	}
}

// Emitted for all driver setup flow state changes.
func (i *Integration) sendDriverSetupChangeEvent(eventType DriverSetupEventType, state DriverSetupState, err DriverSetupError, require_user_action *RequireUserAction) {
	var res interface{}

	now := time.Now()

	if require_user_action == nil {
		res = DriverSetupChangeEvent{
			CommonEvent{Kind: "event", Msg: "driver_setup_change", Cat: "DEVICE", Ts: now.Format(time.RFC3339)},
			DriverSetupChangeData{EventType: eventType, State: state, Error: err},
		}
	} else {
		res = DriverSetupChangeEvent{
			CommonEvent{Kind: "event", Msg: "driver_setup_change", Cat: "DEVICE", Ts: now.Format(time.RFC3339)},
			DriverSetupChangeData{EventType: eventType, State: state, Error: err, RequireUserAction: *require_user_action},
		}
	}

	if err := i.sendEventMessage(&res, websocket.TextMessage); err != nil {
		log.WithError(err).Error("Cannot send Event Message")
	}

}

func (i *Integration) handleConnectEvent(e *ConnectEvent) {

	// Cat should be "DEVICE"

	switch e.Msg {
	case "connect":
		// Call the handler of the client
		if i.handleConnectionFunction != nil {
			i.handleConnectionFunction(e)
		}

	case "disconnect":
		// Call the handler of the client
		if i.handleConnectionFunction != nil {
			i.handleConnectionFunction(e)
		}
	}

}

// If the user aborts the setup process, the Remote Two sends this event.
// Further messages from the integration from the setup process will be ignored afterwards.
func (i *Integration) handleAbortDriverSetupEvent(e *AbortDriverSetupEvent) {
	log.Info("Abort Driver Setup")
	// TODO: implement something?
}

// Emitted when an attribute of an entity changes, e.g. is switched off.
// Either after an entity_command or if the entity is updated manually through a user or an external system.
// This keeps the Remote Two in sync with the real state of the entity without the need of constant polling.
func (i *Integration) SendEntityChangeEvent(e interface{}, a *map[string]interface{}) {

	entity_id := i.getEntityId(e)

	log.WithField("entity_id", entity_id).Debug("Send Entity Change Event if subscribed")
	log.WithField("subscribedEtities", i.SubscribedEntities).Debug("Currently subscribed entities")

	// Only send the event when remote is subscribed to
	if i.Config.IgnoreEntitySubscription || slices.Contains(i.SubscribedEntities, entity_id) {

		var res interface{}

		// UTC time for event timestamp
		loc, _ := time.LoadLocation("UTC")
		now := time.Now().In(loc)
		timeformat := "2006-01-02T15:04:05.999999999Z"

		device_id := i.getDeviceId(e)

		entity_type := i.getEntityType(e)

		var attributes map[string]interface{}
		//if attributes is set, only send thos
		if a == nil {
			attributes = i.getEntityAttributes(e)
		} else {
			attributes = *a
		}

		res = EntityChangeEvent{
			CommonEvent{
				Kind: "event",
				Msg:  "entity_change",
				Cat:  "ENTITY",
				Ts:   now.Format(timeformat),
			},
			EntityChangeData{
				DeviceId:   device_id,
				EntityId:   entity_id,
				EntityType: entity_type.Type,
				Attributes: attributes,
			},
		}

		if err := i.sendEventMessage(&res, websocket.TextMessage); err != nil {
			log.WithError(err).Error("Cannot send Event Message")
		}
	}

}
