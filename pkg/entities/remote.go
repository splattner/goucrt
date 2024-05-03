package entities

type RemoteEntityState EntityState
type RemoteEntityFeatures EntityFeature
type RemoteEntityAttributes EntityAttribute
type RemoteEntityCommand EntityCommand
type RemoteEntityOption EntityOption

const (
	OnRemoteEntityState  RemoteEntityState = "ON"
	OffRemoteEntityState RemoteEntityState = "OFF"
)

const (
	SendCmdRemoteEntityFeatures RemoteEntityFeatures = "send_cmd"
	OnOffRemoteEntityFeatures   RemoteEntityFeatures = "on_off"
	ToggleRemoteEntityyFeatures RemoteEntityFeatures = "toggle"
)

const (
	StateRemoteEntityAttribute RemoteEntityAttributes = "state"
)

const (
	OnRemoteEntityCommand              RemoteEntityCommand = "on"
	OffRemoteEntityCommand             RemoteEntityCommand = "off"
	SendCmdRemoteEntityCommand         RemoteEntityCommand = "send_cmd"
	SendCmdSequenceRemoteEntityCommand RemoteEntityCommand = "send_cmd_sequence"
)

const (
	SimpleCommandsRemoteEntityOption RemoteEntityOption = "simple_commands"
	ButtonMappingRemoteEntityOption  RemoteEntityOption = "button_mapping"
	UserInterfaceRemoteEntityOption  RemoteEntityOption = "user_interface"
)

type RemoteEntity struct {
	Entity
	Commands map[RemoteEntityCommand]func(RemoteEntity, map[string]interface{}) int `json:"-"`
	Options  map[RemoteEntityOption]interface{}                                     `json:"options"`
}

func NewRemoteEntity(id string, name LanguageText, area string) *RemoteEntity {

	remoteEntity := RemoteEntity{}
	remoteEntity.Id = id
	remoteEntity.Name = name
	remoteEntity.Area = area

	remoteEntity.Commands = make(map[RemoteEntityCommand]func(RemoteEntity, map[string]interface{}) int)
	remoteEntity.Attributes = make(map[string]interface{})

	remoteEntity.Options = make(map[RemoteEntityOption]interface{})

	return &remoteEntity
}

func (e *RemoteEntity) UpdateEntity(newEntity RemoteEntity) error {

	e.Name = newEntity.Name
	e.Area = newEntity.Area
	e.Commands = newEntity.Commands
	e.Features = newEntity.Features
	e.Attributes = newEntity.Attributes

	return nil
}

// Register a function for the Entity command
// Based on the Feature, the correct Attributes will be added
func (e RemoteEntity) AddFeature(feature RemoteEntityFeatures) {
	e.Features = append(e.Features, feature)

	// Add Attributes based on enabled features
	// https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_remote.md
	switch feature {
	case OnOffRemoteEntityFeatures:
		e.AddAttribute(string(StateRemoteEntityAttribute), OffRemoteEntityState)

	}
}

// Register a function for the Entity command
func (e *RemoteEntity) AddCommand(command RemoteEntityCommand, function func(RemoteEntity, map[string]interface{}) int) {
	e.Commands[command] = function

}

// Call the registred function for this entity_command
func (e *RemoteEntity) HandleCommand(cmd_id string, params map[string]interface{}) int {
	if e.Commands[RemoteEntityCommand(cmd_id)] != nil {
		return e.Commands[RemoteEntityCommand(cmd_id)](*e, params)
	}

	return 404
}

// Check if an Attribute is available
func (e *RemoteEntity) HasAttribute(attribute RemoteEntityAttributes) bool {
	_, ok := e.Attributes[string(attribute)]

	return ok
}

// Update an Attribute if its available
func (e *RemoteEntity) UpdateAttribute(attribute RemoteEntityAttributes, value interface{}) {

	if e.HasAttribute(attribute) {
		e.Attributes[string(attribute)] = value
	}
}

// Add an option to the MediaPlayer Entity
func (e *RemoteEntity) AddOption(option RemoteEntityOption, value interface{}) {

	e.Options[option] = value

}
