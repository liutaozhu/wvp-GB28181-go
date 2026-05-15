package xml

import "encoding/xml"

// DeviceInfoQuery for querying device information
type DeviceInfoQuery struct {
	XMLName  xml.Name `xml:"Query"`
	CmdType  string   `xml:"CmdType"`
	SN       int      `xml:"SN"`
	DeviceID string   `xml:"DeviceID"`
}

// DeviceInfoResponse from device
type DeviceInfoResponse struct {
	XMLName      xml.Name `xml:"Response"`
	CmdType      string   `xml:"CmdType"`
	SN           int      `xml:"SN"`
	DeviceID     string   `xml:"DeviceId"`
	DeviceName   string   `xml:"DeviceName"`
	Manufacturer string   `xml:"Manufacturer"`
	Model        string   `xml:"Model"`
	Firmware     string   `xml:"Firmware"`
	ChannelCount int      `xml:"Channel"`
}

// MarshalDeviceInfoQuery creates XML for device info query
func MarshalDeviceInfoQuery(sn int, deviceID string) ([]byte, error) {
	query := DeviceInfoQuery{
		CmdType:  "DeviceInfo",
		SN:       sn,
		DeviceID: deviceID,
	}
	return xml.Marshal(query)
}
