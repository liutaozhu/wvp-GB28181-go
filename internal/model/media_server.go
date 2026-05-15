package model

// MediaServer maps to wvp_media_server table
type MediaServer struct {
	ID              string `gorm:"primaryKey;type:varchar(50)" json:"id"`
	IP              string `gorm:"type:varchar(50)" json:"ip"`
	HookIP          string `gorm:"column:hook_ip;type:varchar(50)" json:"hookIp"`
	SDPIP           string `gorm:"column:sdp_ip;type:varchar(50)" json:"sdpIp"`
	StreamIP        string `gorm:"column:stream_ip;type:varchar(50)" json:"streamIp"`
	HTTPPort        int    `gorm:"column:http_port" json:"httpPort"`
	HTTPSSLPort     int    `gorm:"column:http_ssl_port" json:"httpSSlPort"`
	RTMPPort        int    `gorm:"column:rtmp_port" json:"rtmpPort"`
	RTMPSSLPort     int    `gorm:"column:rtmp_ssl_port" json:"rtmpSSlPort"`
	FLVPort         int    `gorm:"column:flv_port" json:"flvPort"`
	FLVSSLPort      int    `gorm:"column:flv_ssl_port" json:"flvSSLPort"`
	MP4Port         int    `gorm:"column:mp4_port" json:"mp4Port"`
	WSFLVPort       int    `gorm:"column:ws_flv_port" json:"wsFlvPort"`
	WSFLVSSLPort    int    `gorm:"column:ws_flv_ssl_port" json:"wsFlvSSLPort"`
	RTSPPort        int    `gorm:"column:rtsp_port" json:"rtspPort"`
	RTSPSSLPort     int    `gorm:"column:rtsp_ssl_port" json:"rtspSSLPort"`
	RTPProxyPort    int    `gorm:"column:rtp_proxy_port" json:"rtpProxyPort"`
	JTTProxyPort    int    `gorm:"column:jtt_proxy_port" json:"jttProxyPort"`
	AutoConfig      bool   `gorm:"column:auto_config" json:"autoConfig"`
	Secret          string `gorm:"type:varchar(255)" json:"secret"`
	HookAliveInterval float64 `gorm:"column:hook_alive_interval" json:"hookAliveInterval"`
	RTPEnable       bool   `gorm:"column:rtp_enable" json:"rtpEnable"`
	Status          bool   `json:"status"`
	RTPPortRange    string `gorm:"column:rtp_port_range;type:varchar(50)" json:"rtpPortRange"`
	SendRTPPortRange string `gorm:"column:send_rtp_port_range;type:varchar(50)" json:"sendRtpPortRange"`
	RecordAssistPort int   `gorm:"column:record_assist_port" json:"recordAssistPort"`
	DefaultServer   bool   `gorm:"column:default_server" json:"defaultServer"`
	RecordDay       int    `gorm:"column:record_day" json:"recordDay"`
	RecordPath      string `gorm:"column:record_path;type:varchar(255)" json:"recordPath"`
	Type            string `gorm:"type:varchar(50)" json:"type"` // zlm, abl
	TranscodeSuffix string `gorm:"column:transcode_suffix;type:varchar(50)" json:"transcodeSuffix"`
	ServerID        string `gorm:"column:server_id;type:varchar(50)" json:"serverId"`
	CreateTime      string `gorm:"type:varchar(50)" json:"createTime"`
	UpdateTime      string `gorm:"type:varchar(50)" json:"updateTime"`
	LastKeepaliveTime string `gorm:"column:last_keepalive_time;type:varchar(50)" json:"lastKeepaliveTime"`
}

func (MediaServer) TableName() string {
	return "wvp_media_server"
}
