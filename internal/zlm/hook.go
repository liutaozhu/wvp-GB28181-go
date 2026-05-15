package zlm

import (
	"net/http"
	"time"

	"wvp-pro-go/internal/config"
	"wvp-pro-go/internal/database"
	"wvp-pro-go/internal/event"
	"wvp-pro-go/internal/model"
	"wvp-pro-go/internal/sip"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HookHandler handles ZLMediaKit HTTP hook callbacks
type HookHandler struct {
	log          *zap.Logger
	cfg          config.MediaConfig
	eventBus     *event.Bus
	sessionMgr   *sip.SessionManager
	ssrcMgr      *sip.SSRCManager
	sipCmd       *sip.Commander
}

func NewHookHandler(log *zap.Logger, cfg config.MediaConfig, eventBus *event.Bus,
	sessionMgr *sip.SessionManager, ssrcMgr *sip.SSRCManager, sipCmd *sip.Commander) *HookHandler {
	return &HookHandler{
		log:        log,
		cfg:        cfg,
		eventBus:   eventBus,
		sessionMgr: sessionMgr,
		ssrcMgr:    ssrcMgr,
		sipCmd:     sipCmd,
	}
}

// Register registers all hook endpoints on the given router
func (h *HookHandler) Register(r *gin.RouterGroup) {
	r.POST("/on_publish", h.OnPublish)
	r.POST("/on_play", h.OnPlay)
	r.POST("/on_stream_changed", h.OnStreamChanged)
	r.POST("/on_stream_none_reader", h.OnStreamNoneReader)
	r.POST("/on_rtp_server_timeout", h.OnRtpServerTimeout)
	r.POST("/on_record_mp4", h.OnRecordMp4)
	r.POST("/on_flow_report", h.OnFlowReport)
	r.POST("/on_http_access", h.OnHTTPAccess)
	r.POST("/on_server_started", h.OnServerStarted)
	r.POST("/on_server_exited", h.OnServerExited)
	r.POST("/on_send_rtp_stopped", h.OnSendRtpStopped)
	r.POST("/on_record_ts", h.OnRecordTs)
}

func (h *HookHandler) parseHookParam(c *gin.Context) (map[string]interface{}, error) {
	var params map[string]interface{}
	if err := c.ShouldBindJSON(&params); err != nil {
		return nil, err
	}
	return params, nil
}

func (h *HookHandler) respond(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{"code": code, "msg": msg})
}

// OnPublish handles stream publish event
// Called when a stream is pushed to ZLMediaKit
func (h *HookHandler) OnPublish(c *gin.Context) {
	param, err := h.parseHookParam(c)
	if err != nil {
		h.respond(c, 0, "success")
		return
	}

	app := getString(param, "app")
	stream := getString(param, "stream")
	ip := getString(param, "ip")

	h.log.Info("on_publish hook",
		zap.String("app", app),
		zap.String("stream", stream),
		zap.String("ip", ip),
	)

	// Check if this stream corresponds to a known GB28181 session
	session, err := h.sessionMgr.Get(stream)
	if err == nil && session != nil {
		// This is a GB28181 RTP stream that arrived
		h.log.Info("GB28181 stream arrived",
			zap.String("deviceID", session.DeviceID),
			zap.String("channelID", session.ChannelID),
			zap.String("stream", stream),
			zap.String("ssrc", session.SSRC),
		)

		// Notify that stream has started
		h.eventBus.Publish(event.EventStreamStart, event.StreamEvent{
			Stream:        stream,
			DeviceID:      session.DeviceID,
			ChannelID:     session.ChannelID,
			App:           app,
			MediaServerID: session.MediaServerID,
			SSRC:          session.SSRC,
		})
	} else {
		// This might be a regular push stream (not GB28181)
		// Update StreamPush model if it exists
		h.handlePushStreamPublish(app, stream, ip)
	}

	h.respond(c, 0, "success")
}

// OnPlay handles stream play event
// Called when a client starts playing a stream
func (h *HookHandler) OnPlay(c *gin.Context) {
	param, err := h.parseHookParam(c)
	if err != nil {
		h.respond(c, 0, "success")
		return
	}

	app := getString(param, "app")
	stream := getString(param, "stream")
	id := getString(param, "id")
	params := getString(param, "params")

	h.log.Info("on_play hook",
		zap.String("app", app),
		zap.String("stream", stream),
		zap.String("id", id),
		zap.String("params", params),
	)

	// If this is an on-demand stream (GB28181), ensure the device stream is active
	session, err := h.sessionMgr.Get(stream)
	if err == nil && session != nil {
		h.log.Info("play request for GB28181 stream",
			zap.String("deviceID", session.DeviceID),
			zap.String("channelID", session.ChannelID),
		)
	}

	h.respond(c, 0, "success")
}

