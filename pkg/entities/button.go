package entities

import "log"

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
	Commands map[string]func(ButtonEntity) `json:"-"`
}

func NewButtonEntity(id string, name LanguageText, area string) *ButtonEntity {

	buttonEntity := ButtonEntity{}
	buttonEntity.Id = id
	buttonEntity.Name = name
	buttonEntity.Area = area

	buttonEntity.EntityType.Type = "button"

	buttonEntity.Commands = make(map[string]func(ButtonEntity))

	// PressButtonEntityyFeatures is always present even if not specified
	// https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_button.md
	buttonEntity.AddFeature(PressButtonEntityFeatures)

	buttonEntity.AddAttribute(string(StateEntityAttribute), AvailableButtonEntityState)

	return &buttonEntity
}

func (e *ButtonEntity) AddFeature(feature ButtonEntityFeatures) {
	e.Features = append(e.Features, feature)

}

func (e *ButtonEntity) AddCommand(command ButtonEntityCommand, function func(ButtonEntity)) {
	e.Commands[string(command)] = function
}

func (e *ButtonEntity) HandleCommand(req *EntityCommandReq) {
	log.Println("Handle Command in Button Entity")

	if e.Commands[req.MsgData.CmdId] != nil {
		e.Commands[req.MsgData.CmdId](*e)
	}

}
