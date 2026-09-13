package entities

import "fmt"

type SensorEntityState EntityState
type SensorEntityFeatures EntityFeature
type SensorEntityAttributes EntityAttribute
type SensorEntityCommand EntityCommand
type SensorDeviceClass string

const (
	OnSensorEntityState SensorEntityState = "ON"
)

const (
	StateSensorEntityyAttribute  SensorEntityAttributes = "state"
	ValueSensortEntityyAttribute SensorEntityAttributes = "value"
	UnitSSensorntityyAttribute   SensorEntityAttributes = "unit"
)

const (
	CustomSensorDeviceClass      SensorDeviceClass = "custom"
	BatterySensorDeviceClass     SensorDeviceClass = "battery"
	CurrentSensorDeviceClass     SensorDeviceClass = "current"
	EnergySensorDeviceClass      SensorDeviceClass = "energy"
	HumiditySensorDeviceClass    SensorDeviceClass = "humidity"
	PowerSensorDeviceClass       SensorDeviceClass = "power"
	TemperatureSensorDeviceClass SensorDeviceClass = "temperature"
	VoltageSensorDeviceClass     SensorDeviceClass = "voltage"
)

type SensorEntity struct {
	BaseEntity
	DeviceClass SensorDeviceClass `json:"device_class,omitempty"`
}

func NewSensorEntity(id string, name LanguageText, area string, deviceClass SensorDeviceClass) *SensorEntity {

	sensorEntity := SensorEntity{}
	sensorEntity.Id = id
	sensorEntity.Name = name
	sensorEntity.Area = area

	sensorEntity.DeviceClass = deviceClass

	sensorEntity.Type = "sensor"

	sensorEntity.Attributes = make(map[string]interface{})

	sensorEntity.AddAttribute("state", OnSensorEntityState)
	sensorEntity.AddAttribute("value", 0)
	sensorEntity.AddAttribute("unit", "")

	switch sensorEntity.DeviceClass {
	case BatterySensorDeviceClass:
		sensorEntity.Attributes["unit"] = "%"
	case CurrentSensorDeviceClass:
		sensorEntity.Attributes["unit"] = "A"
	case EnergySensorDeviceClass:
		sensorEntity.Attributes["unit"] = "kWh"
	case HumiditySensorDeviceClass:
		sensorEntity.Attributes["unit"] = "%"
	case PowerSensorDeviceClass:
		sensorEntity.Attributes["unit"] = "W"
	case TemperatureSensorDeviceClass:
		sensorEntity.Attributes["unit"] = "°C"
	case VoltageSensorDeviceClass:
		sensorEntity.Attributes["unit"] = "V"

	}

	return &sensorEntity
}

func (e *SensorEntity) UpdateEntity(newEntity interface{}) error {
	updated, ok := newEntity.(SensorEntity)
	if !ok {
		return fmt.Errorf("cannot update SensorEntity from %T", newEntity)
	}

	e.Name = updated.Name
	e.Area = updated.Area
	e.Attributes["unit"] = updated.Attributes["unit"]

	return nil
}

// A sensor has no commands; HandleCommand always reports the command as unrecognized. Exists only
// to satisfy the Entity interface.
func (e *SensorEntity) HandleCommand(cmd_id string, params map[string]interface{}) int {
	return 404
}