// OnStreamChanged handles stream state change event
// Called when a stream goes online or offline
func (h *HookHandler) OnStreamChanged(c *gin.Context) {
	param, err := h.parseHookParam(c)
	if err != nil {
		h.respond(c, 0, "success")
		return
	}

	online := getBool(param, "online")
	app := getString(param, "app")
	stream := getString(param, "stream")

	h.log.Info("on_stream_changed hook",
		zap.String("app", app),
		zap.String("stream", stream),
		zap.Bool("online", online),
	)

	// Check if this is a known GB28181 session
	session, err := h.sessionMgr.Get(stream)
	if err == nil && session != nil {
		if online {
			// Stream is ready
			h.log.Info("GB28181 stream is ready",
				zap.String("deviceID", session.DeviceID),
				zap.String("channelID", session.ChannelID),
				zap.String("ssrc", session.SSRC),
			)
			h.eventBus.Publish(event.EventStreamChange, event.StreamEvent{
				Stream:    stream,
				DeviceID:  session.DeviceID,
				ChannelID: session.ChannelID,
				App:       app,
				Online:    true,
				SSRC:      session.SSRC,
			})
		} else {
			// Stream went offline
			h.log.Info("GB28181 stream went offline",
				zap.String("deviceID", session.DeviceID),
				zap.String("channelID", session.ChannelID),
			)
			h.eventBus.Publish(event.EventStreamChange, event.StreamEvent{
				Stream:    stream,
				DeviceID:  session.DeviceID,
				ChannelID: session.ChannelID,
				App:       app,
				Online:    false,
				SSRC:      session.SSRC,
			})

			// Clean up session
			h.sessionMgr.Delete(stream)
			h.ssrcMgr.ReleaseSSRC(session.SSRC)
		}
	}

	h.respond(c, 0, "success")
}

// OnStreamNoneReader handles no viewer event
// Called when there are no viewers for a stream
func (h *HookHandler) OnStreamNoneReader(c *gin.Context) {
	param, err := h.parseHookParam(c)
	if err != nil {
		h.respond(c, 0, "success")
		return
	}

	app := getString(param, "app")
	stream := getString(param, "stream")

	h.log.Info("on_stream_none_reader hook",
		zap.String("app", app),
		zap.String("stream", stream),
	)

	// Check if this stream should be auto-closed when no viewers
	// For GB28181 streams, we typically want to close them to save bandwidth
	session, err := h.sessionMgr.Get(stream)
	if err == nil && session != nil {
		h.log.Info("no viewers for GB28181 stream, closing",
			zap.String("deviceID", session.DeviceID),
			zap.String("channelID", session.ChannelID),
			zap.String("stream", stream),
		)

		// Send BYE to device to stop the stream
		go func() {
			if err := h.sipCmd.StreamByeCmd(session.DeviceID, session.ChannelID, session.CallID); err != nil {
				h.log.Warn("failed to send BYE on none reader", zap.Error(err))
			}
			h.sessionMgr.Delete(stream)
			h.ssrcMgr.ReleaseSSRC(session.SSRC)

			// Close RTP server
			if session.MediaServerID != "" {
				var ms model.MediaServer
				if err := database.DB.Where("id = ?", session.MediaServerID).First(&ms).Error; err == nil {
					// Use zlm client to close RTP server
					// Note: we need access to zlm client here - in production this would be injected
				}
			}

			h.eventBus.Publish(event.EventStreamStop, event.StreamEvent{
				Stream:    stream,
				DeviceID:  session.DeviceID,
				ChannelID: session.ChannelID,
				App:       app,
			})
		}()
	}

	h.respond(c, 0, "success")
}

// OnRtpServerTimeout handles RTP server timeout event
// Called when no RTP data is received within the timeout period
func (h *HookHandler) OnRtpServerTimeout(c *gin.Context) {
	param, err := h.parseHookParam(c)
	if err != nil {
		h.respond(c, 0, "success")
		return
	}

	streamID := getString(param, "stream_id")

	h.log.Warn("on_rtp_server_timeout hook",
		zap.String("stream_id", streamID),
	)

	// Find and clean up the session
	session, err := h.sessionMgr.Get(streamID)
	if err == nil && session != nil {
		h.log.Warn("RTP server timeout, cleaning up session",
			zap.String("deviceID", session.DeviceID),
			zap.String("channelID", session.ChannelID),
			zap.String("stream", streamID),
			zap.String("ssrc", session.SSRC),
		)

		// Release SSRC
		h.ssrcMgr.ReleaseSSRC(session.SSRC)

		// Delete session
		h.sessionMgr.Delete(streamID)

		// Notify event bus
		h.eventBus.Publish(event.EventStreamStop, event.StreamEvent{
			Stream:    streamID,
			DeviceID:  session.DeviceID,
			ChannelID: session.ChannelID,
			SSRC:      session.SSRC,
		})
	}

	h.respond(c, 0, "success")
}

