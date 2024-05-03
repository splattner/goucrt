package entities

type ClimateEntityState EntityState
type ClimateEntityFeatures EntityFeature
type ClimateEntityAttributes EntityAttribute
type ClimateEntityCommand EntityCommand

const (
	OffClimateEntityState      ClimateEntityState = "OFF"
	HeatClimateEntityState     ClimateEntityState = "HEAT"
	CoolClimateEntityState     ClimateEntityState = "Cool"
	HeatCoolClimateEntityState ClimateEntityState = "HEAT_COOL"
	FanClimateEntityState      ClimateEntityState = "FAN"
	AutoClimateEntityStat      ClimateEntityState = "Auto"
)

const (
	OnOffClimateEntityFeatures                 ClimateEntityFeatures = "on_off"
	HeatClimateEntityFeatures                  ClimateEntityFeatures = "heat"
	CoolClimateEntityFeatures                  ClimateEntityFeatures = "cool"
	CurrentTemperatureClimateEntityFeatures    ClimateEntityFeatures = "current_temperature"
	TargetTemperaturClimateEntityFeatures      ClimateEntityFeatures = "target_temperatur"
	TargetTemperaturRangeClimateEntityFeatures ClimateEntityFeatures = "target_temperature_range"
	FanClimateEntityFeatures                   ClimateEntityFeatures = "fan"
)

const (
	OnClimateEntityCommand                     ClimateEntityCommand = "on"
	OffClimateEntityCommand                    ClimateEntityCommand = "off"
	HVACModeClimateEntityCommand               ClimateEntityCommand = "hvac_mode"
	TargetTemperatureClimateEntityCommand      ClimateEntityCommand = "target_temperature"
	TargetTemperatureRangeClimateEntityCommand ClimateEntityCommand = "target_temperature_range"
	FanModeClimateEntityCommand                ClimateEntityCommand = "fan_mode"
)

const (
	StateClimateEntityAttribute                 ClimateEntityAttributes = "state"
	CurrentTemperatureClimateEntityAttribute    ClimateEntityAttributes = "current_temperature"
	TargetTemperatureClimateEntityAttribute     ClimateEntityAttributes = "target_temperature"
	TargetTemperatureHighClimateEntityAttribute ClimateEntityAttributes = "target_temperature_high"
	TargetTemperatureLowClimateEntityAttribute  ClimateEntityAttributes = "target_temperature_low"
	FanModeClimateEntityAttribute               ClimateEntityAttributes = " fan_mode"
)

type ClimateEntity struct {
	Entity
	Commands map[ClimateEntityCommand]func(ClimateEntity, map[string]interface{}) int `json:"-"`
}

func NewClimateEntity(id string, name LanguageText, area string) *ClimateEntity {

	climateEntity := ClimateEntity{}
	climateEntity.Id = id
	climateEntity.Name = name
	climateEntity.Area = area

	climateEntity.EntityType.Type = "climate"

	climateEntity.Commands = make(map[ClimateEntityCommand]func(ClimateEntity, map[string]interface{}) int)
	climateEntity.Attributes = make(map[string]interface{})

	return &climateEntity
}

func (e *ClimateEntity) UpdateEntity(newEntity ClimateEntity) error {

	e.Name = newEntity.Name
	e.Area = newEntity.Area
	e.Commands = newEntity.Commands
	e.Features = newEntity.Features
	e.Attributes = newEntity.Attributes

	return nil
}

// Add a Feature to this Climat Entity
// Based on the Feature, the correct Attributes will be added
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

// Register a function for the Entity command
func (e *ClimateEntity) AddCommand(command ClimateEntityCommand, function func(ClimateEntity, map[string]interface{}) int) {
	e.Commands[command] = function

}

// Call the registred function for this entity_command
func (e *ClimateEntity) HandleCommand(cmd_id string, params map[string]interface{}) int {
	if e.Commands[ClimateEntityCommand(cmd_id)] != nil {
		return e.Commands[ClimateEntityCommand(cmd_id)](*e, params)
	}

	return 404
}
