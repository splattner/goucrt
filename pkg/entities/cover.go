package entities

type CoverEntityState string
type CoverEntityFeatures string
type CoverEntityAttributes string
type CoverEntityCommand string

const (
	OpeningCoverEntityState CoverEntityState = "OPENING"
	OpenCoverEntityState                     = "OPEN"
	ClosingCoverEntityState                  = "CLOSING"
	CloseCoverEntityState                    = "CLOSED"
)

const (
	OpenCoverEntityFeatures         CoverEntityFeatures = "open"
	CloseCoverEntityFeatures                            = "close"
	StopCoverEntityFeatures                             = "stop"
	PositionCoverEntityFeatures                         = "position"
	TiltCoverEntityFeatures                             = "tilt"
	TiltStopCoverEntityFeatures                         = "tilt_stop"
	TiltPositionCoverEntityFeatures                     = "tilt_position"
)

const (
	OpenCoverEntityCommand     CoverEntityCommand = "open"
	CloseCoverEntityCommand                       = "close"
	StopCoverEntityyommand                        = "stop"
	PositionCoverEntityCommand                    = "position"
	TiltCoverEntityCommand                        = "tilt"
	TiltUpCoverEntityCommand                      = "tilt_up"
	TiltDownCoverEntityCommand                    = "tilt_down"
	TiltStopCoverEntityCommand                    = "tilt_stop"
)

const (
	StateCoverEntityAttribute        CoverEntityAttributes = "state"
	PositionCoverEntityAttribute                           = "position"
	TiltPositionCoverEntityAttribute                       = "tilt_position"
)

type CoverEntity struct {
	Entity
	Commands map[string]func(CoverEntity, map[string]interface{}) int `json:"-"`
}

func NewCoverEntity(id string, name LanguageText, area string) *CoverEntity {

	coverEntity := CoverEntity{}
	coverEntity.Id = id
	coverEntity.Name = name
	coverEntity.Area = area

	coverEntity.EntityType.Type = "cover"

	coverEntity.Commands = make(map[string]func(CoverEntity, map[string]interface{}) int)
	coverEntity.Attributes = make(map[string]interface{})

	return &coverEntity
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
	e.Commands[string(command)] = function

}

// Call the registred function for this entity_command
func (e *CoverEntity) HandleCommand(cmd_id string, params map[string]interface{}) int {
	if e.Commands[cmd_id] != nil {
		return e.Commands[cmd_id](*e, params)
	}

	return 404
}
