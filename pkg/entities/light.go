package entities

import (
	log "github.com/sirupsen/logrus"
)

type LightEntityState string
type LightEntityFeatures string
type LightEntityAttributes string
type LightEntityCommand string

const (
	OnLightEntityState  LightEntityState = "ON"
	OffLightEntityState                  = "OFF"
)

const (
	OnOffLightEntityFeatures            LightEntityFeatures = "on_off"
	ToggleLightEntityFeatures                               = "toggle"
	DimLightEntityFeatures                                  = "dim"
	ColorLightEntityFeatures                                = "color"
	ColorTemperatureLightEntityFeatures                     = "color_temperature"
)

const (
	OnLightEntityCommand     LightEntityCommand = "on"
	OffLightEntityCommand                       = "off"
	ToggleLightEntityCommand                    = "toggle"
)

const (
	StateLightEntityAttribute            LightEntityAttributes = "state"
	HueLightEntityAttribute                                    = "hue"
	SaturationLightEntityAttribute                             = "saturation"
	BrightnessLightEntityAttribute                             = "brightness"
	ColorTemperatureLightEntityAttribute                       = "color_temperature"
)

type LightEntity struct {
	Entity
	Commands map[string]func(LightEntity, map[string]interface{}) int `json:"-"`
}

func NewLightEntity(id string, name LanguageText, area string) *LightEntity {
	log.WithFields(log.Fields{
		"ID":   id,
		"Name": name,
		"Area": area,
	}).Debug(("Create new LightEntity"))

	lightEntity := LightEntity{}
	lightEntity.Id = id
	lightEntity.Name = name
	lightEntity.Area = area

	lightEntity.EntityType.Type = "light"

	lightEntity.Commands = make(map[string]func(LightEntity, map[string]interface{}) int)
	lightEntity.Attributes = make(map[string]interface{})

	return &lightEntity
}

// Register a function for the Entity command
// Based on the Feature, the correct Attributes will be added
func (e *LightEntity) AddFeature(feature LightEntityFeatures) {
	e.Features = append(e.Features, feature)

	// Add Attributes based on enabled features
	// https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_light.md
	switch feature {
	case OnOffLightEntityFeatures, ToggleLightEntityFeatures:
		e.AddAttribute(string(StateLightEntityAttribute), OffLightEntityState)

	case ColorLightEntityFeatures:
		e.AddAttribute(string(HueLightEntityAttribute), 0)
		e.AddAttribute(string(SaturationLightEntityAttribute), 0)

	case DimLightEntityFeatures:
		e.AddAttribute(string(BrightnessLightEntityAttribute), 0)

	case ColorTemperatureLightEntityFeatures:
		e.AddAttribute(string(ColorTemperatureLightEntityAttribute), 0)

	}
}

// Register a function for the Entity command
func (e *LightEntity) AddCommand(command LightEntityCommand, function func(LightEntity, map[string]interface{}) int) {
	e.Commands[string(command)] = function

}

// Call the registred function for this entity_command
func (e *LightEntity) HandleCommand(cmd_id string, params map[string]interface{}) int {
	if e.Commands[cmd_id] != nil {
		return e.Commands[cmd_id](*e, params)
	}

	return 404
}

func (e *LightEntity) HasAttribute(attribute LightEntityAttributes) bool {
	_, ok := e.Attributes[string(attribute)]

	return ok
}
