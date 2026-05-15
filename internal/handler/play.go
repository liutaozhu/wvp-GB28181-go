package handler

import (
	"wvp-pro-go/internal/service"
	"wvp-pro-go/internal/utils"

	"github.com/gin-gonic/gin"
)

type PlayHandler struct {
	svcs *service.Services
}

func NewPlayHandler(svcs *service.Services) *PlayHandler {
	return &PlayHandler{svcs: svcs}
}

func (h *PlayHandler) Register(r *gin.RouterGroup) {
	r.GET("/start/:deviceId/:channelId", h.Start)
	r.GET("/stop/:deviceId/:channelId", h.Stop)
	r.POST("/convertStop/:key", h.ConvertStop)
	r.GET("/broadcast/:deviceId/:channelId", h.Broadcast)
	r.POST("/broadcast/:deviceId/:channelId", h.Broadcast)
	r.GET("/broadcast/stop/:deviceId/:channelId", h.BroadcastStop)
	r.POST("/broadcast/stop/:deviceId/:channelId", h.BroadcastStop)
	r.GET("/ssrc", h.GetSSRC)
	r.GET("/snap", h.Snap)
}

func (h *PlayHandler) Start(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	si, err := h.svcs.Play.StartPlay(deviceID, channelID)
	if err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.Success(si))
}

func (h *PlayHandler) Stop(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	if err := h.svcs.Play.StopPlay(deviceID, channelID); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PlayHandler) ConvertStop(c *gin.Context) {
	key := c.Param("key")
	if err := h.svcs.Play.StopPlay("", key); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PlayHandler) Broadcast(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	si, err := h.svcs.Play.AudioBroadcast(deviceID, channelID)
	if err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.Success(si))
}

func (h *PlayHandler) BroadcastStop(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	if err := h.svcs.Play.StopAudioBroadcast(deviceID, channelID); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PlayHandler) GetSSRC(c *gin.Context) {
	c.JSON(200, utils.Success(map[string]string{}))
}

func (h *PlayHandler) Snap(c *gin.Context) {
	deviceID := c.Query("deviceId")
	channelID := c.Query("channelId")
	if deviceID == "" || channelID == "" {
		c.JSON(200, utils.Fail(400, "缺少参数"))
		return
	}
	path, err := h.svcs.Play.GetSnap(deviceID, channelID)
	if err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.Success(path))
}
