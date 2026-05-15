package service

import (
	"fmt"

	"wvp-pro-go/internal/database"
	"wvp-pro-go/internal/model"
	"wvp-pro-go/internal/zlm"

	"go.uber.org/zap"
)

type MediaService struct {
	zlmClient *zlm.Client
	zlmServer *zlm.Server
	log       *zap.Logger
}

func NewMediaService(zlmClient *zlm.Client, zlmServer *zlm.Server, log *zap.Logger) *MediaService {
	return &MediaService{
		zlmClient: zlmClient,
		zlmServer: zlmServer,
		log:       log,
	}
}

func (s *MediaService) GetMediaServerForMinimumLoad(hasAssist *bool) (*model.MediaServer, error) {
	servers, err := s.ListMediaServers()
	if err != nil {
		return nil, err
	}

	var best *model.MediaServer
	for i := range servers {
		if !servers[i].Status {
			continue
		}
		if hasAssist != nil {
			if *hasAssist && servers[i].RTPProxyPort == 0 {
				continue
			}
			if !*hasAssist && servers[i].RTPProxyPort > 0 {
				continue
			}
		}
		if best == nil {
			best = &servers[i]
		}
	}
	return best, nil
}

func (s *MediaService) GetMediaServer(id string) (*model.MediaServer, error) {
	var ms model.MediaServer
	err := database.DB.Where("id = ?", id).First(&ms).Error
	if err != nil {
		return nil, err
	}
	return &ms, nil
}

func (s *MediaService) ListMediaServers() ([]model.MediaServer, error) {
	var servers []model.MediaServer
	err := database.DB.Find(&servers).Error
	return servers, err
}

func (s *MediaService) GetOnlineMediaServers() ([]model.MediaServer, error) {
	var servers []model.MediaServer
	err := database.DB.Where("status = ?", true).Find(&servers).Error
	return servers, err
}

func (s *MediaService) SaveMediaServer(ms *model.MediaServer) error {
	return database.DB.Save(ms).Error
}

func (s *MediaService) DeleteMediaServer(id string) error {
	return database.DB.Delete(&model.MediaServer{}, "id = ?", id).Error
}

func (s *MediaService) AutoConfigMediaServer(ms *model.MediaServer) error {
	return s.zlmServer.AutoConfig()
}

func (s *MediaService) GetStreamInfo(ms *model.MediaServer, app, stream string) (*StreamInfo, error) {
	list, err := s.zlmClient.GetMediaList("", app, stream, "", 0)
	if err != nil {
		return nil, err
	}

	si := &StreamInfo{
		Stream:        stream,
		App:           app,
		MediaServerID: ms.ID,
	}

	if len(list) > 0 {
		flvPort := ms.FLVPort
		if flvPort == 0 {
			flvPort = ms.HTTPPort
		}
		si.FLV = fmt.Sprintf("http://%s:%d/%s/%s.live.flv", ms.StreamIP, flvPort, app, stream)
		si.WSFLV = fmt.Sprintf("ws://%s:%d/%s/%s.live.flv", ms.StreamIP, flvPort, app, stream)
		si.HLS = fmt.Sprintf("http://%s:%d/%s/%s/hls.m3u8", ms.StreamIP, flvPort, app, stream)
		si.RTMP = fmt.Sprintf("rtmp://%s:%d/%s/%s", ms.StreamIP, ms.RTMPPort, app, stream)
		si.RTSP = fmt.Sprintf("rtsp://%s:%d/%s/%s", ms.StreamIP, ms.RTSPPort, app, stream)
	}

	return si, nil
}

func (s *MediaService) IsStreamReady(ms *model.MediaServer, app, stream string) bool {
	list, err := s.zlmClient.GetMediaList("", app, stream, "", 0)
	if err != nil {
		return false
	}
	return len(list) > 0
}

func (s *MediaService) CloseStream(ms *model.MediaServer, app, stream, schema string) error {
	_, err := s.zlmClient.CloseStream("", app, stream, schema)
	return err
}

func (s *MediaService) GetSnap(ms *model.MediaServer, app, stream, path, fileName string) error {
	_, err := s.zlmClient.GetSnap("", app, stream, "15", "1")
	return err
}

func (s *MediaService) PauseRtpCheck(ms *model.MediaServer, streamKey string) (bool, error) {
	return s.zlmClient.Pause(streamKey)
}

func (s *MediaService) ResumeRtpCheck(ms *model.MediaServer, streamKey string) (bool, error) {
	return s.zlmClient.Resume(streamKey)
}

func (s *MediaService) SeekRtp(ms *model.MediaServer, streamKey string, stampSec string) (bool, error) {
	return s.zlmClient.Seek(streamKey, stampSec)
}

func (s *MediaService) SetSpeed(ms *model.MediaServer, streamKey string, speed string) (bool, error) {
	return s.zlmClient.Speed(streamKey, speed)
}

func (s *MediaService) OpenRtpServer(ms *model.MediaServer, port int, streamID string, tcpMode int, isUDP bool) (*zlm.RtpServerResult, error) {
	return s.zlmClient.OpenRtpServer(fmt.Sprintf("%d", port), streamID, tcpMode, isUDP)
}

func (s *MediaService) CloseRtpServer(ms *model.MediaServer, streamID string) (bool, error) {
	return s.zlmClient.CloseRtpServer(streamID)
}

func (s *MediaService) GetLoad(ms *model.MediaServer) (*MediaServerLoad, error) {
	return &MediaServerLoad{
		ID:            ms.ID,
		IP:            ms.IP,
		PushStreamCount: 0,
	}, nil
}

type MediaServerLoad struct {
	ID            string `json:"id"`
	IP            string `json:"ip"`
	PushStreamCount int   `json:"pushStreamCount"`
	ProxyCount    int    `json:"proxyCount"`
	GBReceiveCount int   `json:"gbReceiveCount"`
	GBSendCount   int    `json:"gbSendCount"`
}
