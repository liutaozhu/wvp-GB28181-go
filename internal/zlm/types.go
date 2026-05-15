package zlm

// ServerConfig is the ZLM server configuration response
type ServerConfig struct {
	API struct {
		Debug     string `json:"apiDebug"`
		Secret    string `json:"secret"`
		SnapRoot  string `json:"snapRoot"`
		DownLoad  string `json:"downLoad"`
		CharSet   string `json:"charSet"`
		StreamNoneReaderDelayMS string `json:"streamNoneReaderDelayMS"`
		MaxStreamWaitMS string `json:"maxStreamWaitMS"`
	} `json:"api"`
	Hook struct {
		Enable          string `json:"enable"`
		TimeoutSec      string `json:"timeoutSec"`
		OnPublish       string `json:"on_publish"`
		OnPlay          string `json:"on_play"`
		OnStreamChanged string `json:"on_stream_changed"`
		OnFlowReport    string `json:"on_flow_report"`
		OnHTTPAccess    string `json:"on_http_access"`
		OnServerStarted string `json:"on_server_started"`
		OnServerExited  string `json:"on_server_exited"`
		OnRecvRtpAux    string `json:"on_recv_rtp_aux"`
	} `json:"hook"`
	Media struct {
		FlowThreshold  string `json:"flowThreshold"`
		FlowSeconds    string `json:"flowSeconds"`
		AddMuteAudio   string `json:"addMuteAudio"`
		WaitTrackReadyMS string `json:"waitTrackReadyMS"`
	} `json:"media"`
	Protocol struct {
		AutoClose  string `json:"autoClose"`
		ParseTxt   string `json:"parseTxt"`
		Enabled    string `json:"enabled"`
	} `json:"protocol"`
}

// StreamProxyResult is the response from addStreamProxy
type StreamProxyResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Key  string `json:"key"`
}

// MediaInfo contains media stream information
type MediaInfo struct {
	App          string `json:"app"`
	Stream       string `json:"stream"`
	Schema       string `json:"schema"`
	Vhost        string `json:"vhost"`
	Online       int    `json:"online"`
	ReaderCount  int    `json:"readerCount"`
	TotalReaderCount int `json:"totalReaderCount"`
	OriginType   int    `json:"originType"`
	OriginTypeStr string `json:"originTypeStr"`
	OriginURL    string `json:"originUrl"`
	CreateSecond float64 `json:"createSecond"`
	AliveSecond  float64 `json:"aliveSecond"`
	BytesSpeed   int    `json:"bytesSpeed"`
	Tracks       []TrackInfo `json:"tracks"`
}

// TrackInfo contains track information
type TrackInfo struct {
	CodecID     int    `json:"codec_id"`
	CodecType   int    `json:"codec_type"`
	Ready       bool   `json:"ready"`
	CodecIDName string `json:"codec_id_name"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	FPS         int    `json:"fps"`
	Channels    int    `json:"channels"`
	SampleRate  int    `json:"sample_rate"`
	Bitrate     int    `json:"bitrate"`
}

// MediaListResponse is the response from getMediaList
type MediaListResponse struct {
	Code int         `json:"code"`
	Data []MediaInfo `json:"data"`
}

// OnlineMediaListResponse is the response from getOnlineMediaList
type OnlineMediaListResponse struct {
	Code int `json:"code"`
	Data struct {
		OnlineMediaList []MediaInfo `json:"onLine"`
	} `json:"data"`
}

// RtpInfo contains RTP server info
type RtpInfo struct {
	Port    int    `json:"port"`
	Stream  string `json:"stream"`
	Exists  bool   `json:"exist"`
}

// RtpInfoResponse is the response from getRtpInfo
type RtpInfoResponse struct {
	Code   int     `json:"code"`
	Exist  bool    `json:"exist"`
	RtpInfo RtpInfo `json:"data"`
}

// RtpServerResult is the response from openRtpServer
type RtpServerResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Port int    `json:"port"`
}

// RtpServer contains RTP server info
type RtpServer struct {
	Port   int    `json:"port"`
	Stream string `json:"stream"`
}

// RtpServerListResponse is the response from listRtpServer
type RtpServerListResponse struct {
	Code int         `json:"code"`
	Data []RtpServer `json:"data"`
}

// SendRtpResult is the response from startSendRtp
type SendRtpResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// HookType represents ZLM hook event types
type HookType string

const (
	HookOnPublish         HookType = "on_publish"
	HookOnPlay            HookType = "on_play"
	HookOnStreamChanged   HookType = "on_stream_changed"
	HookOnStreamNoneReader HookType = "on_stream_none_reader"
	HookOnRtpServerTimeout HookType = "on_rtp_server_timeout"
	HookOnRecordMp4       HookType = "on_record_mp4"
	HookOnRecordTs        HookType = "on_record_ts"
	HookOnFlowReport      HookType = "on_flow_report"
	HookOnHTTPAccess      HookType = "on_http_access"
	HookOnServerStarted   HookType = "on_server_started"
	HookOnServerExited    HookType = "on_server_exited"
	HookOnSendRtpStopped  HookType = "on_send_rtp_stopped"
)

// HookPublishParam is the parameter for on_publish hook
type HookPublishParam struct {
	MediaServerID string    `json:"mediaServerId"`
	App           string    `json:"app"`
	Stream        string    `json:"stream"`
	IP            string    `json:"ip"`
	Params        string    `json:"params"`
	Schema        string    `json:"schema"`
	Vhost         string    `json:"vhost"`
	ID            string    `json:"id"`
	OriginType    int       `json:"originType"`
}

// HookStreamChangedParam is the parameter for on_stream_changed hook
type HookStreamChangedParam struct {
	MediaServerID string      `json:"mediaServerId"`
	App           string      `json:"app"`
	Stream        string      `json:"stream"`
	Registrable   bool        `json:"registable"`
	Schema        string      `json:"schema"`
	Vhost         string      `json:"vhost"`
	Secret        string      `json:"secret"`
	Tracks        []TrackInfo `json:"tracks"`
	Online        bool        `json:"online"`
}

// HookStreamNoneReaderParam is the parameter for on_stream_none_reader hook
type HookStreamNoneReaderParam struct {
	MediaServerID string `json:"mediaServerId"`
	App           string `json:"app"`
	Stream        string `json:"stream"`
	Schema        string `json:"schema"`
	Vhost         string `json:"vhost"`
}

// HookRecordMp4Param is the parameter for on_record_mp4 hook
type HookRecordMp4Param struct {
	MediaServerID string    `json:"mediaServerId"`
	App           string    `json:"app"`
	Stream        string    `json:"stream"`
	FileName      string    `json:"file_name"`
	Folder        string    `json:"folder"`
	FilePath      string    `json:"file_path"`
	Time          float64   `json:"time"`
}

// HookResponse is the response to ZLM hook requests
type HookResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	// For on_play hook: addClose=1 to close the player
	AddClose int `json:"addClose,omitempty"`
	// For on_publish hook: addMuteAudio=1 to add mute audio
	AddMuteAudio int `json:"addMuteAudio,omitempty"`
}

func OKHookResponse() HookResponse {
	return HookResponse{Code: 0, Msg: "success"}
}
