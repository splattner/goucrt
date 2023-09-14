package integration

import (
	"fmt"
	"log"

	"github.com/splattner/goucrt/pkg/entities"
)

// Return the ID of an entity
func GetEntityId(entity interface{}) string {
	log.Println("Get ID of entity with type: " + fmt.Sprintf("%T", entity))

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

// Call the correct HandleCommand function depending on the entity type
func HandleCommand(entity interface{}, req interface{}) {
	log.Println("Handle the entiy_command in the correct entity: " + fmt.Sprintf("%T", entity))

	// Ugly.. I guess but I don't know how better
	switch e := entity.(type) {
	case *entities.ButtonEntity:
		e.HandleCommand(req.(*entities.EntityCommandReq))

	case *entities.LightEntity:
		e.HandleCommand(req.(*entities.EntityCommandReq))

	case *entities.SwitchsEntity:
		e.HandleCommand(req.(*entities.EntityCommandReq))

	case *entities.MediaPlayerEntity:
		e.HandleCommand(req.(*entities.EntityCommandReq))

	case *entities.ClimateEntity:
		e.HandleCommand(req.(*entities.EntityCommandReq))

	case *entities.CoverEntity:
		e.HandleCommand(req.(*entities.EntityCommandReq))
	}
}
