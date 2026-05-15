package service

import (
	"fmt"
	"time"

	"wvp-pro-go/internal/database"
	"wvp-pro-go/internal/event"
	"wvp-pro-go/internal/model"
	"wvp-pro-go/internal/sip"
	"wvp-pro-go/internal/zlm"

	"go.uber.org/zap"
)

type PlayService struct {
	cmd        *sip.Commander
	zlmClient  *zlm.Client
	ssrcMgr    *sip.SSRCManager
	sessionMgr *sip.SessionManager
	subscribe  *sip.Subscribe
	mediaSvc   *MediaService
	eventBus   *event.Bus
	log        *zap.Logger
}

func NewPlayService(
	cmd *sip.Commander,
	zlmClient *zlm.Client,
	ssrcMgr *sip.SSRCManager,
	sessionMgr *sip.SessionManager,
	subscribe *sip.Subscribe,
	mediaSvc *MediaService,
	eventBus *event.Bus,
	log *zap.Logger,
) *PlayService {
	return &PlayService{
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

// StartPlay starts live streaming for a device channel
func (s *PlayService) StartPlay(deviceID, channelID string) (*StreamInfo, error) {
	// 1. Get device from DB
	device, err := s.getDevice(deviceID)
	if err != nil {
		return nil, fmt.Errorf("device not found: %s", deviceID)
	}
	if !device.Online {
		return nil, fmt.Errorf("device offline: %s", deviceID)
	}

	// 2. Get channel from DB (validate existence)
	if _, err := s.getChannel(deviceID, channelID); err != nil {
		return nil, fmt.Errorf("channel not found: %s/%s", deviceID, channelID)
	}

	// 3. Select media server
	ms, err := s.selectMediaServer(device)
	if err != nil {
		return nil, fmt.Errorf("no available media server: %v", err)
	}

	// 4. Allocate SSRC
	ssrc, err := s.ssrcMgr.AllocateSSRC("play")
	if err != nil {
		return nil, fmt.Errorf("failed to allocate SSRC: %v", err)
	}

	streamID := fmt.Sprintf("%s_%s", deviceID, channelID)

	// 5. Open RTP server on ZLM
	rtpPort := ms.RTPProxyPort
	if rtpPort == 0 {
		rtpPort = 10000 // default
	}

	// Determine TCP mode based on stream mode
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

	rtpResult, err := s.mediaSvc.OpenRtpServer(ms, rtpPort, streamID, tcpMode, isUDP)
	if err != nil {
		s.ssrcMgr.ReleaseSSRC(ssrc)
		return nil, fmt.Errorf("failed to open RTP server: %v", err)
	}

	port := rtpResult.Port

	// 6. Build SDP and send INVITE
	sdpIP := device.SDPIP
	if sdpIP == "" {
		sdpIP = ms.SDPIP
	}
	if sdpIP == "" {
		sdpIP = ms.IP
	}

	// Set SSRC on builder and send INVITE
	result, err := s.cmd.PlayStreamCmd(deviceID, channelID, ssrc, sdpIP, port)
	if err != nil {
		s.mediaSvc.CloseRtpServer(ms, streamID)
		s.ssrcMgr.ReleaseSSRC(ssrc)
		return nil, fmt.Errorf("failed to send INVITE: %v", err)
	}

	// 7. Wait for response
	if result == nil {
		// Async mode - will be handled by hook callback
		s.log.Info("play initiated, waiting for stream",
			zap.String("deviceID", deviceID),
			zap.String("channelID", channelID),
			zap.String("stream", streamID),
			zap.String("ssrc", ssrc),
		)

		// Save session
		transaction := &model.SsrcTransaction{
			DeviceID:     deviceID,
			ChannelID:    channelID,
			SSRC:         ssrc,
			AllocatedSSRC: ssrc,
			Stream:       streamID,
			MediaServerID: ms.ID,
			Type:         "PLAY",
		}
		s.sessionMgr.Save(streamID, transaction)

		// Return stream info (stream URL will be available after hook)
		si := &StreamInfo{
			DeviceID:      deviceID,
			ChannelID:     channelID,
			Stream:        streamID,
			App:           "rtp",
			SSRC:          ssrc,
			MediaServerID: ms.ID,
		}
		s.buildStreamURLs(si, ms)
		return si, nil
	}

	// 8. Build stream info
	si := &StreamInfo{
		DeviceID:      deviceID,
		ChannelID:     channelID,
		Stream:        streamID,
		App:           "rtp",
		SSRC:          ssrc,
		MediaServerID: ms.ID,
	}
	s.buildStreamURLs(si, ms)

	// 9. Auto snapshot
	go func() {
		if err := s.mediaSvc.GetSnap(ms, "rtp", streamID, "snap", fmt.Sprintf("%s_%s.jpg", deviceID, channelID)); err != nil {
			s.log.Warn("snapshot failed", zap.Error(err))
		}
	}()

	return si, nil
}

// StopPlay stops live streaming
func (s *PlayService) StopPlay(deviceID, channelID string) error {
	streamID := fmt.Sprintf("%s_%s", deviceID, channelID)

	// Get session
	session, err := s.sessionMgr.Get(streamID)
	if err != nil {
		// Already stopped
		s.log.Info("stream already stopped", zap.String("stream", streamID))
		return nil
	}

	// Send BYE
	if err := s.cmd.StreamByeCmd(deviceID, channelID, session.CallID); err != nil {
		s.log.Warn("failed to send BYE", zap.Error(err))
	}

	// Clean up
	s.sessionMgr.Delete(streamID)
	s.ssrcMgr.ReleaseSSRC(session.SSRC)

	// Close RTP server if we have media server info
	if session.MediaServerID != "" {
		ms, err := s.mediaSvc.GetMediaServer(session.MediaServerID)
		if err == nil {
			s.mediaSvc.CloseRtpServer(ms, streamID)
		}
	}

	s.log.Info("stream stopped", zap.String("stream", streamID))
	return nil
}

// GetSnap gets a snapshot from a channel
func (s *PlayService) GetSnap(deviceID, channelID string) (string, error) {
	streamID := fmt.Sprintf("%s_%s.jpg", deviceID, channelID)

	// Try to get snapshot from existing stream first
	session, err := s.sessionMgr.Get(fmt.Sprintf("%s_%s", deviceID, channelID))
	if err == nil && session.MediaServerID != "" {
		ms, err := s.mediaSvc.GetMediaServer(session.MediaServerID)
		if err == nil {
			path := "snap"
			fileName := streamID
			if err := s.mediaSvc.GetSnap(ms, "rtp", session.Stream, path, fileName); err == nil {
				return fmt.Sprintf("%s/%s", path, fileName), nil
			}
		}
	}

	// If no stream, start play then snapshot
	_, err = s.StartPlay(deviceID, channelID)
	if err != nil {
		return "", err
	}

	// Wait a moment for stream to be ready
	time.Sleep(2 * time.Second)

	ms, err := s.mediaSvc.GetMediaServerForMinimumLoad(nil)
	if err != nil {
		return "", err
	}

	path := "snap"
	if err := s.mediaSvc.GetSnap(ms, "rtp", fmt.Sprintf("%s_%s", deviceID, channelID), path, streamID); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", path, streamID), nil
}

// AudioBroadcast starts audio broadcast
func (s *PlayService) AudioBroadcast(deviceID, channelID string) (*StreamInfo, error) {
	if _, err := s.getDevice(deviceID); err != nil {
		return nil, fmt.Errorf("device not found: %s", deviceID)
	}
	if _, err := s.getChannel(deviceID, channelID); err != nil {
		return nil, fmt.Errorf("channel not found: %s/%s", deviceID, channelID)
	}

	ms, err := s.mediaSvc.GetMediaServerForMinimumLoad(nil)
	if err != nil {
		return nil, err
	}

	streamID := fmt.Sprintf("%s_%s", deviceID, channelID)

	// Send audio broadcast command
	if err := s.cmd.AudioBroadcastCmd(deviceID, channelID); err != nil {
		return nil, err
	}

	si := &StreamInfo{
		DeviceID:      deviceID,
		ChannelID:     channelID,
		Stream:        streamID,
		App:           "broadcast",
		MediaServerID: ms.ID,
	}
	s.buildStreamURLs(si, ms)
	return si, nil
}

// StopAudioBroadcast stops audio broadcast
func (s *PlayService) StopAudioBroadcast(deviceID, channelID string) error {
	streamID := fmt.Sprintf("%s_%s", deviceID, channelID)
	session, _ := s.sessionMgr.Get(streamID)
	if session != nil {
		s.cmd.StreamByeCmd(deviceID, channelID, session.CallID)
		s.sessionMgr.Delete(streamID)
	}
	return nil
}

func (s *PlayService) getDevice(deviceID string) (*model.Device, error) {
	var device model.Device
	err := database.DB.Where("device_id = ?", deviceID).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (s *PlayService) getChannel(deviceID, channelID string) (*model.DeviceChannel, error) {
	var ch model.DeviceChannel
	err := database.DB.Where("device_id = ? AND gb_device_id = ?", deviceID, channelID).First(&ch).Error
	if err != nil {
		return nil, err
	}
	return &ch, nil
}

func (s *PlayService) selectMediaServer(device *model.Device) (*model.MediaServer, error) {
	if device.MediaServerID == "" || device.MediaServerID == "auto" {
		return s.mediaSvc.GetMediaServerForMinimumLoad(nil)
	}
	return s.mediaSvc.GetMediaServer(device.MediaServerID)
}

func (s *PlayService) buildStreamURLs(si *StreamInfo, ms *model.MediaServer) {
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
