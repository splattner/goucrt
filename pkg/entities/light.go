package entities

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

	lightEntity := LightEntity{}
	lightEntity.Id = id
	lightEntity.Name = name
	lightEntity.Area = area

	lightEntity.EntityType.Type = "light"

	lightEntity.Commands = make(map[string]func(LightEntity, map[string]interface{}) int)
	lightEntity.Attributes = make(map[string]interface{})

	return &lightEntity
}

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

	case ColorTemperatureLightEntityFeatures:
		e.AddAttribute(string(ColorTemperatureLightEntityAttribute), 0)

	}
}

func (e *LightEntity) AddCommand(command LightEntityCommand, function func(LightEntity, map[string]interface{}) int) {
	e.Commands[string(command)] = function

}

func (e *LightEntity) HandleCommand(cmd_id string, params map[string]interface{}) int {
	if e.Commands[cmd_id] != nil {
		return e.Commands[cmd_id](*e, params)
	}

	return 404
}
