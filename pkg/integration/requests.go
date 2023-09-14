package integration

import (
	"encoding/json"
	"log"

	"github.com/splattner/goucrt/pkg/integration/entities"
)

func (i *integration) handleRequest(req *RequestMessage, p []byte) interface{} {
	log.Println("Handle Request Message")

	var res interface{}

	switch req.Msg {
	case "auth_required":
		log.Println("auth required")

		authRequiredReq := AuthRequestMessage{}
		json.Unmarshal(p, &authRequiredReq)

		// TODO
		//res = i.handleAuthRequired(&authRequiredReq)

	case "get_driver_version":
		log.Println("get driver version")

		driverVersionReq := DriverVersionReq{}
		json.Unmarshal(p, &driverVersionReq)

		res = i.handleGetDriverVersionRequest(&driverVersionReq)

	case "get_driver_metadata":
		log.Println("get driver metadata")

		driverMetadataReq := DriverMetadataReq{}
		json.Unmarshal(p, &driverMetadataReq)

		res = i.getDriverMetadata(&driverMetadataReq)

	case "get_device_state":
		log.Println("get device state")

		deviceStateMessageReq := DeviceStateMessageReq{}
		json.Unmarshal(p, &deviceStateMessageReq)

		i.handleGetDeviceStateRequest(&deviceStateMessageReq)

	case "get_available_entities":
		log.Println("get_available_entities")

		availableEntityMessageReq := AvailableEntityMessageReq{}
		json.Unmarshal(p, &availableEntityMessageReq)

		res = i.handleGetAvailableEntitiesRequest(&availableEntityMessageReq)
	case "subscribe_events":
		log.Println("subscribe_events")

		subscribeEventMessageReq := SubscribeEventMessageReq{}
		json.Unmarshal(p, &subscribeEventMessageReq)

		res = i.handleSubscribeEventRequest(&subscribeEventMessageReq)
	case "unsubscribe_events":
		log.Println("unsubscribe_events")

		unsubscribeEventMessageReq := UnubscribeEventMessageReq{}
		json.Unmarshal(p, &unsubscribeEventMessageReq)

		res = i.handleUnsubscribeEventsRequest(&unsubscribeEventMessageReq)

	case "get_entity_state":
		log.Println("get_entity_state")

		entityStateReq := GetEntityStateMessageReq{}
		json.Unmarshal(p, &entityStateReq)

		res = i.getEntityState(&entityStateReq)

	case "entity_command":
		log.Println("entity_command")

		entityCommandReq := EntityCommandReq{}
		json.Unmarshal(p, &entityCommandReq)

		res = i.handleEntityCommandRequest(&entityCommandReq)

	case "setup_driver":
		log.Println("setup_driver")

		setupDriverReq := SetupDriverMessageReq{}
		json.Unmarshal(p, &setupDriverReq)

		res = i.handleSetupDriverRequest(&setupDriverReq)

	case "set_driver_user_data":
		log.Println("set_driver_user_data")

		setUserData := SetDriverUserDataReq{}
		json.Unmarshal(p, &setUserData)

		res = i.handleSetUserDataRequest(&setUserData)

	default:
		log.Println("mesage not know")
	}

	return res
}

func (i *integration) handleGetDeviceStateRequest(req *DeviceStateMessageReq) {

	i.sendDeviceStateEvent()

}

func (i *integration) handleGetDriverVersionRequest(req *DriverVersionReq) *ResponseMessage {

	msg_data := DriverVersionData{
		Name: i.Metadata.Name.En,
		Version: Version{
			Api:    "test",
			Driver: i.Metadata.Version,
		},
	}

	res := ResponseMessage{
		CommonResp{
			Kind: "resp",
			Id:   req.Id,
			Msg:  "driver_version",
			Code: 200,
		},
		msg_data,
	}

	return &res

}

func (i *integration) getDriverMetadata(req *DriverMetadataReq) *DriverMetadataReponse {

	res := DriverMetadataReponse{
		CommonResp{
			Kind: "resp",
			Id:   req.Id,
			Msg:  "driver_metadata",
			Code: 200,
		},
		i.Metadata,
	}

	return &res

}

func (i *integration) handleGetAvailableEntitiesRequest(req *AvailableEntityMessageReq) *AvailableEntityMessage {

	res := AvailableEntityMessage{
		CommonResp{Kind: "resp", Id: req.Id, Msg: "available_entities", Code: 200},
		AvailableEntityData{
			Filter:            req.MsgData.Filter,
			AvailableEntities: i.Entities,
		},
	}

	return &res

}

func (i *integration) handleSetupDriverRequest(req *SetupDriverMessageReq) *ResponseMessage {

	// Todo: Call Setup driver
	// Send Event?

	res := ResponseMessage{
		CommonResp{
			Kind: "resp",
			Id:   req.Id,
			Msg:  req.Msg,
			Code: 200,
		},
		nil,
	}

	return &res

}

// Subscribe to entity state change events to receive entity_change events from the integration driver.
// If no entity IDs are specified then events for all available entities are sent to the Remote Two.
func (i *integration) handleSubscribeEventRequest(req *SubscribeEventMessageReq) *SubscribeEventMessage {

	res := SubscribeEventMessage{
		CommonResp{Kind: "resp", Id: req.Id, Msg: req.Msg, Code: 200},
	}

	// TODO: Implement

	return &res

}

// If no entity IDs are specified then all events for all available entities are stopped.
// This message is sent by the Remote Two if a previously configured entity is no longer used and therefore no longer interested in entity updates. If the integration driver keeps sending events for the unsubscribed entities then they are simply discarded.
func (i *integration) handleUnsubscribeEventsRequest(req *UnubscribeEventMessageReq) *UnubscribeEventMessage {

	res := UnubscribeEventMessage{
		CommonResp{Kind: "resp", Id: req.Id, Msg: req.Msg, Code: 200},
	}

	// TODO: implement

	return &res

}

func (i *integration) getEntityState(req *GetEntityStateMessageReq) *GetEntityStateMessage {

	// TODO, get the entities
	var entityState = []entities.EntityStateData{}

	entityA := entities.EntityStateData{
		DeviceId:   i.DeviceId,
		EntityType: entities.EntityType{"Dummy"},
		EntityId:   "newid",
		Attributes: nil,
	}

	entityState = append(entityState, entityA)

	res := GetEntityStateMessage{
		CommonResp{Kind: "resp", Id: req.Id, Msg: "entity_states", Code: 200},
		entityState,
	}

	return &res

}

func (i *integration) handleEntityCommandRequest(req *EntityCommandReq) *EntityCommandResponse {

	res := EntityCommandResponse{
		CommonResp{Kind: "resp", Id: req.Id, Msg: req.Msg, Code: 200},
	}

	return &res

}

func (i *integration) handleSetUserDataRequest(req *SetDriverUserDataReq) *EntityCommandResponse {

	res := EntityCommandResponse{
		CommonResp{Kind: "resp", Id: req.Id, Msg: req.Msg, Code: 200},
	}

	return &res
}
