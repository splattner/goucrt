package entities

type SwitchEntityState EntityState
type SwitchEntityFeatures EntityFeature
type SwitchEntityAttributes EntityAttribute
type SwitchEntityCommand EntityCommand

const (
	OnSwitchtEntityState  SwitchEntityState = "ON"
	OffSwitchtEntityState SwitchEntityState = "OFF"
)

const (
	OnOffSwitchEntityyFeatures  SwitchEntityFeatures = "on_off"
	ToggleSwitchEntityyFeatures SwitchEntityFeatures = "toggle"
)

const (
	StateSwitchEntityyAttribute SwitchEntityAttributes = "state"
)

const (
	OnSwitchEntityCommand     SwitchEntityCommand = "on"
	OffSwitchEntityCommand    SwitchEntityCommand = "off"
	ToggleSwitchEntityCommand SwitchEntityCommand = "toggle"
)

type SwitchsEntity struct {
	Entity
	Commands map[SwitchEntityCommand]func(SwitchsEntity, map[string]interface{}) int `json:"-"`
}

func NewSwitchEntity(id string, name LanguageText, area string) *SwitchsEntity {

	switchEntity := SwitchsEntity{}
	switchEntity.Id = id
	switchEntity.Name = name
	switchEntity.Area = area

	switchEntity.EntityType.Type = "switch"

	switchEntity.Commands = make(map[SwitchEntityCommand]func(SwitchsEntity, map[string]interface{}) int)
	switchEntity.Attributes = make(map[string]interface{})

	return &switchEntity
}

func (e *SwitchsEntity) UpdateEntity(newEntity SwitchsEntity) error {

	e.Name = newEntity.Name
	e.Area = newEntity.Area
	e.Commands = newEntity.Commands
	e.Features = newEntity.Features
	e.Attributes = newEntity.Attributes

	return nil
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
func (e *SwitchsEntity) AddCommand(command SwitchEntityCommand, function func(SwitchsEntity, map[string]interface{}) int) {
	e.Commands[command] = function

}

func (e *SwitchsEntity) MapCommandWithParams(command SwitchEntityCommand, f func(map[string]interface{}) error) {

	e.AddCommand(command, func(entity SwitchsEntity, params map[string]interface{}) int {

		if err := f(params); err != nil {
			return 404
		}
		return 200
	})
}

func (e *SwitchsEntity) MapCommand(command SwitchEntityCommand, f func() error) {

	e.AddCommand(command, func(entity SwitchsEntity, params map[string]interface{}) int {

		if err := f(); err != nil {
			return 404
		}
		return 200
	})

}

// Call the registred function for this entity_command
func (e *SwitchsEntity) HandleCommand(cmd_id string, params map[string]interface{}) int {
	if e.Commands[SwitchEntityCommand(cmd_id)] != nil {
		return e.Commands[SwitchEntityCommand(cmd_id)](*e, params)
	}

	return 404
}
