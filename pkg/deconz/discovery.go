package deconz

import (
	log "github.com/sirupsen/logrus"
)

// Get All Lights, all Groupd, all Sensors
// Add them to the available devices
func (d *Deconz) StartDiscovery(enableGroups bool) {

	log.WithField("DeCONZ Host", d.host).Info("Starting Deconz device discovery")

	if d.apikey == "" {
		log.Fatal("API Key is not set, you first need to aquire a API Key")
		return
	}

	// Lights
	allLights, err := d.GetAllLights()
	if err != nil {
		log.WithError(err).Debug("Error getting all Lights from Deconz")
	}
	log.WithField("lights", allLights).Trace("Deconz Discovery")
	for _, light := range allLights {
		d.lightsDiscovery(light)
	}
	// Remote those devices that were not discovered anymore
	d.removeDevice(allLights)

	// Groups
	if enableGroups {
		allGroups, err := d.GetAllGroups()
		if err != nil {
			log.WithError(err).Debug("Error getting all Groups from Deconz")
		}
		log.WithField("groups", allGroups).Trace("Deconz Discovery")
		for _, group := range allGroups {
			d.groupsDiscovery(group)
		}
		// Remote those devices that were not discovered anymore
		d.removeDevice(allGroups)
	}

	// Sensors
	allSensors, err := d.GetAllSensors()
	if err != nil {
		log.WithError(err).Debug("Error getting all Sensors from Deconz")
	}
	log.WithField("sesors", allSensors).Trace("Deconz Discovery")
	for _, sensor := range allSensors {
		d.sensorDiscovery(sensor)
	}
	// Remote those devices that were not discovered anymore
	d.removeDevice(allSensors)

	log.Info("Deconz, Device Discovery finished")
}

// Handle a new discovered group
func (d *Deconz) groupsDiscovery(group DeconzGroup) {

	log.WithFields(log.Fields{
		"Name": group.Name,
		"ID":   group.ID,
	}).Debug("Found new Group")

	if len(group.Lights) > 0 {

		deconzDevice := new(DeconzDevice)

		deconzDevice.Type = GroupDeconzDeviceType
		deconzDevice.Group = group

		log.WithField("Name", group.Name).Debug("Deconz, Group discovered")

		deconzDevice.NewDeconzDevice(d)
	}
}

// Handle a new discovered light
func (d *Deconz) lightsDiscovery(light DeconzLight) {
	log.WithFields(log.Fields{
		"Name":     light.Name,
		"Type":     light.Type,
		"UniqueID": light.UniqueID,
		"ID":       light.ID,
	}).Debug("Found new Light")

	if light.Type != "Configuration tool" { // filter this out
		deconzDevice := new(DeconzDevice)

		deconzDevice.Type = LightDeconzDeviceType
		deconzDevice.Light = light

		deconzDevice.NewDeconzDevice(d)

	}

}

// Handle a new discovered sensor
func (d *Deconz) sensorDiscovery(sensor DeconzSensor) {

	log.WithFields(log.Fields{
		"Name":     sensor.Name,
		"Type":     sensor.Type,
		"UniqueID": sensor.UniqueID,
		"ID":       sensor.ID,
	}).Debug("Found new Sensor")

	// See https://dresden-elektronik.github.io/deconz-rest-doc/endpoints/sensors/#supported-sensor-types-and-states
	switch sensor.Type {
	case "ZHAOpenClose", "ZHATemperature", "ZHAHumidity", "ZHAPressure":
		// Currently we only look at those
		// Todo: Maybee this can be done more generic by looking the State and figure out what we can use?
		deconzDevice := new(DeconzDevice)
		deconzDevice.Type = SensorDeconzDeviceType
		deconzDevice.Sensor = sensor

		log.WithField("Name", sensor.Name).Debug("Deconz, Sensor discovered")

		deconzDevice.NewDeconzDevice(d)
	}

}
