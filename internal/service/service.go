package service

import (
	"go.uber.org/zap"

	"wvp-pro-go/internal/event"
	"wvp-pro-go/internal/sip"
	"wvp-pro-go/internal/zlm"
)

// StreamInfo represents stream playback/preview information
type StreamInfo struct {
	DeviceID      string  `json:"deviceId"`
	ChannelID     string  `json:"channelId"`
	Stream        string  `json:"stream"`
	App           string  `json:"app"`
	SSRC          string  `json:"ssrc"`
	MediaServerID string  `json:"mediaServerId"`
	FLV           string  `json:"flv"`
	WSFLV         string  `json:"ws_flv"`
	HLS           string  `json:"hls"`
	RTMP          string  `json:"rtmp"`
	RTSP          string  `json:"rtsp"`
	StartTime     string  `json:"startTime,omitempty"`
	EndTime       string  `json:"endTime,omitempty"`
	Progress      float64 `json:"progress,omitempty"`
	Pause         bool    `json:"pause,omitempty"`
}

// RecordItem represents a record query result
type RecordItem struct {
	DeviceID  string `json:"deviceId"`
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Sponsor   string `json:"sponsor"`
	Secrecy   int    `json:"secrecy"`
	Type      string `json:"type"`
	FileSize  float64 `json:"fileSize"`
}

// Services holds all service instances
type Services struct {
	Device      *DeviceService
	Channel     *ChannelService
	Play        *PlayService
	PTZ         *PTZService
	Playback    *PlaybackService
	Platform    *PlatformService
	Media       *MediaService
	StreamProxy *StreamProxyService
	Region      *RegionService
	Group       *GroupService
	Alarm       *AlarmService
	MobilePos   *MobilePositionService
	User        *UserService
}

// InitServices initializes all services
func InitServices(
	sipCmd *sip.Commander,
	zlmClient *zlm.Client,
	zlmServer *zlm.Server,
	ssrcMgr *sip.SSRCManager,
	sessionMgr *sip.SessionManager,
	subscribe *sip.Subscribe,
	eventBus *event.Bus,
	logger *zap.Logger,
) *Services {
	mediaSvc := NewMediaService(zlmClient, zlmServer, logger)

	return &Services{
		Device:      NewDeviceService(logger),
		Channel:     NewChannelService(logger),
		Play:        NewPlayService(sipCmd, zlmClient, ssrcMgr, sessionMgr, subscribe, mediaSvc, eventBus, logger),
		PTZ:         NewPTZService(sipCmd, logger),
		Playback:    NewPlaybackService(sipCmd, zlmClient, ssrcMgr, sessionMgr, subscribe, mediaSvc, eventBus, logger),
		Platform:    NewPlatformService(logger),
		Media:       mediaSvc,
		StreamProxy: NewStreamProxyService(zlmClient, mediaSvc, logger),
		Region:      NewRegionService(logger),
		Group:       NewGroupService(logger),
		Alarm:       NewAlarmService(logger),
		MobilePos:   NewMobilePositionService(logger),
		User:        NewUserService(logger),
	}
}
