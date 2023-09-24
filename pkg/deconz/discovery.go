package deconz

import (
	"fmt"

	deconzgroup "github.com/jurgen-kluft/go-conbee/groups"
	deconzlight "github.com/jurgen-kluft/go-conbee/lights"
	deconzsensor "github.com/jurgen-kluft/go-conbee/sensors"
	log "github.com/sirupsen/logrus"
)

func (d *Deconz) StartDiscovery(enableGroups bool) {

	log.WithField("DeCONZ Host", d.host).Info("Starting Deconz device discovery")

	if d.apikey == "" {
		log.Fatal("API Key is not set, you first need to aquire a API Key")
		return
	}

	deconzHost := d.host + ":" + fmt.Sprint(d.port)

	// Lights
	dl := deconzlight.New(deconzHost, d.apikey)
	allLights, err := dl.GetAllLights()
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
		dg := deconzgroup.New(deconzHost, d.apikey)
		allGroups, err := dg.GetAllGroups()
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
	ds := deconzsensor.New(deconzHost, d.apikey)
	allSensors, err := ds.GetAllSensors()
	if err != nil {
		log.WithError(err).Debug("Error getting all Sensors from Deconz")
	}
	log.WithField("sonsors", allSensors).Trace("Deconz Discovery")
	for _, sensor := range allSensors {
		d.sensorDiscovery(sensor)
	}
	// Remote those devices that were not discovered anymore
	d.removeDevice(allSensors)

	log.Info("Deconz, Device Discovery finished")
}

func (d *Deconz) groupsDiscovery(group deconzgroup.Group) {

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

func (d *Deconz) lightsDiscovery(light deconzlight.Light) {
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

		log.WithField("Name", light.Name).Debug("Deconz, Lights discovered")

		deconzDevice.NewDeconzDevice(d)

	}

}

func (d *Deconz) sensorDiscovery(sensor deconzsensor.Sensor) {

	log.WithFields(log.Fields{
		"Name":     sensor.Name,
		"Type":     sensor.Type,
		"UniqueID": sensor.UniqueID,
		"ID":       sensor.ID,
	}).Debug("Found new Sensor")

	// See https://dresden-elektronik.github.io/deconz-rest-doc/endpoints/sensors/#supported-sensor-types-and-states
	switch sensor.Type {
	case "ZHAOpenClose", "ZHATemperature", "ZHAHumidity", "ZHAPressure":
		deconzDevice := new(DeconzDevice)
		deconzDevice.Type = SensorDeconzDeviceType
		deconzDevice.Sensor = sensor

		log.WithField("Name", sensor.Name).Debug("Deconz, Sensor discovered")

		deconzDevice.NewDeconzDevice(d)
	}

}
