package denonavr

import "regexp"

type ReceiverType struct {
	Type string
	Port int
}

type DescriptionType struct {
	Port int
	Url  string
}

const (

	// Receiver types
	AVR_NAME        string = "avr"
	AVR_X_NAME      string = "avr-x"
	AVR_X_2016_NAME string = "avr-x-2016"

	// General URLs
	APPCOMMAND_URL            string = "/goform/AppCommand.xml"
	APPCOMMAND0300_URL        string = "/goform/AppCommand0300.xml"
	DEVICEINFO_URL            string = "/goform/Deviceinfo.xml"
	NETAUDIOSTATUS_URL        string = "/goform/formNetAudio_StatusXml.xml"
	TUNERSTATUS_URL           string = "/goform/formTuner_TunerXml.xml"
	HDTUNERSTATUS_URL         string = "/goform/formTuner_HdXml.xml"
	COMMAND_NETAUDIO_POST_URL string = "/NetAudio/index.put.asp"
	COMMAND_PAUSE             string = "/goform/formiPhoneAppDirect.xml?NS9B"
	COMMAND_PLAY              string = "/goform/formiPhoneAppDirect.xml?NS9A"

	// Main Zone URLs
	STATUS_URL                string = "/goform/formMainZone_MainZoneXmlStatus.xml"
	MAINZONE_URL              string = "/goform/formMainZone_MainZoneXml.xml"
	COMMAND_SEL_SRC_URL       string = "/goform/formiPhoneAppDirect.xml?SI"
	COMMAND_FAV_SRC_URL       string = "/goform/formiPhoneAppDirect.xml?ZM"
	COMMAND_POWER_ON_URL      string = "/goform/formiPhoneAppPower.xml?1+PowerOn"
	COMMAND_POWER_STANDBY_URL string = "/goform/formiPhoneAppPower.xml?1+PowerStandby"
	COMMAND_VOLUME_UP_URL     string = "/goform/formiPhoneAppDirect.xml?MVUP"
	COMMAND_VOLUME_DOWN_URL   string = "/goform/formiPhoneAppDirect.xml?MVDOWN"
	COMMAND_SET_VOLUME_URL    string = "/goform/formiPhoneAppVolume.xml?1+{volume=.1f}"
	COMMAND_MUTE_ON_URL       string = "/goform/formiPhoneAppMute.xml?1+MuteOn"
	COMMAND_MUTE_OFF_URL      string = "/goform/formiPhoneAppMute.xml?1+MuteOff"
	COMMAND_SEL_SM_URL        string = "/goform/formiPhoneAppDirect.xml?MS"
	COMMAND_SET_ZST_URL       string = "/goform/formiPhoneAppDirect.xml?MN"

	// Zone 2 URLs
	STATUS_Z2_URL                string = "/goform/formZone2_Zone2XmlStatus.xml"
	COMMAND_SEL_SRC_Z2_URL       string = "/goform/formiPhoneAppDirect.xml?Z2"
	COMMAND_FAV_SRC_Z2_URL       string = "/goform/formiPhoneAppDirect.xml?Z2"
	COMMAND_POWER_ON_Z2_URL      string = "/goform/formiPhoneAppPower.xml?2+PowerOn"
	COMMAND_POWER_STANDBY_Z2_URL string = "/goform/formiPhoneAppPower.xml?2+PowerStandby"
	COMMAND_VOLUME_UP_Z2_URL     string = "/goform/formiPhoneAppDirect.xml?Z2UP"
	COMMAND_VOLUME_DOWN_Z2_URL   string = "/goform/formiPhoneAppDirect.xml?Z2DOWN"
	COMMAND_SET_VOLUME_Z2_URL    string = "/goform/formiPhoneAppVolume.xml?2+{volume=.1f}"
	COMMAND_MUTE_ON_Z2_URL       string = "/goform/formiPhoneAppMute.xml?2+MuteOn"
	COMMAND_MUTE_OFF_Z2_URL      string = "/goform/formiPhoneAppMute.xml?2+MuteOff"

	// Zone 3 URLs
	STATUS_Z3_URL                string = "/goform/formZone3_Zone3XmlStatus.xml"
	COMMAND_SEL_SRC_Z3_URL       string = "/goform/formiPhoneAppDirect.xml?Z3"
	COMMAND_FAV_SRC_Z3_URL       string = "/goform/formiPhoneAppDirect.xml?Z3"
	COMMAND_POWER_ON_Z3_URL      string = "/goform/formiPhoneAppPower.xml?3+PowerOn"
	COMMAND_POWER_STANDBY_Z3_URL string = "/goform/formiPhoneAppPower.xml?3+PowerStandby"
	COMMAND_VOLUME_UP_Z3_URL     string = "/goform/formiPhoneAppDirect.xml?Z3UP"
	COMMAND_VOLUME_DOWN_Z3_URL   string = "/goform/formiPhoneAppDirect.xml?Z3DOWN"
	COMMAND_SET_VOLUME_Z3_URL    string = "/goform/formiPhoneAppVolume.xml?3+{volume=.1f}"
	COMMAND_MUTE_ON_Z3_URL       string = "/goform/formiPhoneAppMute.xml?3+MuteOn"
	COMMAND_MUTE_OFF_Z3_URL      string = "/goform/formiPhoneAppMute.xml?3+MuteOff"
)

