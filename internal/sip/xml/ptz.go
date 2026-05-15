package xml

import "encoding/xml"

// Control represents PTZ control command XML
type Control struct {
	XMLName      xml.Name `xml:"Control"`
	CmdType      string   `xml:"CmdType"`
	SN           int      `xml:"SN"`
	DeviceID     string   `xml:"DeviceID"`
	PTZCmd       string   `xml:"PTZCmd"`
	Info         *Info    `xml:"Info,omitempty"`
	CombineCode2 int      `xml:"CombineCode2,omitempty"`
}

// Info for front-end control
type Info struct {
	ControlPriority int `xml:"ControlPriority"`
}

// PTZCmdByte creates a PTZ command byte array from parameters
// command: 0=stop, direction codes for various movements
// horizonSpeed, verticalSpeed, zoomSpeed: 0-255
func PTZCmdByte(command byte, horizonSpeed, verticalSpeed, zoomSpeed byte) []byte {
	return []byte{command, horizonSpeed, verticalSpeed, zoomSpeed}
}

// MarshalPTZCmd creates XML for PTZ control
func MarshalPTZCmd(sn int, deviceID string, ptzCmd []byte) ([]byte, error) {
	cmd := hexEncode(ptzCmd)
	control := Control{
		CmdType:  "DeviceControl",
		SN:       sn,
		DeviceID: deviceID,
		PTZCmd:   cmd,
	}
	return xml.Marshal(control)
}

func hexEncode(b []byte) string {
	result := make([]byte, len(b)*2)
	for i, v := range b {
		result[i*2] = hexChar(v >> 4)
		result[i*2+1] = hexChar(v & 0x0f)
	}
	return string(result)
}

func hexChar(b byte) byte {
	if b < 10 {
		return '0' + b
	}
	return 'A' + b - 10
}
