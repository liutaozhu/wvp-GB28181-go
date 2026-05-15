package model

// DeviceAlarm maps to wvp_device_alarm table
type DeviceAlarm struct {
	ID             uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID       string `gorm:"column:device_id;type:varchar(50);not null" json:"deviceId"`
	ChannelID      string `gorm:"column:channel_id;type:varchar(50);not null" json:"channelId"`
	AlarmPriority  string `gorm:"column:alarm_priority;type:varchar(50)" json:"alarmPriority"`
	AlarmMethod    string `gorm:"column:alarm_method;type:varchar(50)" json:"alarmMethod"`
	AlarmTime      string `gorm:"column:alarm_time;type:varchar(50)" json:"alarmTime"`
	AlarmDescription string `gorm:"column:alarm_description;type:varchar(255)" json:"alarmDescription"`
	Longitude      float64 `json:"longitude"`
	Latitude       float64 `json:"latitude"`
	AlarmType      string `gorm:"column:alarm_type;type:varchar(50)" json:"alarmType"`
	CreateTime     string `gorm:"column:create_time;type:varchar(50);not null" json:"createTime"`
}

func (DeviceAlarm) TableName() string {
	return "wvp_device_alarm"
}
