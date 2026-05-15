package model

// SsrcTransaction tracks SSRC transactions (stored in Redis primarily)
type SsrcTransaction struct {
	DeviceID      string `json:"deviceId"`
	PlatformID    string `json:"platformId"`
	ChannelID     string `json:"channelId"`
	CallID        string `json:"callId"`
	App           string `json:"app"`
	Stream        string `json:"stream"`
	MediaServerID string `json:"mediaServerId"`
	SSRC          string `json:"ssrc"`
	AllocatedSSRC string `json:"allocatedSsrc"`
	Type          string `json:"type"` // PLAY, PLAYBACK, DOWNLOAD, BROADCAST, TALK

	// SIP transaction info
	CallIDStr   string `json:"callIdStr"`
	FromTag     string `json:"fromTag"`
	ToTag       string `json:"toTag"`
	ViaBranch   string `json:"viaBranch"`
	Expires     int    `json:"expires"`
	User        string `json:"user"`
	EventID     string `json:"eventId"`
	AsSender    bool   `json:"asSender"`
}

// SipTransactionInfo contains SIP transaction details
type SipTransactionInfo struct {
	CallID    string `json:"callId"`
	FromTag   string `json:"fromTag"`
	ToTag     string `json:"toTag"`
	ViaBranch string `json:"viaBranch"`
	Expires   int    `json:"expires"`
	User      string `json:"user"`
	EventID   string `json:"eventId"`
	AsSender  bool   `json:"asSender"`
}
