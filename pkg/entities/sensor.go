package entities

type SensorEntityState string
type SensorEntityFeatures string
type SensorEntityAttributes string
type SensorEntityCommand string

const (
	OnSensorEntityState SensorEntityState = "ON"
)

const (
	StateSensorEntityyAttribute  SensorEntityAttributes = "state"
	ValueSensortEntityyAttribute                        = "value"
	UnitSSensorntityyAttribute                          = "unit"
)

type SensorEntity struct {
	Entity
}

func NewSensorEntity(id string, name LanguageText, area string) *SensorEntity {

	sensorEntity := SensorEntity{}
	sensorEntity.Id = id
	sensorEntity.Name = name
	sensorEntity.Area = area

	sensorEntity.EntityType.Type = "sensor"

	return &sensorEntity
}
