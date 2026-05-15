package xml

import "encoding/xml"

// MobilePositionNotify XML for GPS position reports
type MobilePositionNotify struct {
	XMLName    xml.Name `xml:"Notify"`
	CmdType    string   `xml:"CmdType"`
	SN         int      `xml:"SN"`
	DeviceID   string   `xml:"DeviceID"`
	Time       string   `xml:"Time"`
	Longitude  float64  `xml:"Longitude"`
	Latitude   float64  `xml:"Latitude"`
	Speed      float64  `xml:"Speed"`
	Direction  float64  `xml:"Direction"`
	Altitude   float64  `xml:"Altitude,omitempty"`
}

// MarshalMobilePositionSubscribe creates XML for mobile position subscription
func MarshalMobilePositionSubscribe(sn int, deviceID string, interval int) ([]byte, error) {
	type Subscribe struct {
		XMLName  xml.Name `xml:"Control"`
		CmdType  string   `xml:"CmdType"`
		SN       int      `xml:"SN"`
		DeviceID string   `xml:"DeviceID"`
		Interval int      `xml:"Interval"`
	}
	sub := Subscribe{
		CmdType:  "MobilePosition",
		SN:       sn,
		DeviceID: deviceID,
		Interval: interval,
	}
	return xml.Marshal(sub)
}
