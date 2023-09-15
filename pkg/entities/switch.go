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
	Commands map[string]func(SwitchsEntity) `json:"-"`
}

func NewSwitchEntity(id string, name LanguageText, area string) *SwitchsEntity {

	switchEntity := SwitchsEntity{}
	switchEntity.Id = id
	switchEntity.Name = name
	switchEntity.Area = area

	switchEntity.EntityType.Type = "switch"

	switchEntity.Commands = make(map[string]func(SwitchsEntity))
	switchEntity.Attributes = make(map[string]interface{})

	return &switchEntity
}

func (e *SwitchsEntity) AddFeature(feature SwitchEntityFeatures) {
	e.Features = append(e.Features, feature)

	// Add Attributes based on enabled features
	// https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_switch.md
	switch feature {
	case OnOffSwitchEntityyFeatures, ToggleSwitchEntityyFeatures:
		e.AddAttribute(string(StateSwitchEntityyAttribute), OffSwitchtEntityState)

	}
}

func (e *SwitchsEntity) AddCommand(command LightEntityCommand, function func(SwitchsEntity)) {
	e.Commands[string(command)] = function

}

func (e *SwitchsEntity) HandleCommand(cmd_id string, params interface{}) {
	if e.Commands[cmd_id] != nil {
		e.Commands[cmd_id](*e)
	}
}
