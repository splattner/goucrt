package entities

type ButtonEntityState string
type ButtonEntityFeatures string
type ButtonEntityAttribute string
type ButtonEntityCommand string

const (
	AvailableButtonEntityState EntityState = "AVAILABLE"
)

const (
	PressButtonEntityFeatures ButtonEntityFeatures = "press"
)

const (
	PushButtonEntityCommand ButtonEntityCommand = "push"
)

const (
	StateEntityAttribute ButtonEntityAttribute = "state"
)

type ButtonEntity struct {
	Entity
	Commands map[string]func(ButtonEntity) int `json:"-"`
}

func NewButtonEntity(id string, name LanguageText, area string) *ButtonEntity {

	buttonEntity := ButtonEntity{}
	buttonEntity.Id = id
	buttonEntity.Name = name
	buttonEntity.Area = area

	buttonEntity.EntityType.Type = "button"

	buttonEntity.Commands = make(map[string]func(ButtonEntity) int)
	buttonEntity.Attributes = make(map[string]interface{})

	// PressButtonEntityyFeatures is always present even if not specified
	// https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_button.md
	buttonEntity.AddFeature(PressButtonEntityFeatures)

	buttonEntity.AddAttribute(string(StateEntityAttribute), AvailableButtonEntityState)

	return &buttonEntity
}

// Add a feature to this Button
func (e *ButtonEntity) AddFeature(feature ButtonEntityFeatures) {
	e.Features = append(e.Features, feature)

}

// Register a function for the Entity command
func (e *ButtonEntity) AddCommand(command ButtonEntityCommand, function func(ButtonEntity) int) {
	e.Commands[string(command)] = function
}

// Map a Light EntityCommand to a function call without params
func (e *ButtonEntity) MapCommand(command ButtonEntityCommand, f func() error) {

	e.AddCommand(command, func(entity ButtonEntity) int {

		if err := f(); err != nil {
			return 404
		}
		return 200
	})
}

// Call the registred function for this entity_command
func (e *ButtonEntity) HandleCommand(cmd_id string) int {

	if e.Commands[cmd_id] != nil {
		return e.Commands[cmd_id](*e)
	}

	return 404

}
