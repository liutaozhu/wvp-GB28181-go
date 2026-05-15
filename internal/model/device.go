package model

// Device maps to wvp_device table
type Device struct {
	ID                              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID                        string     `gorm:"column:device_id;type:varchar(50);uniqueIndex;not null" json:"deviceId"`
	Name                            string     `gorm:"type:varchar(255)" json:"name"`
	Manufacturer                    string     `gorm:"type:varchar(255)" json:"manufacturer"`
	Model                           string     `gorm:"type:varchar(255)" json:"model"`
	Firmware                        string     `gorm:"type:varchar(255)" json:"firmware"`
	Transport                       string     `gorm:"type:varchar(50)" json:"transport"`
	StreamMode                      string     `gorm:"type:varchar(50)" json:"streamMode"`
	Online                          bool       `gorm:"column:on_line;default:false" json:"onLine"`
	IP                              string     `gorm:"type:varchar(50)" json:"ip"`
	Port                            int        `json:"port"`
	Expires                         int        `json:"expires"`
	HostAddress                     string     `gorm:"type:varchar(50)" json:"hostAddress"`
	Charset                         string     `gorm:"type:varchar(50)" json:"charset"`
	SSRCCheck                       bool       `gorm:"column:ssrc_check;default:false" json:"ssrcCheck"`
	GeoCoordSys                     string     `gorm:"type:varchar(50)" json:"geoCoordSys"`
	MediaServerID                   string     `gorm:"type:varchar(50);default:'auto'" json:"mediaServerId"`
	CustomName                      string     `gorm:"type:varchar(255)" json:"customName"`
	SDPIP                           string     `gorm:"column:sdp_ip;type:varchar(50)" json:"sdpIp"`
	LocalIP                         string     `gorm:"type:varchar(50)" json:"localIp"`
	Password                        string     `gorm:"type:varchar(255)" json:"password"`
	AsMessageChannel                bool       `gorm:"default:false" json:"asMessageChannel"`
	HeartBeatInterval               int        `json:"heartBeatInterval"`
	HeartBeatCount                  int        `json:"heartBeatCount"`
	PositionCapability              int        `json:"positionCapability"`
	ChannelCount                    int        `json:"channelCount"`
	SubscribeCycleForCatalog        int        `gorm:"default:0" json:"subscribeCycleForCatalog"`
	SubscribeCycleForMobilePosition int        `gorm:"default:0" json:"subscribeCycleForMobilePosition"`
	MobilePositionSubmissionInterval int       `gorm:"default:5" json:"mobilePositionSubmissionInterval"`
	SubscribeCycleForAlarm          int        `gorm:"default:0" json:"subscribeCycleForAlarm"`
	BroadcastPushAfterAck           bool       `gorm:"default:false" json:"broadcastPushAfterAck"`
	ServerID                        string     `gorm:"type:varchar(50)" json:"serverId"`
	CreateTime                      string     `gorm:"type:varchar(50)" json:"createTime"`
	UpdateTime                      string     `gorm:"type:varchar(50)" json:"updateTime"`

	// Non-persisted fields
	RegisterTimeStamp  int64 `gorm:"-" json:"registerTimeStamp"`
	KeepaliveTimeStamp int64 `gorm:"-" json:"keepaliveTimeStamp"`
}

func (Device) TableName() string {
	return "wvp_device"
}
