package service

import (
	"fmt"

	"wvp-pro-go/internal/database"
	"wvp-pro-go/internal/event"
	"wvp-pro-go/internal/model"
	"wvp-pro-go/internal/sip"
	"wvp-pro-go/internal/zlm"

	"go.uber.org/zap"
)

type PlaybackService struct {
	cmd        *sip.Commander
	zlmClient  *zlm.Client
	ssrcMgr    *sip.SSRCManager
	sessionMgr *sip.SessionManager
	subscribe  *sip.Subscribe
	mediaSvc   *MediaService
	eventBus   *event.Bus
	log        *zap.Logger
}

func NewPlaybackService(
	cmd *sip.Commander,
	zlmClient *zlm.Client,
	ssrcMgr *sip.SSRCManager,
	sessionMgr *sip.SessionManager,
	subscribe *sip.Subscribe,
	mediaSvc *MediaService,
	eventBus *event.Bus,
	log *zap.Logger,
) *PlaybackService {
	return &PlaybackService{
		cmd:        cmd,
		zlmClient:  zlmClient,
		ssrcMgr:    ssrcMgr,
		sessionMgr: sessionMgr,
		subscribe:  subscribe,
		mediaSvc:   mediaSvc,
		eventBus:   eventBus,
		log:        log,
	}
}

func (s *PlaybackService) StartPlayback(deviceID, channelID, startTime, endTime string) (*StreamInfo, error) {
	device, err := s.getDevice(deviceID)
	if err != nil {
		return nil, fmt.Errorf("device not found: %s", deviceID)
	}
	if !device.Online {
		return nil, fmt.Errorf("device offline: %s", deviceID)
	}

	if _, err := s.getChannel(deviceID, channelID); err != nil {
		return nil, fmt.Errorf("channel not found: %s/%s", deviceID, channelID)
	}

	ms, err := s.mediaSvc.GetMediaServerForMinimumLoad(nil)
	if err != nil {
		return nil, fmt.Errorf("no available media server: %v", err)
	}

	ssrc, err := s.ssrcMgr.AllocateSSRC("playback")
	if err != nil {
		return nil, fmt.Errorf("failed to allocate SSRC: %v", err)
	}

	streamID := fmt.Sprintf("%s_%s_%s_%s", deviceID, channelID, startTime, endTime)

	sdpIP := device.SDPIP
	if sdpIP == "" {
		sdpIP = ms.SDPIP
	}
	if sdpIP == "" {
		sdpIP = ms.IP
	}

	tcpMode := 0
	isUDP := true
	switch device.StreamMode {
	case "TCP":
		tcpMode = 1
		isUDP = false
	case "TCP-ACTIVE":
		tcpMode = 2
		isUDP = false
	}

	rtpPort := ms.RTPProxyPort
	if rtpPort == 0 {
		rtpPort = 10000
	}

	_, err = s.mediaSvc.OpenRtpServer(ms, rtpPort, streamID, tcpMode, isUDP)
	if err != nil {
		s.ssrcMgr.ReleaseSSRC(ssrc)
		return nil, fmt.Errorf("failed to open RTP server: %v", err)
	}

	_, err = s.cmd.PlaybackStreamCmd(deviceID, channelID, startTime, endTime, ssrc, sdpIP, rtpPort)
	if err != nil {
		s.mediaSvc.CloseRtpServer(ms, streamID)
		s.ssrcMgr.ReleaseSSRC(ssrc)
		return nil, fmt.Errorf("failed to send playback INVITE: %v", err)
	}

	transaction := &model.SsrcTransaction{
		DeviceID:      deviceID,
		ChannelID:     channelID,
		SSRC:          ssrc,
		AllocatedSSRC: ssrc,
		Stream:        streamID,
		MediaServerID: ms.ID,
		Type:          "PLAYBACK",
	}
	s.sessionMgr.Save(streamID, transaction)

	si := &StreamInfo{
		DeviceID:      deviceID,
		ChannelID:     channelID,
		Stream:        streamID,
		App:           "rtp",
		SSRC:          ssrc,
		MediaServerID: ms.ID,
		StartTime:     startTime,
		EndTime:       endTime,
	}
	s.buildStreamURLs(si, ms)
	return si, nil
}

