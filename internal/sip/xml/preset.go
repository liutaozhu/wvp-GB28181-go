package xml

import "encoding/xml"

// PresetQuery for querying preset positions
type PresetQuery struct {
	XMLName  xml.Name `xml:"Query"`
	CmdType  string   `xml:"CmdType"`
	SN       int      `xml:"SN"`
	DeviceID string   `xml:"DeviceId"`
}

// PresetResponse from device
type PresetResponse struct {
	XMLName  xml.Name      `xml:"Response"`
	CmdType  string        `xml:"CmdType"`
	SN       int           `xml:"SN"`
	DeviceID string        `xml:"DeviceId"`
	ItemList []PresetItem `xml:"Item"`
}

// PresetItem represents a single preset
type PresetItem struct {
	PresetID  int    `xml:"PresetID"`
	PresetName string `xml:"PresetName"`
}

// MarshalPresetQuery creates XML for preset query
func MarshalPresetQuery(sn int, deviceID string) ([]byte, error) {
	query := PresetQuery{
		CmdType:  "PresetQuery",
		SN:       sn,
		DeviceID: deviceID,
	}
	return xml.Marshal(query)
}
