package entities

type SensorEntityState string
type SensorEntityFeatures string
type SensorEntityAttributes string
type SensorEntityCommand string
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
	CustomSensorDeviceClass     SensorDeviceClass = "custom"
	BatterySensorDeviceClass    SensorDeviceClass = "battery"
	CurrentSensorDeviceClass    SensorDeviceClass = "current"
	EnegrySensorDeviceClass     SensorDeviceClass = "energy"
	HumiditySensorDeviceClass   SensorDeviceClass = "humidity"
	PowerSensorDeviceClass      SensorDeviceClass = "power"
	TemperaturSensorDeviceClass SensorDeviceClass = "temperatur"
	VoltageSensorDeviceClass    SensorDeviceClass = "voltage"
)

type SensorEntity struct {
	Entity
	DeviceClass SensorDeviceClass
}

func NewSensorEntity(id string, name LanguageText, area string, deviceClass SensorDeviceClass) *SensorEntity {

	sensorEntity := SensorEntity{}
	sensorEntity.Id = id
	sensorEntity.Name = name
	sensorEntity.Area = area

	sensorEntity.DeviceClass = deviceClass

	sensorEntity.EntityType.Type = "sensor"

	sensorEntity.Attributes = make(map[string]interface{})

	sensorEntity.AddAttribute("state", OnSensorEntityState)
	sensorEntity.AddAttribute("value", 0)
	sensorEntity.AddAttribute("unit", "")

	switch sensorEntity.DeviceClass {
	case BatterySensorDeviceClass:
		sensorEntity.Attributes["unit"] = "%"
	case CurrentSensorDeviceClass:
		sensorEntity.Attributes["unit"] = "A"
	case EnegrySensorDeviceClass:
		sensorEntity.Attributes["unit"] = "kWh"
	case HumiditySensorDeviceClass:
		sensorEntity.Attributes["unit"] = "%"
	case PowerSensorDeviceClass:
		sensorEntity.Attributes["unit"] = "W"
	case TemperaturSensorDeviceClass:
		sensorEntity.Attributes["unit"] = "°C"
	case VoltageSensorDeviceClass:
		sensorEntity.Attributes["unit"] = "V"

	}

	return &sensorEntity
}

func (e *SensorEntity) UpdateEntity(newEntity SensorEntity) error {

	e.Name = newEntity.Name
	e.Area = newEntity.Area
	e.Attributes["unit"] = newEntity.Attributes["unit"]

	return nil
}
