package model

// Record maps to wvp_record table for MP4 recording info
type Record struct {
	ID         uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	App        string  `gorm:"column:app;type:varchar(255);not null" json:"app"`
	Stream     string  `gorm:"column:stream;type:varchar(255);not null" json:"stream"`
	FilePath   string  `gorm:"column:file_path;type:varchar(255)" json:"filePath"`
	Folder     string  `gorm:"column:folder;type:varchar(255)" json:"folder"`
	FileName   string  `gorm:"column:file_name;type:varchar(255)" json:"fileName"`
	URL        string  `gorm:"column:url;type:varchar(255)" json:"url"`
	Duration   float64 `gorm:"column:duration" json:"duration"`
	StartTime  string  `gorm:"column:start_time;type:varchar(50)" json:"startTime"`
	EndTime    string  `gorm:"column:end_time;type:varchar(50)" json:"endTime"`
	FileSize   int64   `gorm:"column:file_size" json:"fileSize"`
	DeviceID   string  `gorm:"column:device_id;type:varchar(50)" json:"deviceId"`
	ChannelID  string  `gorm:"column:channel_id;type:varchar(50)" json:"channelId"`
	CreateTime string  `gorm:"column:create_time;type:varchar(50);not null" json:"createTime"`
}

func (Record) TableName() string {
	return "wvp_record"
}
