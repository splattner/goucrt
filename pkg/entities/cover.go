package entities

type CoverEntityState string
type CoverEntityFeatures string
type CoverEntityAttributes string
type CoverEntityCommand string

const (
	OpeningCoverEntityState CoverEntityState = "OPENING"
	OpenCoverEntityState    CoverEntityState = "OPEN"
	ClosingCoverEntityState CoverEntityState = "CLOSING"
	CloseCoverEntityState   CoverEntityState = "CLOSED"
)

const (
	OpenCoverEntityFeatures         CoverEntityFeatures = "open"
	CloseCoverEntityFeatures        CoverEntityFeatures = "close"
	StopCoverEntityFeatures         CoverEntityFeatures = "stop"
	PositionCoverEntityFeatures     CoverEntityFeatures = "position"
	TiltCoverEntityFeatures         CoverEntityFeatures = "tilt"
	TiltStopCoverEntityFeatures     CoverEntityFeatures = "tilt_stop"
	TiltPositionCoverEntityFeatures CoverEntityFeatures = "tilt_position"
)

const (
	OpenCoverEntityCommand     CoverEntityCommand = "open"
	CloseCoverEntityCommand    CoverEntityCommand = "close"
	StopCoverEntityyommand     CoverEntityCommand = "stop"
	PositionCoverEntityCommand CoverEntityCommand = "position"
	TiltCoverEntityCommand     CoverEntityCommand = "tilt"
	TiltUpCoverEntityCommand   CoverEntityCommand = "tilt_up"
	TiltDownCoverEntityCommand CoverEntityCommand = "tilt_down"
	TiltStopCoverEntityCommand CoverEntityCommand = "tilt_stop"
)

const (
	StateCoverEntityAttribute        CoverEntityAttributes = "state"
	PositionCoverEntityAttribute     CoverEntityAttributes = "position"
	TiltPositionCoverEntityAttribute CoverEntityAttributes = "tilt_position"
)

type CoverEntity struct {
	Entity
	Commands map[CoverEntityCommand]func(CoverEntity, map[string]interface{}) int `json:"-"`
}

func NewCoverEntity(id string, name LanguageText, area string) *CoverEntity {

	coverEntity := CoverEntity{}
	coverEntity.Id = id
	coverEntity.Name = name
	coverEntity.Area = area

	coverEntity.EntityType.Type = "cover"

	coverEntity.Commands = make(map[CoverEntityCommand]func(CoverEntity, map[string]interface{}) int)
	coverEntity.Attributes = make(map[string]interface{})

	return &coverEntity
}

func (e *CoverEntity) UpdateEntity(newEntity CoverEntity) error {

	e.Name = newEntity.Name
	e.Area = newEntity.Area
	e.Commands = newEntity.Commands
	e.Features = newEntity.Features
	e.Attributes = newEntity.Attributes

	return nil

}

// Register a function for the Entity command
// Based on the Feature, the correct Attributes will be added
func (e *CoverEntity) AddFeature(feature CoverEntityFeatures) {
	e.Features = append(e.Features, feature)

	// Add Attributes based on enabled features
	// https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_cover.md
	switch feature {
	case OpenCoverEntityFeatures:
		e.AddAttribute(string(StateCoverEntityAttribute), OpenCoverEntityState)
		e.AddAttribute(string(PositionCoverEntityAttribute), 0)

	case CloseCoverEntityFeatures:
		e.AddAttribute(string(StateCoverEntityAttribute), OpenCoverEntityState)
		e.AddAttribute(string(PositionCoverEntityAttribute), 0)

	case StopCoverEntityFeatures:
		e.AddAttribute(string(StateClimateEntityAttribute), OpenCoverEntityState)

	case PositionCoverEntityFeatures:
		e.AddAttribute(string(PositionCoverEntityAttribute), 0)

	case TiltPositionCoverEntityFeatures, TiltStopCoverEntityFeatures, TiltCoverEntityFeatures:
		e.AddAttribute(string(TiltPositionCoverEntityAttribute), 0)

	}
}

// Register a function for the Entity command
func (e *CoverEntity) AddCommand(command CoverEntityCommand, function func(CoverEntity, map[string]interface{}) int) {
	e.Commands[command] = function

}

// Call the registred function for this entity_command
func (e *CoverEntity) HandleCommand(cmd_id string, params map[string]interface{}) int {
	if e.Commands[CoverEntityCommand(cmd_id)] != nil {
		return e.Commands[CoverEntityCommand(cmd_id)](*e, params)
	}

	return 404
}
