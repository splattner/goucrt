package entities

type SwitchEntityState string
type SwitchEntityFeatures string
type SwitchEntityAttributes string
type SwitchEntityCommand string

const (
	OnSwitchtEntityState  SwitchEntityState = "ON"
	OffSwitchtEntityState                   = "OFF"
)

const (
	OnOffSwitchEntityyFeatures  SwitchEntityFeatures = "on_off"
	ToggleSwitchEntityyFeatures                      = "toggle"
)

const (
	StateSwitchEntityyAttribute SwitchEntityAttributes = "state"
)

const (
	OnSwitchEntityCommand     SwitchEntityCommand = "on"
	OffSwitchEntityCommand                        = "off"
	ToggleSwitchEntityCommand                     = "toggle"
)

type SwitchsEntity struct {
	Entity
	Commands map[string]func(SwitchsEntity, map[string]interface{}) int `json:"-"`
}

func NewSwitchEntity(id string, name LanguageText, area string) *SwitchsEntity {

	switchEntity := SwitchsEntity{}
	switchEntity.Id = id
	switchEntity.Name = name
	switchEntity.Area = area

	switchEntity.EntityType.Type = "switch"

	switchEntity.Commands = make(map[string]func(SwitchsEntity, map[string]interface{}) int)
	switchEntity.Attributes = make(map[string]interface{})

	return &switchEntity
}

// Register a function for the Entity command
// Based on the Feature, the correct Attributes will be added
func (e *SwitchsEntity) AddFeature(feature SwitchEntityFeatures) {
	e.Features = append(e.Features, feature)

	// Add Attributes based on enabled features
	// https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_switch.md
	switch feature {
	case OnOffSwitchEntityyFeatures, ToggleSwitchEntityyFeatures:
		e.AddAttribute(string(StateSwitchEntityyAttribute), OffSwitchtEntityState)

	}
}

// Register a function for the Entity command
func (e *SwitchsEntity) AddCommand(command LightEntityCommand, function func(SwitchsEntity, map[string]interface{}) int) {
	e.Commands[string(command)] = function

}

// Call the registred function for this entity_command
func (e *SwitchsEntity) HandleCommand(cmd_id string, params map[string]interface{}) int {
	if e.Commands[cmd_id] != nil {
		return e.Commands[cmd_id](*e, params)
	}

	return 404
}
