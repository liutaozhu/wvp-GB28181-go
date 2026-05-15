package xml

import "encoding/xml"

// Keepalive XML for device heartbeat
type Keepalive struct {
	XMLName  xml.Name `xml:"Notify"`
	CmdType  string   `xml:"CmdType"`
	SN       int      `xml:"SN"`
	DeviceID string   `xml:"DeviceID"`
	Status   string   `xml:"Status"`
}

func MarshalKeepalive(sn int, deviceID string) ([]byte, error) {
	ka := Keepalive{
		CmdType:  "Keepalive",
		SN:       sn,
		DeviceID: deviceID,
		Status:   "OK",
	}
	return xml.Marshal(ka)
}
