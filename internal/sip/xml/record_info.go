package xml

import "encoding/xml"

// RecordInfoQuery for querying recorded video
type RecordInfoQuery struct {
	XMLName   xml.Name `xml:"Query"`
	CmdType   string   `xml:"CmdType"`
	SN        int      `xml:"SN"`
	DeviceID  string   `xml:"DeviceId"`
	StartTime string   `xml:"StartTime"`
	EndTime   string   `xml:"EndTime"`
	Secrecy   int      `xml:"Secrecy,omitempty"`
	Type      string   `xml:"Type,omitempty"`
}

// RecordInfoResponse from device
type RecordInfoResponse struct {
	XMLName   xml.Name      `xml:"Response"`
	CmdType   string        `xml:"CmdType"`
	SN        int           `xml:"SN"`
	DeviceID  string        `xml:"DeviceId"`
	SumNum    int           `xml:"SumNum"`
	ItemList  []RecordItem  `xml:"Item"`
}

// RecordItem represents a single recording
type RecordItem struct {
	DeviceID  string  `xml:"DeviceID"`
	Name      string  `xml:"Name"`
	FilePath  string  `xml:"FilePath"`
	Address   string  `xml:"Address"`
	StartTime string  `xml:"StartTime"`
	EndTime   string  `xml:"EndTime"`
	Secrecy   int     `xml:"Secrecy"`
	Type      string  `xml:"Type"`
	FileSize  float64 `xml:"FileSize,omitempty"`
}

// MarshalRecordInfoQuery creates XML for record info query
func MarshalRecordInfoQuery(sn int, deviceID, startTime, endTime string) ([]byte, error) {
	query := RecordInfoQuery{
		CmdType:   "RecordInfo",
		SN:        sn,
		DeviceID:  deviceID,
		StartTime: startTime,
		EndTime:   endTime,
	}
	return xml.Marshal(query)
}
