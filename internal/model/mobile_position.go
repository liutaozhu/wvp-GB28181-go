package model

// MobilePosition maps to wvp_device_mobile_position table
type MobilePosition struct {
	ID           uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID     string  `gorm:"column:device_id;type:varchar(50);not null" json:"deviceId"`
	ChannelID    string  `gorm:"column:channel_id;type:varchar(50);not null" json:"channelId"`
	DeviceName   string  `gorm:"column:device_name;type:varchar(255)" json:"deviceName"`
	Time         string  `gorm:"type:varchar(50)" json:"time"`
	Longitude    float64 `json:"longitude"`
	Latitude     float64 `json:"latitude"`
	Altitude     float64 `json:"altitude"`
	Speed        float64 `json:"speed"`
	Direction    float64 `json:"direction"`
	ReportSource string  `gorm:"column:report_source;type:varchar(50)" json:"reportSource"`
	CreateTime   string  `gorm:"type:varchar(50)" json:"createTime"`
}

func (MobilePosition) TableName() string {
	return "wvp_device_mobile_position"
}
