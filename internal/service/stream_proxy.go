package service

import (
	"fmt"
	"time"

	"wvp-pro-go/internal/database"
	"wvp-pro-go/internal/model"
	"wvp-pro-go/internal/utils"
	"wvp-pro-go/internal/zlm"

	"go.uber.org/zap"
)

type StreamProxyService struct {
	zlmClient *zlm.Client
	mediaSvc  *MediaService
	log       *zap.Logger
}

func NewStreamProxyService(zlmClient *zlm.Client, mediaSvc *MediaService, log *zap.Logger) *StreamProxyService {
	return &StreamProxyService{
		zlmClient: zlmClient,
		mediaSvc:  mediaSvc,
		log:       log,
	}
}

func (s *StreamProxyService) List(page, count int) (*utils.PageInfo[any], error) {
	var proxies []model.StreamProxy
	var total int64

	db := database.DB.Model(&model.StreamProxy{})
	db.Count(&total)
	if err := db.Offset((page - 1) * count).Limit(count).Order("id DESC").Find(&proxies).Error; err != nil {
		return nil, err
	}

	list := make([]any, len(proxies))
	for i := range proxies {
		list[i] = proxies[i]
	}
	return utils.NewPageInfo[any](total, list, page, count), nil
}

func (s *StreamProxyService) GetOne(id uint) (*model.StreamProxy, error) {
	var p model.StreamProxy
	err := database.DB.Where("id = ?", id).First(&p).Error
	return &p, err
}

func (s *StreamProxyService) Add(p *model.StreamProxy) (*StreamInfo, error) {
	// Set default values
	if p.App == "" {
		p.App = "live"
	}
	if p.Stream == "" {
		p.Stream = fmt.Sprintf("proxy_%d", time.Now().UnixMilli())
	}
	p.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	p.UpdateTime = p.CreateTime

	// Save to database
	if err := database.DB.Create(p).Error; err != nil {
		return nil, err
	}

	// If enabled, start pulling immediately
	if p.Enable {
		si, err := s.startPull(p)
		if err != nil {
			s.log.Warn("failed to start stream proxy after add", zap.Error(err), zap.Uint("id", p.ID))
			return nil, err
		}
		return si, nil
	}

	return nil, nil
}

func (s *StreamProxyService) Update(p *model.StreamProxy) error {
	p.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	return database.DB.Save(p).Error
}

func (s *StreamProxyService) Delete(id uint) error {
	proxy, err := s.GetOne(id)
	if err != nil {
		return err
	}

	// Stop pulling if active
	if proxy.Pulling {
		s.stopPull(proxy)
	}

	return database.DB.Delete(&model.StreamProxy{}, id).Error
}

func (s *StreamProxyService) Start(id uint) (*StreamInfo, error) {
	proxy, err := s.GetOne(id)
	if err != nil {
		return nil, err
	}

	if proxy.URL == "" {
		return nil, fmt.Errorf("stream proxy URL is empty")
	}

	// If already pulling, stop first to avoid "stream already exists"
	if proxy.Pulling {
		s.stopPull(proxy)
	}

	return s.startPull(proxy)
}

func (s *StreamProxyService) Stop(id uint) error {
	proxy, err := s.GetOne(id)
	if err != nil {
		return err
	}

	s.stopPull(proxy)
	return nil
}

// startPull calls ZLMediaKit addStreamProxy API to start pulling the stream
func (s *StreamProxyService) startPull(proxy *model.StreamProxy) (*StreamInfo, error) {
	timeout := proxy.Timeout
	if timeout <= 0 {
		timeout = 30
	}

	result, err := s.zlmClient.AddStreamProxy("__defaultVhost__", proxy.App, proxy.Stream, proxy.URL, proxy.EnableAudio, proxy.EnableMP4, timeout)
	if err != nil {
		return nil, fmt.Errorf("ZLM addStreamProxy failed: %v", err)
	}

	// If stream already exists, delete it first and retry
	if result.Code != 0 {
		s.log.Warn("addStreamProxy returned error, attempting to delete and retry",
			zap.Int("code", result.Code), zap.String("msg", result.Msg))
		key := fmt.Sprintf("__defaultVhost__/%s/%s", proxy.App, proxy.Stream)
		s.zlmClient.DelStreamProxy(key, "__defaultVhost__", proxy.App, proxy.Stream)
		// Retry
		result, err = s.zlmClient.AddStreamProxy("__defaultVhost__", proxy.App, proxy.Stream, proxy.URL, proxy.EnableAudio, proxy.EnableMP4, timeout)
		if err != nil {
			return nil, fmt.Errorf("ZLM addStreamProxy retry failed: %v", err)
		}
		if result.Code != 0 {
			return nil, fmt.Errorf("ZLM addStreamProxy error after retry: code=%d, msg=%s", result.Code, result.Msg)
		}
	}

	s.log.Info("stream proxy started",
		zap.Uint("id", proxy.ID),
		zap.String("app", proxy.App),
		zap.String("stream", proxy.Stream),
		zap.String("url", proxy.URL),
	)

	// Update pulling status
	proxy.Pulling = true
	proxy.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	database.DB.Save(proxy)

	// Build stream info
	ms, err := s.mediaSvc.GetMediaServerForMinimumLoad(nil)
	if err != nil {
		// Return basic info even if we can't get media server details
		return &StreamInfo{
			Stream: proxy.Stream,
			App:    proxy.App,
		}, nil
	}

	si := &StreamInfo{
		Stream:        proxy.Stream,
		App:           proxy.App,
		MediaServerID: ms.ID,
	}
	flvPort := ms.FLVPort
	if flvPort == 0 {
		flvPort = ms.HTTPPort
	}
	si.FLV = fmt.Sprintf("http://%s:%d/%s/%s.live.flv", ms.StreamIP, flvPort, proxy.App, proxy.Stream)
	si.WSFLV = fmt.Sprintf("ws://%s:%d/%s/%s.live.flv", ms.StreamIP, flvPort, proxy.App, proxy.Stream)
	si.HLS = fmt.Sprintf("http://%s:%d/%s/%s/hls.m3u8", ms.StreamIP, flvPort, proxy.App, proxy.Stream)
	si.RTMP = fmt.Sprintf("rtmp://%s:%d/%s/%s", ms.StreamIP, ms.RTMPPort, proxy.App, proxy.Stream)
	si.RTSP = fmt.Sprintf("rtsp://%s:%d/%s/%s", ms.StreamIP, ms.RTSPPort, proxy.App, proxy.Stream)

	return si, nil
}

// stopPull calls ZLMediaKit delStreamProxy API to stop pulling the stream
func (s *StreamProxyService) stopPull(proxy *model.StreamProxy) {
	key := fmt.Sprintf("__defaultVhost__/%s/%s", proxy.App, proxy.Stream)
	_, err := s.zlmClient.DelStreamProxy(key, "__defaultVhost__", proxy.App, proxy.Stream)
	if err != nil {
		s.log.Warn("failed to stop stream proxy", zap.Error(err), zap.Uint("id", proxy.ID))
	}

	// Update status
	proxy.Pulling = false
	proxy.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	database.DB.Save(proxy)

	s.log.Info("stream proxy stopped",
		zap.Uint("id", proxy.ID),
		zap.String("app", proxy.App),
		zap.String("stream", proxy.Stream),
	)
}