// OnRecordMp4 handles MP4 recording complete event
// Called when an MP4 recording file is generated
func (h *HookHandler) OnRecordMp4(c *gin.Context) {
	param, err := h.parseHookParam(c)
	if err != nil {
		h.respond(c, 0, "success")
		return
	}

	app := getString(param, "app")
	stream := getString(param, "stream")
	filePath := getString(param, "file_path")
	filePathName := getString(param, "file_path_name")
	fileName := getString(param, "file_name")
	folder := getString(param, "folder")
	url := getString(param, "url")
	timeLen := getFloat64(param, "time")
	startTime := getFloat64(param, "start_time")
	duration := getFloat64(param, "duration")

	h.log.Info("on_record_mp4 hook",
		zap.String("app", app),
		zap.String("stream", stream),
		zap.String("file_path", filePath),
		zap.String("file_name", fileName),
		zap.Float64("duration", duration),
		zap.Float64("time", timeLen),
	)

	// Save recording info to database
	rec := &model.Record{
		App:        app,
		Stream:     stream,
		FilePath:   filePathName,
		Folder:     folder,
		FileName:   fileName,
		URL:        url,
		Duration:   duration,
		StartTime:  time.Unix(int64(startTime), 0).Format("2006-01-02 15:04:05"),
		EndTime:    time.Unix(int64(startTime+duration), 0).Format("2006-01-02 15:04:05"),
		FileSize:   getInt64(param, "file_size"),
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := database.DB.Create(rec).Error; err != nil {
		h.log.Warn("failed to save record info", zap.Error(err))
	} else {
		h.log.Info("recording info saved",
			zap.String("stream", stream),
			zap.String("file", fileName),
			zap.Float64("duration", duration),
		)
	}

	h.respond(c, 0, "success")
}

// OnFlowReport handles flow report event
func (h *HookHandler) OnFlowReport(c *gin.Context) {
	_, err := h.parseHookParam(c)
	if err != nil {
		h.respond(c, 0, "success")
		return
	}
	h.respond(c, 0, "success")
}

// OnHTTPAccess handles HTTP access event
func (h *HookHandler) OnHTTPAccess(c *gin.Context) {
	_, err := h.parseHookParam(c)
	if err != nil {
		h.respond(c, 0, "success")
		return
	}
	h.respond(c, 0, "success")
}

// OnServerStarted handles server started event
func (h *HookHandler) OnServerStarted(c *gin.Context) {
	h.log.Info("ZLM server started")
	h.respond(c, 0, "success")
}

// OnServerExited handles server exited event
func (h *HookHandler) OnServerExited(c *gin.Context) {
	h.log.Warn("ZLM server exited")
	h.respond(c, 0, "success")
}

// OnSendRtpStopped handles send RTP stopped event
func (h *HookHandler) OnSendRtpStopped(c *gin.Context) {
	param, err := h.parseHookParam(c)
	if err != nil {
		h.respond(c, 0, "success")
		return
	}
	app := getString(param, "app")
	stream := getString(param, "stream")
	ssrc := getString(param, "ssrc")

	h.log.Info("on_send_rtp_stopped hook",
		zap.String("app", app),
		zap.String("stream", stream),
		zap.String("ssrc", ssrc),
	)

	// Check if this is a GB28181 cascade stream
	session, err := h.sessionMgr.Get(stream)
	if err == nil && session != nil {
		h.eventBus.Publish(event.EventStreamStop, event.StreamEvent{
			Stream:    stream,
			DeviceID:  session.DeviceID,
			ChannelID: session.ChannelID,
			SSRC:      ssrc,
		})
	}

	h.respond(c, 0, "success")
}

// OnRecordTs handles TS recording event
func (h *HookHandler) OnRecordTs(c *gin.Context) {
	_, err := h.parseHookParam(c)
	if err != nil {
		h.respond(c, 0, "success")
		return
	}
	h.respond(c, 0, "success")
}

// handlePushStreamPublish handles non-GB28181 stream publish (no-op since push stream feature removed)
func (h *HookHandler) handlePushStreamPublish(app, stream, ip string) {
	h.log.Debug("non-GB28181 stream publish (ignored)",
		zap.String("app", app),
		zap.String("stream", stream),
	)
}

// Helper functions
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
		if f, ok := v.(float64); ok {
			return f != 0
		}
	}
	return false
}

func getFloat64(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return 0
}

func getInt64(m map[string]interface{}, key string) int64 {
	if v, ok := m[key]; ok {
		if f, ok := v.(float64); ok {
			return int64(f)
		}
		if i, ok := v.(int64); ok {
			return i
		}
	}
	return 0
}
