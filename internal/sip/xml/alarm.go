package xml

import "encoding/xml"

// AlarmNotify XML for alarm notifications
type AlarmNotify struct {
	XMLName         xml.Name `xml:"Notify"`
	CmdType         string   `xml:"CmdType"`
	SN              int      `xml:"SN"`
	DeviceID        string   `xml:"DeviceID"`
	AlarmPriority   string   `xml:"AlarmPriority"`
	AlarmMethod     string   `xml:"AlarmMethod"`
	AlarmTime       string   `xml:"AlarmTime"`
	AlarmType       string   `xml:"AlarmType"`
	Longitude       float64  `xml:"Longitude,omitempty"`
	Latitude        float64  `xml:"Latitude,omitempty"`
	Info            string   `xml:"Info,omitempty"`
}

// MarshalAlarmNotify creates XML for alarm notification
func MarshalAlarmNotify(sn int, deviceID, alarmType string) ([]byte, error) {
	notify := AlarmNotify{
		CmdType:   "Alarm",
		SN:        sn,
		DeviceID:  deviceID,
		AlarmType: alarmType,
	}
	return xml.Marshal(notify)
}
