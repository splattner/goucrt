package entities

import "fmt"

type SensorEntityState EntityState
type SensorEntityFeatures EntityFeature
type SensorEntityAttributes EntityAttribute
type SensorEntityCommand EntityCommand
type SensorEntityOption EntityOption
type SensorDeviceClass string

const (
	OnSensorEntityState SensorEntityState = "ON"
)

const (
	StateSensorEntityAttribute SensorEntityAttributes = "state"
	ValueSensorEntityAttribute SensorEntityAttributes = "value"
	UnitSensorEntityAttribute  SensorEntityAttributes = "unit"
)

// Deprecated: misspelled aliases kept for backward compatibility, will be removed in a future release.
const (
	// Deprecated: use StateSensorEntityAttribute instead.
	StateSensorEntityyAttribute = StateSensorEntityAttribute
	// Deprecated: use ValueSensorEntityAttribute instead.
	ValueSensortEntityyAttribute = ValueSensorEntityAttribute
	// Deprecated: use UnitSensorEntityAttribute instead.
	UnitSSensorntityyAttribute = UnitSensorEntityAttribute
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
	// BinarySensorDeviceClass: a sensor with two states ("on"/"off") in the value attribute. The
	// specific binary sensor type (e.g. "window", "motion") goes in the unit attribute - unlike the
	// other device classes, there's no single default unit to set here, so callers set it directly.
	BinarySensorDeviceClass SensorDeviceClass = "binary"
)

const (
	// CustomLabelSensorEntityOption: LanguageText. Label for a custom sensor, if device_class isn't
	// set, or to override a default device class label.
	CustomLabelSensorEntityOption SensorEntityOption = "custom_label"
	// CustomUnitSensorEntityOption: LanguageText. Unit label for a custom sensor, if device_class
	// isn't set, or to override a default unit.
	CustomUnitSensorEntityOption SensorEntityOption = "custom_unit"
	// NativeUnitSensorEntityOption: string. The sensor's native unit of measurement, to perform
	// automatic conversion. Applies to the temperature device class.
	NativeUnitSensorEntityOption SensorEntityOption = "native_unit"
	// DecimalsSensorEntityOption: int, default 0. Number of decimal places to show in the UI for a
	// numeric value; not applicable to string values.
	DecimalsSensorEntityOption SensorEntityOption = "decimals"
)

type SensorEntity struct {
	BaseEntity
	DeviceClass SensorDeviceClass                  `json:"device_class,omitempty"`
	Options     map[SensorEntityOption]interface{} `json:"options,omitempty"`
}

func NewSensorEntity(id string, name LanguageText, area string, deviceClass SensorDeviceClass) *SensorEntity {

	sensorEntity := SensorEntity{}
	sensorEntity.Id = id
	sensorEntity.Name = name
	sensorEntity.Area = area

	sensorEntity.DeviceClass = deviceClass

	sensorEntity.Type = "sensor"

	sensorEntity.Attributes = make(map[string]interface{})
	sensorEntity.Options = make(map[SensorEntityOption]interface{})

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
	e.Options = updated.Options

	return nil
}

// Add an option to the Sensor Entity
func (e *SensorEntity) AddOption(option SensorEntityOption, value interface{}) {
	e.Options[option] = value
}

// A sensor has no commands; HandleCommand always reports the command as unrecognized. Exists only
// to satisfy the Entity interface.
func (e *SensorEntity) HandleCommand(cmd_id string, params map[string]interface{}) int {
	return 404
}
