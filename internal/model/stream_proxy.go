package model

// StreamProxy maps to wvp_stream_proxy table
type StreamProxy struct {
	ID                uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Type              int    `json:"type"` // 0: 拉流代理
	App               string `gorm:"type:varchar(255)" json:"app"`
	Stream            string `gorm:"type:varchar(255)" json:"stream"`
	URL               string `gorm:"type:varchar(255)" json:"url"`
	FFmpegCMD         string `gorm:"column:ffmpeg_cmd;type:text" json:"ffmpegCmd"`
	EnableAudio       bool   `gorm:"column:enable_audio" json:"enableAudio"`
	EnableMP4         bool   `gorm:"column:enable_mp4" json:"enableMp4"`
	Enable            bool   `json:"enable"`
	Timeout           int    `json:"timeout"` // 超时时间(秒)
	Pulling           bool   `json:"pulling"` // 是否正在拉流
	EnableRemoveKey   bool   `gorm:"column:enable_remove_key" json:"enableRemoveKey"`
	RemoveKey         string `gorm:"column:remove_key;type:varchar(255)" json:"removeKey"`
	MediaServerID     string `gorm:"column:media_server_id;type:varchar(50)" json:"mediaServerId"`
	ChannelID         uint   `gorm:"column:channel_id" json:"channelId"`
	DeviceID          string `gorm:"column:device_id;type:varchar(50)" json:"deviceId"`
	Name              string `gorm:"type:varchar(255)" json:"name"`
	Description       string `gorm:"type:varchar(255)" json:"description"`
	CreateTime        string `gorm:"type:varchar(50)" json:"createTime"`
	UpdateTime        string `gorm:"type:varchar(50)" json:"updateTime"`
}

func (StreamProxy) TableName() string {
	return "wvp_stream_proxy"
}
