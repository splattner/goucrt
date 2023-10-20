package denonavr

import "encoding/xml"

type ValueLists struct {
	Index string `xml:"index,attr"`
	Table string `xml:"table,attr"`
}

type DenonXML struct {
	XMLName          xml.Name     `xml:"item"`
	FriendlyName     string       `xml:"FriendlyName>value"`
	Power            string       `xml:"Power>value"`
	ZonePower        string       `xml:"ZonePower>value"`
	RenameZone       string       `xml:"RenameZone>value"`
	TopMenuLink      string       `xml:"TopMenuLink>value"`
	VideoSelectDisp  string       `xml:"VideoSelectDisp>value"`
	VideoSelect      string       `xml:"VideoSelect>value"`
	VideoSelectOnOff string       `xml:"VideoSelectOnOff>value"`
	VideoSelectList  []ValueLists `xml:"VideoSelectLists>value"`
	ECOModeDisp      string       `xml:"ECOModeDisp>value"`
	ECOMode          string       `xml:"ECOMode>value"`
	ECOModeList      []ValueLists `xml:"ECOModeLists>value"`
	AddSourceDisplay string       `xml:"AddSourceDisplay>value"`
	ModelId          string       `xml:"ModelId>value"`
	BrandId          string       `xml:"BrandId>value"`
	SalesArea        string       `xml:"SalesArea>value"`
	InputFuncSelect  string       `xml:"InputFuncSelect>value"`
	NetFuncSelect    string       `xml:"NetFuncSelect>value"`
	SelectSurround   string       `xml:"selectSurround>value"`
	VolumeDisplay    string       `xml:"VolumeDisplay>value"`
	MasterVolume     string       `xml:"MasterVolume>value"`
	Mute             string       `xml:"Mute>value"`
}

type DeviceInfoXML struct {
	XMLName          xml.Name `xml:"Device_Info"`
	DeviceInfoVers   string   `xml:"DeviceInfoVers"`
	CommApiVers      string   `xml:"CommApiVers"`
	BrandCode        string   `xml:"BrandCode"`
	ProductCategory  string   `xml:"ProductCategory"`
	CategoryName     string   `xml:"CategoryName"`
	ManualModelName  string   `xml:"ManualModelName"`
	DeliveryCode     string   `xml:"DeliveryCode"`
	ModelName        string   `xml:"ModelName"`
	MacAddress       string   `xml:"MacAddress"`
	UpgradeVersion   string   `xml:"UpgradeVersion"`
	ReloadDeviceInfo string   `xml:"ReloadDeviceInfo"`
	DeviceZones      string   `xml:"DeviceZones"`
}
