package denonavr

import (
	"fmt"
	"hash/fnv"

	"k8s.io/utils/strings/slices"
)

// Set an attribute and return true uf the attributed has changed
func (d *DenonAVR) SetAttribute(name string, value interface{}) {

	changed := d.attributes[name] != nil && d.attributes[name] == value

	d.attributes[name] = value

	if changed {
		d.callEntityChangeFunction(name, d.attributes[name])
	}

}

// Get the current Media Title
// Title of the Playing media or the current Input Function
// Return true if the d.media_title has changed
func (d *DenonAVR) getMediaTitle() string {
	media_title := ""

	if d.IsOn() {
		if slices.Contains(PLAYING_SOURCES, d.mainZoneData.InputFuncSelect) {
			// This is a source that is playing audio
			// fot the moment, also set this to the input func
			media_title = d.mainZoneData.InputFuncSelect
		} else {
			// Not a playing source
			media_title = d.mainZoneData.InputFuncSelect
		}
	}

	d.SetAttribute("media_title", media_title)
	return d.attributes["media_title"].(string)
}

// Get the current Media Title
// Title of the Playing media or the current Input Function
// Return true if the d.media_image_url has changed
func (d *DenonAVR) getMediaImageURL() string {
	media_image_url := ""

	if d.IsOn() {
		if slices.Contains(PLAYING_SOURCES, d.mainZoneData.InputFuncSelect) {
			// This is a source that is playing audio
			// fot the moment, also set this to the input func

			hash := fnv.New32a()
			hash.Write([]byte(d.getMediaTitle()))
			media_image_url = fmt.Sprintf("http://%s:%d/NetAudio/art.asp-jpg?%d", d.Host, 80, hash.Sum32())
		} else {
			media_image_url = fmt.Sprintf("http://%s:%d/", d.Host, 80) + "img/album%20art_S.png"
		}
	}

	d.SetAttribute("media_image_url", media_image_url)

	return d.attributes["media_image_url"].(string)
}
