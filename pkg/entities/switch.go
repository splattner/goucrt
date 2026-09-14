package entities

import "fmt"

type SwitchEntityState EntityState
type SwitchEntityFeatures EntityFeature
type SwitchEntityAttributes EntityAttribute
type SwitchEntityCommand EntityCommand
type SwitchEntityOption EntityOption

const (
	OnSwitchEntityState  SwitchEntityState = "ON"
	OffSwitchEntityState SwitchEntityState = "OFF"
)

const (
	OnOffSwitchEntityFeatures  SwitchEntityFeatures = "on_off"
	ToggleSwitchEntityFeatures SwitchEntityFeatures = "toggle"
)

const (
	StateSwitchEntityAttribute SwitchEntityAttributes = "state"
)

// Deprecated: misspelled aliases kept for backward compatibility, will be removed in a future release.
const (
	// Deprecated: use OnSwitchEntityState instead.
	OnSwitchtEntityState = OnSwitchEntityState
	// Deprecated: use OffSwitchEntityState instead.
	OffSwitchtEntityState = OffSwitchEntityState
	// Deprecated: use OnOffSwitchEntityFeatures instead.
	OnOffSwitchEntityyFeatures = OnOffSwitchEntityFeatures
	// Deprecated: use ToggleSwitchEntityFeatures instead.
	ToggleSwitchEntityyFeatures = ToggleSwitchEntityFeatures
	// Deprecated: use StateSwitchEntityAttribute instead.
	StateSwitchEntityyAttribute = StateSwitchEntityAttribute
)

// Deprecated: use SwitchEntity instead.
type SwitchsEntity = SwitchEntity

const (
	OnSwitchEntityCommand     SwitchEntityCommand = "on"
	OffSwitchEntityCommand    SwitchEntityCommand = "off"
	ToggleSwitchEntityCommand SwitchEntityCommand = "toggle"
)

const (
	// ReadableSwitchEntityOption: bool, default true. If set to false, the current state of the
	// switch cannot be read - the switch becomes stateless and the UI may ask the user for it.
	ReadableSwitchEntityOption SwitchEntityOption = "readable"
)

type SwitchEntity struct {
	BaseEntity
	Commands map[SwitchEntityCommand]func(SwitchEntity, map[string]interface{}) int `json:"-"`
	Options  map[SwitchEntityOption]interface{}                                     `json:"options,omitempty"`
}

func NewSwitchEntity(id string, name LanguageText, area string) *SwitchEntity {

	switchEntity := SwitchEntity{}
	switchEntity.Id = id
	switchEntity.Name = name
	switchEntity.Area = area

	switchEntity.Type = "switch"

	switchEntity.Commands = make(map[SwitchEntityCommand]func(SwitchEntity, map[string]interface{}) int)
	switchEntity.Attributes = make(map[string]interface{})
	switchEntity.Options = make(map[SwitchEntityOption]interface{})

	return &switchEntity
}

// Add an option to the Switch Entity
func (e *SwitchEntity) AddOption(option SwitchEntityOption, value interface{}) {
	e.Options[option] = value
}

func (e *SwitchEntity) UpdateEntity(newEntity interface{}) error {
	updated, ok := newEntity.(SwitchEntity)
	if !ok {
		return fmt.Errorf("cannot update SwitchEntity from %T", newEntity)
	}

	e.Name = updated.Name
	e.Area = updated.Area
	e.Commands = updated.Commands
	e.Features = updated.Features
	e.Attributes = updated.Attributes
	e.Options = updated.Options

	return nil
}

// Register a function for the Entity command
// Based on the Feature, the correct Attributes will be added
func (e *SwitchEntity) AddFeature(feature SwitchEntityFeatures) {
	e.Features = append(e.Features, feature)

	// Add Attributes based on enabled features
	// https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_switch.md
	switch feature {
	case OnOffSwitchEntityFeatures, ToggleSwitchEntityFeatures:
		e.AddAttribute(string(StateSwitchEntityAttribute), OffSwitchEntityState)

	}
}

// Register a function for the Entity command
func (e *SwitchEntity) AddCommand(command SwitchEntityCommand, function func(SwitchEntity, map[string]interface{}) int) {
	e.Commands[command] = function

}

func (e *SwitchEntity) MapCommandWithParams(command SwitchEntityCommand, f func(map[string]interface{}) error) {

	e.AddCommand(command, func(entity SwitchEntity, params map[string]interface{}) int {

		if err := f(params); err != nil {
			return 404
		}
		return 200
	})
}

func (e *SwitchEntity) MapCommand(command SwitchEntityCommand, f func() error) {

	e.AddCommand(command, func(entity SwitchEntity, params map[string]interface{}) int {

		if err := f(); err != nil {
			return 404
		}
		return 200
	})

}

// Call the registred function for this entity_command
func (e *SwitchEntity) HandleCommand(cmd_id string, params map[string]interface{}) int {
	if e.Commands[SwitchEntityCommand(cmd_id)] != nil {
		return e.Commands[SwitchEntityCommand(cmd_id)](*e, params)
	}

	return 404
}
