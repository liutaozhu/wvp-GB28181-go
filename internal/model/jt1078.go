package model

// JTTerminal maps to wvp_jt_terminal table (JT/T 1078)
type JTTerminal struct {
	ID                uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	PhoneNumber       string `gorm:"column:phone_number;type:varchar(20);uniqueIndex;not null" json:"phoneNumber"`
	PlateNumber       string `gorm:"column:plate_number;type:varchar(20)" json:"plateNumber"`
	PlateColor        int    `gorm:"column:plate_color" json:"plateColor"`
	SIMCardID         string `gorm:"column:sim_card_id;type:varchar(20)" json:"simCardId"`
	TerminalID        string `gorm:"column:terminal_id;type:varchar(20)" json:"terminalId"`
	TerminalModel     string `gorm:"column:terminal_model;type:varchar(20)" json:"terminalModel"`
	ManufacturerID    string `gorm:"column:manufacturer_id;type:varchar(20)" json:"manufacturerId"`
	ProvinceID        int    `gorm:"column:province_id" json:"provinceId"`
	CityID            int    `gorm:"column:city_id" json:"cityId"`
	Online            bool   `gorm:"default:false" json:"online"`
	MediaServerID     string `gorm:"column:media_server_id;type:varchar(50)" json:"mediaServerId"`
	CreateTime        string `gorm:"type:varchar(50)" json:"createTime"`
	UpdateTime        string `gorm:"type:varchar(50)" json:"updateTime"`
}

func (JTTerminal) TableName() string {
	return "wvp_jt_terminal"
}

// JTChannel maps to wvp_jt_channel table
type JTChannel struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	PhoneNumber string `gorm:"column:phone_number;type:varchar(20);index;not null" json:"phoneNumber"`
	ChannelID   string `gorm:"column:channel_id;type:varchar(20);not null" json:"channelId"`
	ChannelName string `gorm:"column:channel_name;type:varchar(255)" json:"channelName"`
	ChannelType int    `gorm:"column:channel_type" json:"channelType"`
}

func (JTChannel) TableName() string {
	return "wvp_jt_channel"
}
