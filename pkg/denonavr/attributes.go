package denonavr

import "reflect"

// Set an attribute and call entity Change function if changed
func (d *DenonAVR) SetAttribute(name string, value interface{}) {

	changed := d.attributes[name] != nil && !reflect.DeepEqual(d.attributes[name], value)

	d.attributes[name] = value

	if changed {
		d.callEntityChangeFunction(name, d.attributes[name])
	}

}