// AVR-X search patterns
var DEVICEINFO_AVR_X_PATTERN = regexp.MustCompile("(.*AV(C|R)-(X|S).*|.*SR500[6-9]|.*SR60(07|08|09|10|11|12|13)|.*SR70(07|08|09|10|11|12|13)|.*SR501[3-4]|.*NR1604|.*NR1710)")
var DEVICEINFO_COMMAPI_PATTERN = regexp.MustCompile("(0210|0220|0250|0300|0301)")

type ReceiverURLs map[string]string

var DENONVAR_URLS = ReceiverURLs{
	"appcommand":                  APPCOMMAND_URL,
	"appcommand0300":              APPCOMMAND0300_URL,
	"status":                      STATUS_URL,
	"mainzone":                    MAINZONE_URL,
	"deviceinfo":                  DEVICEINFO_URL,
	"netaudiostatus":              NETAUDIOSTATUS_URL,
	"tunerstatus":                 TUNERSTATUS_URL,
	"hdtunerstatus":               HDTUNERSTATUS_URL,
	"command_sel_src":             COMMAND_SEL_SRC_URL,
	"command_fav_src":             COMMAND_FAV_SRC_URL,
	"command_power_on":            COMMAND_POWER_ON_URL,
	"command_power_standby":       COMMAND_POWER_STANDBY_URL,
	"command_volume_up":           COMMAND_VOLUME_UP_URL,
	"command_volume_down":         COMMAND_VOLUME_DOWN_URL,
	"command_set_volume":          COMMAND_SET_VOLUME_URL,
	"command_mute_on":             COMMAND_MUTE_ON_URL,
	"command_mute_off":            COMMAND_MUTE_OFF_URL,
	"command_sel_sound_mode":      COMMAND_SEL_SM_URL,
	"command_netaudio_post":       COMMAND_NETAUDIO_POST_URL,
	"command_set_all_zone_stereo": COMMAND_SET_ZST_URL,
	"command_pause":               COMMAND_PAUSE,
	"command_play":                COMMAND_PLAY,
}

var ZONE2_URLS = ReceiverURLs{
	"appcommand":                  APPCOMMAND_URL,
	"appcommand0300":              APPCOMMAND0300_URL,
	"status":                      STATUS_Z2_URL,
	"mainzone":                    MAINZONE_URL,
	"deviceinfo":                  DEVICEINFO_URL,
	"netaudiostatus":              NETAUDIOSTATUS_URL,
	"tunerstatus":                 TUNERSTATUS_URL,
	"hdtunerstatus":               HDTUNERSTATUS_URL,
	"command_sel_src":             COMMAND_SEL_SRC_Z2_URL,
	"command_fav_src":             COMMAND_FAV_SRC_Z2_URL,
	"command_power_on":            COMMAND_POWER_ON_Z2_URL,
	"command_power_standby":       COMMAND_POWER_STANDBY_Z2_URL,
	"command_volume_up":           COMMAND_VOLUME_UP_Z2_URL,
	"command_volume_down":         COMMAND_VOLUME_DOWN_Z2_URL,
	"command_set_volume":          COMMAND_SET_VOLUME_Z2_URL,
	"command_mute_on":             COMMAND_MUTE_ON_Z2_URL,
	"command_mute_off":            COMMAND_MUTE_OFF_Z2_URL,
	"command_sel_sound_mode":      COMMAND_SEL_SM_URL,
	"command_netaudio_post":       COMMAND_NETAUDIO_POST_URL,
	"command_set_all_zone_stereo": COMMAND_SET_ZST_URL,
	"command_pause":               COMMAND_PAUSE,
	"command_play":                COMMAND_PLAY,
}

var ZONE3_URLS = ReceiverURLs{
	"appcommand":                  APPCOMMAND_URL,
	"appcommand0300":              APPCOMMAND0300_URL,
	"status":                      STATUS_Z3_URL,
	"mainzone":                    MAINZONE_URL,
	"deviceinfo":                  DEVICEINFO_URL,
	"netaudiostatus":              NETAUDIOSTATUS_URL,
	"tunerstatus":                 TUNERSTATUS_URL,
	"hdtunerstatus":               HDTUNERSTATUS_URL,
	"command_sel_src":             COMMAND_SEL_SRC_Z3_URL,
	"command_fav_src":             COMMAND_FAV_SRC_Z3_URL,
	"command_power_on":            COMMAND_POWER_ON_Z3_URL,
	"command_power_standby":       COMMAND_POWER_STANDBY_Z3_URL,
	"command_volume_up":           COMMAND_VOLUME_UP_Z3_URL,
	"command_volume_down":         COMMAND_VOLUME_DOWN_Z3_URL,
	"command_set_volume":          COMMAND_SET_VOLUME_Z3_URL,
	"command_mute_on":             COMMAND_MUTE_ON_Z3_URL,
	"command_mute_off":            COMMAND_MUTE_OFF_Z3_URL,
	"command_sel_sound_mode":      COMMAND_SEL_SM_URL,
	"command_netaudio_post":       COMMAND_NETAUDIO_POST_URL,
	"command_set_all_zone_stereo": COMMAND_SET_ZST_URL,
	"command_pause":               COMMAND_PAUSE,
	"command_play":                COMMAND_PLAY,
}

var AVR = ReceiverType{
	Type: AVR_NAME,
	Port: 80,
}

var AVR_X = ReceiverType{
	Type: AVR_X_NAME,
	Port: 80,
}

var AVR_X_2016 = ReceiverType{
	Type: AVR_X_2016_NAME,
	Port: 8090,
}
