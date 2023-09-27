package deconz

import (
	log "github.com/sirupsen/logrus"
)

type Deconz struct {
	host          string
	port          int
	websocketport int
	apikey        string

	// Array with all lights, groups, sensors
	allDeconzDevices []*DeconzDevice

	controlChannel chan string

	handleDeviceDiscoveredFunc func(*DeconzDevice)
	handleDeviceRemoveFunc     func(*DeconzDevice)
}

// Create a new DeCONZ client
func NewDeconz(host string, port int, websocketport int, apikey string) *Deconz {

	deconz := Deconz{}
	deconz.host = host
	deconz.port = port
	deconz.websocketport = websocketport
	deconz.apikey = apikey

	deconz.controlChannel = make(chan string)

	return &deconz
}

// Set the function that get called when a new Deconz Device is discovered
func (d *Deconz) SetDeviceDiscoveredHandler(f func(*DeconzDevice)) {
	d.handleDeviceDiscoveredFunc = f
}

// Set the function that get called when a Deconz Device is no longer available anymore
func (d *Deconz) SetDeviceRemoveHandler(f func(*DeconzDevice)) {
	d.handleDeviceRemoveFunc = f
}

// Add a new device if not already available
// Call handleDeviceDiscovered function
func (d *Deconz) addDevice(newDevice *DeconzDevice) {

	for _, d := range d.allDeconzDevices {
		if d.GetID() == newDevice.GetID() && d.Type == newDevice.Type {
			log.WithFields(log.Fields{
				"ID":   newDevice.GetID,
				"Type": newDevice.Type,
			}).Debug("Device already available")
			return
		}
	}
	log.WithFields(log.Fields{
		"ID":   newDevice.GetID(),
		"Type": newDevice.Type,
	}).Debug("Add Device and call handleDeviceDiscovered Func")
	d.allDeconzDevices = append(d.allDeconzDevices, newDevice)

	if d.handleDeviceDiscoveredFunc != nil {
		d.handleDeviceDiscoveredFunc(newDevice)
	}

}

// Check all devices existing if they were still discovered
// Otherwise remove it
// Call handleDeviceRemoveFunc
func (d *Deconz) removeDevice(allDevices interface{}) {

	var toRemove []*DeconzDevice

	switch allDevices := allDevices.(type) {
	case []DeconzSensor:

		// Loop trought the existing devices
		for _, dd := range d.allDeconzDevices {
			if dd.Type == SensorDeconzDeviceType {
				// Check if in discovered device
				remove := true
				for _, device := range allDevices {

					// Compare by id
					if dd.GetID() == device.ID {
						//Still available, so don't remove
						remove = false

						break
					}
				}

				if remove {
					toRemove = append(toRemove, dd)
				}
			}
		}

	case []DeconzLight:
		// Loop trought the existing devices
		for _, dd := range d.allDeconzDevices {
			if dd.Type == LightDeconzDeviceType {
				// Check if in discovered device
				remove := true
				for _, device := range allDevices {

					// Compare by id
					if dd.GetID() == device.ID {
						// Still available, so don't remove
						remove = false

						break
					}
				}
				if remove {
					toRemove = append(toRemove, dd)
				}
			}
		}

	case []DeconzGroup:
		// Loop trought the existing devices
		for _, dd := range d.allDeconzDevices {
			if dd.Type == GroupDeconzDeviceType {
				// Check if in discovered device
				remove := true
				for _, device := range allDevices {

					// Compare by id
					if dd.GetID() == device.ID {
						// Still available, so don't remove
						remove = false
						break
					}

				}
				if remove {
					toRemove = append(toRemove, dd)
				}

			}
		}
	}

	// Finaly remote thos who are not needed anymore and call device removed handler
	for ix, device := range toRemove {
		d.allDeconzDevices[ix] = d.allDeconzDevices[len(d.allDeconzDevices)-1] // Copy last element to index i.
		d.allDeconzDevices[len(d.allDeconzDevices)-1] = nil                    // Erase last element (write zero value).
		d.allDeconzDevices = d.allDeconzDevices[:len(d.allDeconzDevices)-1]    // Truncate slice.

		if d.handleDeviceRemoveFunc != nil {
			d.handleDeviceRemoveFunc(device)
		}
	}

}
