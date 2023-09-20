package entities

import (
	log "github.com/sirupsen/logrus"
)

type EntityState string

const (
	UnavailableEntityState EntityState = "UNAVAILABLE"
	UnkownEntityState                  = "UNKNOWN"
)

type Entity struct {
	Id string `json:"entity_id"`
	EntityType
	DeviceId               string                 `json:"device_id,omitempty"`
	Features               []interface{}          `json:"features"`
	Name                   LanguageText           `json:"name"`
	Area                   string                 `json:"area,omitempty"`
	DeviceClass            string                 `json:"-"`
	Attributes             map[string]interface{} `json:"-"`
	handleEntityChangeFunc func(interface{})      `json:"-"`
}

type EntityType struct {
	Type string `json:"entity_type,omitempty"`
}

type EntityStateData struct {
	DeviceId string `json:"device_id,omitempty"`
	EntityType
	EntityId   string                 `json:"entity_id"`
	Attributes map[string]interface{} `json:"attributes"`
}

// Add an attribute if not already available
func (e *Entity) AddAttribute(name string, value interface{}) {
	if e.Attributes[name] != nil {
		e.Attributes[name] = value
	}
}

func (e *Entity) GetEntityState() *EntityStateData {

	entityState := EntityStateData{
		DeviceId:   e.DeviceId,
		EntityType: e.EntityType,
		EntityId:   e.Id,
		Attributes: e.Attributes,
	}

	return &entityState
}

func (e *Entity) SetHandleEntityChangeFunc(f func(interface{})) {
	log.WithField("entity_id", e.Id).Debug("Entity Change Handler Set")
	e.handleEntityChangeFunc = f
}

func (e *Entity) SetAttributes(attributes map[string]interface{}) {

	log.WithField("Attributes", attributes).Info("Handle attribute change")

	for k, v := range attributes {
		e.Attributes[k] = v
	}

	// Handle the entity Change
	log.WithField("entity_id", e.Id).Debug("Calling the Entity Change Handler if set")
	if e.handleEntityChangeFunc != nil {
		e.handleEntityChangeFunc(e)
	}

}
