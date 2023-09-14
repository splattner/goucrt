package entities

type ClimateEntityState string
type ClimateEntityFeatures string
type ClimateEntityAttributes string
type ClimateEntityCommand string

const (
	OffClimateEntityState      ClimateEntityState = "OFF"
	HeatClimateEntityState                        = "HEAT"
	CoolClimateEntityState                        = "Cool"
	HeatCoolClimateEntityState                    = "HEAT_COOL"
	FanClimateEntityState                         = "FAN"
	AutoClimateEntityState                        = "Auto"
)

const (
	OnOffClimateEntityFeatures                 ClimateEntityFeatures = "on_off"
	HeatClimateEntityFeatures                                        = "heat"
	CoolClimateEntityFeatures                                        = "cool"
	CurrentTemperatureClimateEntityFeatures                          = "current_temperature"
	TargetTemperaturClimateEntityFeatures                            = "target_temperatur"
	TargetTemperaturRangeClimateEntityFeatures                       = "target_temperature_range"
	FanClimateEntityFeatures                                         = "fan"
)

const (
	OnClimateEntityCommand                     ClimateEntityCommand = "on"
	OffClimateEntityCommand                                         = "off"
	HVACModeClimateEntityCommand                                    = "hvac_mode"
	TargetTemperatureClimateEntityCommand                           = "target_temperature"
	TargetTemperatureRangeClimateEntityCommand                      = "target_temperature_range"
	FanModeClimateEntityCommand                                     = "fan_mode"
)

const (
	StateClimateEntityAttribute                 ClimateEntityAttributes = "state"
	CurrentTemperatureClimateEntityAttribute                            = "current_temperature"
	TargetTemperatureClimateEntityAttribute                             = "target_temperature"
	TargetTemperatureHighClimateEntityAttribute                         = "target_temperature_high"
	TargetTemperatureLowClimateEntityAttribute                          = "target_temperature_low"
	FanModeClimateEntityAttribute                                       = " fan_mode"
)

type ClimateEntity struct {
	Entity
	Commands map[string]func(ClimateEntity) `json:"-"`
}

func NewClimateEntity(id string, name LanguageText, area string) *ClimateEntity {

	climateEntity := ClimateEntity{}
	climateEntity.Id = id
	climateEntity.Name = name
	climateEntity.Area = area

	climateEntity.EntityType.Type = "climate"

	climateEntity.Commands = make(map[string]func(ClimateEntity))

	return &climateEntity
}

func (e *ClimateEntity) AddFeature(feature ClimateEntityFeatures) {
	e.Features = append(e.Features, feature)

	// Add Attributes based on enabled features
	// https://github.com/unfoldedcircle/core-api/blob/main/doc/entities/entity_climate.md
	switch feature {
	case OnOffClimateEntityFeatures, HeatClimateEntityFeatures, CoolClimateEntityFeatures, FanClimateEntityFeatures:
		e.AddAttribute(string(StateClimateEntityAttribute), OffClimateEntityState)
	case CurrentTemperatureClimateEntityFeatures:
		e.AddAttribute(string(CurrentTemperatureClimateEntityAttribute), 0)
	case TargetTemperaturClimateEntityFeatures:
		e.AddAttribute(string(TargetTemperatureClimateEntityAttribute), 0)
	case TargetTemperaturRangeClimateEntityFeatures:
		e.AddAttribute(string(TargetTemperatureHighClimateEntityAttribute), 0)
		e.AddAttribute(string(TargetTemperatureLowClimateEntityAttribute), 0)
	}
}

func (e *ClimateEntity) AddCommand(command ClimateEntityCommand, function func(ClimateEntity)) {
	e.Commands[string(command)] = function

}

func (e *ClimateEntity) HandleCommand(req *EntityCommandReq) {
	if e.Commands[req.MsgData.CmdId] != nil {
		e.Commands[req.MsgData.CmdId](*e)
	}
}
