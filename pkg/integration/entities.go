package integration

import (
	"github.com/splattner/goucrt/pkg/entities"
)

// Return the ID of an entity
func GetEntityId(entity interface{}) string {
	var id string

	// Ugly.. I guess but I don't know how better
	switch e := entity.(type) {
	case *entities.ButtonEntity:
		id = e.Id

	case *entities.LightEntity:
		id = e.Id

	case *entities.SwitchsEntity:
		id = e.Id

	case *entities.MediaPlayerEntity:
		id = e.Id

	case *entities.SensorEntity:
		id = e.Id

	case *entities.ClimateEntity:
		id = e.Id

	case *entities.CoverEntity:
		id = e.Id
	}

	return id
}

// Return the DeviceId of an entity
func GetDeviceId(entity interface{}) string {
	var device_id string

	// Ugly.. I guess but I don't know how better
	switch e := entity.(type) {
	case *entities.ButtonEntity:
		device_id = e.DeviceId

	case *entities.LightEntity:
		device_id = e.DeviceId

	case *entities.SwitchsEntity:
		device_id = e.DeviceId

	case *entities.MediaPlayerEntity:
		device_id = e.DeviceId

	case *entities.SensorEntity:
		device_id = e.DeviceId

	case *entities.ClimateEntity:
		device_id = e.DeviceId

	case *entities.CoverEntity:
		device_id = e.DeviceId
	}

	return device_id
}

// Return the EntityType of an entity
func GetEntityType(entity interface{}) entities.EntityType {
	var entity_type entities.EntityType

	// Ugly.. I guess but I don't know how better
	switch e := entity.(type) {
	case *entities.ButtonEntity:
		entity_type = e.EntityType

	case *entities.LightEntity:
		entity_type = e.EntityType

	case *entities.SwitchsEntity:
		entity_type = e.EntityType

	case *entities.MediaPlayerEntity:
		entity_type = e.EntityType

	case *entities.SensorEntity:
		entity_type = e.EntityType

	case *entities.ClimateEntity:
		entity_type = e.EntityType

	case *entities.CoverEntity:
		entity_type = e.EntityType
	}

	return entity_type
}

// Return the EntityType of an entity
func GetEntityAttributes(entity interface{}) map[string]interface{} {
	var attributes map[string]interface{}

	// Ugly.. I guess but I don't know how better
	switch e := entity.(type) {
	case *entities.ButtonEntity:
		attributes = e.Attributes

	case *entities.LightEntity:
		attributes = e.Attributes

	case *entities.SwitchsEntity:
		attributes = e.Attributes

	case *entities.MediaPlayerEntity:
		attributes = e.Attributes

	case *entities.SensorEntity:
		attributes = e.Attributes

	case *entities.ClimateEntity:
		attributes = e.Attributes

	case *entities.CoverEntity:
		attributes = e.Attributes
	}

	return attributes
}

// Call the correct HandleCommand function depending on the entity type
func HandleCommand(entity interface{}, req *EntityCommandReq) {
	cmd_id := req.MsgData.CmdId
	params := req.MsgData.Params

	// Ugly.. I guess but I don't know how better
	switch e := entity.(type) {
	case *entities.ButtonEntity:
		e.HandleCommand(cmd_id, params)

	case *entities.LightEntity:
		e.HandleCommand(cmd_id, params)

	case *entities.SwitchsEntity:
		e.HandleCommand(cmd_id, params)

	case *entities.MediaPlayerEntity:
		e.HandleCommand(cmd_id, params)

	case *entities.ClimateEntity:
		e.HandleCommand(cmd_id, params)

	case *entities.CoverEntity:
		e.HandleCommand(cmd_id, params)
	}
}
