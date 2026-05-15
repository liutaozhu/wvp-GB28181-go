package xml

import "encoding/xml"

// Generic XML envelope for parsing CmdType from any message
type CmdEnvelope struct {
	XMLName xml.Name
	CmdType string `xml:"CmdType"`
}

// GetCmdType extracts the CmdType from XML
func GetCmdType(raw string) (string, error) {
	var env CmdEnvelope
	if err := xml.Unmarshal([]byte(raw), &env); err != nil {
		return "", err
	}
	return env.CmdType, nil
}

// Catalog with DeviceList for parsing incoming catalog responses
type Catalog struct {
	XMLName    xml.Name      `xml:"Response"`
	CmdType    string        `xml:"CmdType"`
	SN         int           `xml:"SN"`
	DeviceID   string        `xml:"DeviceID"`
	SumNum     int           `xml:"SumNum"`
	DeviceList []CatalogItem `xml:"Item"`
}

func ParseCatalog(raw string) (*Catalog, error) {
	var catalog Catalog
	if err := xml.Unmarshal([]byte(raw), &catalog); err != nil {
		return nil, err
	}
	return &catalog, nil
}

// DeviceInfo for parsing incoming device info responses
type DeviceInfo struct {
	XMLName      xml.Name `xml:"Response"`
	CmdType      string   `xml:"CmdType"`
	SN           int      `xml:"SN"`
	DeviceID     string   `xml:"DeviceID"`
	DeviceName   string   `xml:"DeviceName"`
	Name         string   `xml:"Name"`
	Manufacturer string   `xml:"Manufacturer"`
	Model        string   `xml:"Model"`
	Firmware     string   `xml:"Firmware"`
	Channels     int      `xml:"Channel"`
}

func ParseDeviceInfo(raw string) (*DeviceInfo, error) {
	var info DeviceInfo
	if err := xml.Unmarshal([]byte(raw), &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// Alarm for parsing incoming alarm notifications
type Alarm struct {
	XMLName        xml.Name `xml:"Notify"`
	CmdType        string   `xml:"CmdType"`
	SN             int      `xml:"SN"`
	DeviceID       string   `xml:"DeviceID"`
	AlarmPriority  string   `xml:"AlarmPriority"`
	AlarmMethod    string   `xml:"AlarmMethod"`
	AlarmTime      string   `xml:"AlarmTime"`
	AlarmType      string   `xml:"AlarmType"`
	Description    string   `xml:"Info"`
	Longitude      float64  `xml:"Longitude"`
	Latitude       float64  `xml:"Latitude"`
}

func ParseAlarm(raw string) (*Alarm, error) {
	var alarm Alarm
	if err := xml.Unmarshal([]byte(raw), &alarm); err != nil {
		return nil, err
	}
	return &alarm, nil
}

// MobilePosition for parsing incoming mobile position notifications
type MobilePosition struct {
	XMLName   xml.Name `xml:"Notify"`
	CmdType   string   `xml:"CmdType"`
	SN        int      `xml:"SN"`
	DeviceID  string   `xml:"DeviceID"`
	Longitude float64  `xml:"Longitude"`
	Latitude  float64  `xml:"Latitude"`
	Speed     float64  `xml:"Speed"`
	Direction float64  `xml:"Direction"`
	Time      string   `xml:"Time"`
}

func ParseMobilePosition(raw string) (*MobilePosition, error) {
	var pos MobilePosition
	if err := xml.Unmarshal([]byte(raw), &pos); err != nil {
		return nil, err
	}
	return &pos, nil
}

// MediaStatus for parsing incoming media status notifications
type MediaStatus struct {
	XMLName    xml.Name `xml:"Notify"`
	CmdType    string   `xml:"CmdType"`
	SN         int      `xml:"SN"`
	DeviceID   string   `xml:"DeviceID"`
	NotifyType string   `xml:"NotifyType"`
}

func ParseMediaStatus(raw string) (*MediaStatus, error) {
	var status MediaStatus
	if err := xml.Unmarshal([]byte(raw), &status); err != nil {
		return nil, err
	}
	return &status, nil
}