func (s *PlaybackService) StopPlayback(deviceID, channelID, stream string) error {
	session, err := s.sessionMgr.Get(stream)
	if err != nil {
		return nil
	}

	s.cmd.StreamByeCmd(deviceID, channelID, session.CallID)
	s.sessionMgr.Delete(stream)
	s.ssrcMgr.ReleaseSSRC(session.SSRC)

	if session.MediaServerID != "" {
		ms, err := s.mediaSvc.GetMediaServer(session.MediaServerID)
		if err == nil {
			s.mediaSvc.CloseRtpServer(ms, stream)
		}
	}
	return nil
}

func (s *PlaybackService) PlaybackPause(stream string) error {
	session, err := s.sessionMgr.Get(stream)
	if err != nil {
		return fmt.Errorf("stream not found: %s", stream)
	}

	ms, err := s.mediaSvc.GetMediaServer(session.MediaServerID)
	if err != nil {
		return fmt.Errorf("media server not found: %s", session.MediaServerID)
	}

	_, err = s.mediaSvc.PauseRtpCheck(ms, stream)
	return err
}

func (s *PlaybackService) PlaybackResume(stream string) error {
	session, err := s.sessionMgr.Get(stream)
	if err != nil {
		return fmt.Errorf("stream not found: %s", stream)
	}

	ms, err := s.mediaSvc.GetMediaServer(session.MediaServerID)
	if err != nil {
		return fmt.Errorf("media server not found: %s", session.MediaServerID)
	}

	_, err = s.mediaSvc.ResumeRtpCheck(ms, stream)
	return err
}

func (s *PlaybackService) PlaybackSeek(stream string, seekTime int64) error {
	session, err := s.sessionMgr.Get(stream)
	if err != nil {
		return fmt.Errorf("stream not found: %s", stream)
	}

	ms, err := s.mediaSvc.GetMediaServer(session.MediaServerID)
	if err != nil {
		return fmt.Errorf("media server not found: %s", session.MediaServerID)
	}

	_, err = s.mediaSvc.SeekRtp(ms, stream, fmt.Sprintf("%d", seekTime))
	return err
}

func (s *PlaybackService) PlaybackSpeed(stream string, speed float64) error {
	session, err := s.sessionMgr.Get(stream)
	if err != nil {
		return fmt.Errorf("stream not found: %s", stream)
	}

	ms, err := s.mediaSvc.GetMediaServer(session.MediaServerID)
	if err != nil {
		return fmt.Errorf("media server not found: %s", session.MediaServerID)
	}

	_, err = s.mediaSvc.SetSpeed(ms, stream, fmt.Sprintf("%.1f", speed))
	return err
}

func (s *PlaybackService) QueryRecord(deviceID, channelID, startTime, endTime string) ([]RecordItem, error) {
	// Query record info from device via SIP MESSAGE
	var records []RecordItem
	// In a real implementation, this would send a RecordInfoQuery SIP command
	// and wait for the XML response with record list
	return records, nil
}

func (s *PlaybackService) getDevice(deviceID string) (*model.Device, error) {
	var device model.Device
	err := database.DB.Where("device_id = ?", deviceID).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (s *PlaybackService) getChannel(deviceID, channelID string) (*model.DeviceChannel, error) {
	var ch model.DeviceChannel
	err := database.DB.Where("device_id = ? AND gb_device_id = ?", deviceID, channelID).First(&ch).Error
	if err != nil {
		return nil, err
	}
	return &ch, nil
}

func (s *PlaybackService) buildStreamURLs(si *StreamInfo, ms *model.MediaServer) {
	flvPort := ms.FLVPort
	if flvPort == 0 {
		flvPort = ms.HTTPPort
	}
	si.FLV = fmt.Sprintf("http://%s:%d/%s/%s.live.flv", ms.StreamIP, flvPort, si.App, si.Stream)
	si.WSFLV = fmt.Sprintf("ws://%s:%d/%s/%s.live.flv", ms.StreamIP, flvPort, si.App, si.Stream)
	si.HLS = fmt.Sprintf("http://%s:%d/%s/%s/hls.m3u8", ms.StreamIP, flvPort, si.App, si.Stream)
	si.RTMP = fmt.Sprintf("rtmp://%s:%d/%s/%s", ms.StreamIP, ms.RTMPPort, si.App, si.Stream)
	si.RTSP = fmt.Sprintf("rtsp://%s:%d/%s/%s", ms.StreamIP, ms.RTSPPort, si.App, si.Stream)
}
