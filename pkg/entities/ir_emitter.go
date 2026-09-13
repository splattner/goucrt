package entities

import "fmt"

type IrEmitterEntityState EntityState
type IrEmitterEntityFeatures EntityFeature
type IrEmitterEntityAttribute EntityAttribute
type IrEmitterEntityCommand EntityCommand
type IrEmitterEntityOption EntityOption

// IrFormat is an IR code format an emitter supports besides PRONTO, which every IR-emitter must
// support regardless. Currently the spec only defines HEX (learned codes from the Unfolded Circle
// Dock, format "<protocol>;<hex-ir-code>;<bits>;<repeat-count>" from the IRremoteESP8266 library).
type IrFormat string

const (
	OnIrEmitterEntityState IrEmitterEntityState = "ON"
)

const (
	// SendIrEmitterEntityFeatures is always present even if not specified.
	SendIrEmitterEntityFeatures IrEmitterEntityFeatures = "send_ir"
)

const (
	StateIrEmitterEntityAttribute IrEmitterEntityAttribute = "state"
)

const (
	// SendIrEmitterEntityCommand params: code (string, required), format (string, optional,
	// default PRONTO), port (string, optional - only needed with multiple outputs), repeat
	// (number, optional, default 1).
	SendIrEmitterEntityCommand IrEmitterEntityCommand = "send_ir"
	// StopIrEmitterEntityCommand params: port (string, optional).
	StopIrEmitterEntityCommand IrEmitterEntityCommand = "stop_ir"
)

const (
	// PortsIrEmitterEntityOption: []EmitterPort, default none. Can be omitted if the emitter only
	// has a single output port.
	PortsIrEmitterEntityOption IrEmitterEntityOption = "ports"
	// IrFormatsIrEmitterEntityOption: []IrFormat, default none. IR formats supported besides
	// PRONTO, which must always be supported regardless of this option.
	IrFormatsIrEmitterEntityOption IrEmitterEntityOption = "ir_formats"
)

const (
	HexIrFormat IrFormat = "HEX"
)

// EmitterPort describes one addressable output port of an IR-emitter entity, referenced by the
// send_ir/stop_ir commands' optional "port" parameter.
type EmitterPort struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// IrEmitterEntity sends IR commands in PRONTO hex format (and optionally other formats, see
// IrFormatsIrEmitterEntityOption) to integrate external IR blasters/emitters.
type IrEmitterEntity struct {
	BaseEntity
	Commands map[IrEmitterEntityCommand]func(IrEmitterEntity, map[string]interface{}) int `json:"-"`
	Options  map[IrEmitterEntityOption]interface{}                                        `json:"options,omitempty"`
}

func NewIrEmitterEntity(id string, name LanguageText, area string) *IrEmitterEntity {

	irEmitterEntity := IrEmitterEntity{}
	irEmitterEntity.Id = id
	irEmitterEntity.Name = name
	irEmitterEntity.Area = area

	irEmitterEntity.Type = "ir_emitter"

	irEmitterEntity.Commands = make(map[IrEmitterEntityCommand]func(IrEmitterEntity, map[string]interface{}) int)
	irEmitterEntity.Attributes = make(map[string]interface{})
	irEmitterEntity.Options = make(map[IrEmitterEntityOption]interface{})

	// https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_ir_emitter.md
	irEmitterEntity.AddFeature(SendIrEmitterEntityFeatures)

	return &irEmitterEntity
}

func (e *IrEmitterEntity) UpdateEntity(newEntity interface{}) error {
	updated, ok := newEntity.(IrEmitterEntity)
	if !ok {
		return fmt.Errorf("cannot update IrEmitterEntity from %T", newEntity)
	}

	e.Name = updated.Name
	e.Area = updated.Area
	e.Commands = updated.Commands
	e.Features = updated.Features
	e.Attributes = updated.Attributes
	e.Options = updated.Options

	return nil
}

// Add a feature to this IR-emitter entity
// Based on the Feature, the correct Attributes will be added
func (e *IrEmitterEntity) AddFeature(feature IrEmitterEntityFeatures) {
	e.Features = append(e.Features, feature)

	switch feature {
	case SendIrEmitterEntityFeatures:
		e.AddAttribute(string(StateIrEmitterEntityAttribute), OnIrEmitterEntityState)
	}
}

// Add an option to the IR-emitter Entity
func (e *IrEmitterEntity) AddOption(option IrEmitterEntityOption, value interface{}) {
	e.Options[option] = value
}

// Register a function for the Entity command
func (e *IrEmitterEntity) AddCommand(command IrEmitterEntityCommand, function func(IrEmitterEntity, map[string]interface{}) int) {
	e.Commands[command] = function
}

// Map an IrEmitterEntityCommand to a function call with params (send_ir, stop_ir)
func (e *IrEmitterEntity) MapCommandWithParams(command IrEmitterEntityCommand, f func(map[string]interface{}) error) {

	e.AddCommand(command, func(entity IrEmitterEntity, params map[string]interface{}) int {

		if err := f(params); err != nil {
			return 404
		}
		return 200
	})
}

// Call the registred function for this entity_command
func (e *IrEmitterEntity) HandleCommand(cmd_id string, params map[string]interface{}) int {
	if e.Commands[IrEmitterEntityCommand(cmd_id)] != nil {
		return e.Commands[IrEmitterEntityCommand(cmd_id)](*e, params)
	}

	return 404
}
