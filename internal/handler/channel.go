package handler

import (
	"strconv"

	"wvp-pro-go/internal/service"
	"wvp-pro-go/internal/utils"

	"github.com/gin-gonic/gin"
)

type ChannelHandler struct {
	svcs *service.Services
}

func NewChannelHandler(svcs *service.Services) *ChannelHandler {
	return &ChannelHandler{svcs: svcs}
}

func (h *ChannelHandler) Register(r *gin.RouterGroup) {
	r.GET("/one", h.GetOne)
	r.GET("/industry/list", h.GetIndustryList)
	r.GET("/type/list", h.GetTypeList)
	r.GET("/network/identification/list", h.GetNetworkIdentificationList)
	r.POST("/update", h.Update)
	r.POST("/reset", h.Reset)
	r.POST("/add", h.Add)
	r.GET("/list", h.GetList)
	r.GET("/civilcode/list", h.GetCivilCodeList)
	r.GET("/civilCode/unusual/list", h.GetCivilCodeUnusualList)
	r.GET("/parent/unusual/list", h.GetParentUnusualList)
	r.POST("/civilCode/unusual/clear", h.ClearCivilCodeUnusual)
	r.POST("/parent/unusual/clear", h.ClearParentUnusual)
	r.GET("/parent/list", h.GetParentList)
	r.POST("/region/add", h.AddRegion)
	r.POST("/region/delete", h.DeleteRegion)
	r.POST("/region/device/add", h.AddRegionByDevice)
	r.POST("/region/device/delete", h.DeleteRegionByDevice)
	r.POST("/group/add", h.AddGroup)
	r.POST("/group/delete", h.DeleteGroup)
	r.POST("/group/device/add", h.AddGroupByDevice)
	r.POST("/group/device/delete", h.DeleteGroupByDevice)
	r.GET("/play", h.Play)
	r.GET("/play/stop", h.PlayStop)
	r.GET("/playback/query", h.PlaybackQuery)
	r.GET("/playback", h.Playback)
	r.GET("/playback/stop", h.PlaybackStop)
	r.GET("/playback/pause", h.PlaybackPause)
	r.GET("/playback/resume", h.PlaybackResume)
	r.GET("/playback/seek", h.PlaybackSeek)
	r.GET("/playback/speed", h.PlaybackSpeed)
	r.GET("/map/list", h.GetMapList)
	r.POST("/map/reset-level", h.ResetMapLevel)
	r.POST("/map/thin/draw", h.DrawThin)
	r.GET("/map/thin/clear", h.ClearThin)
	r.GET("/map/thin/save", h.SaveThin)
	r.GET("/map/thin/progress", h.GetThinProgress)
	r.GET("/map/tile/:z/:x/:y", h.GetTile)
	r.GET("/map/thin/tile/:z/:x/:y", h.GetThinTile)
}

func (h *ChannelHandler) GetOne(c *gin.Context) {
	deviceID := c.Query("deviceId")
	channelID := c.Query("channelId")
	if deviceID == "" || channelID == "" {
		c.JSON(200, utils.Fail(400, "缺少参数"))
		return
	}
	ch, err := h.svcs.Channel.GetChannel(deviceID, channelID)
	if err != nil {
		c.JSON(200, utils.Fail(404, "通道不存在"))
		return
	}
	c.JSON(200, utils.Success(ch))
}

func (h *ChannelHandler) GetIndustryList(c *gin.Context) {
	c.JSON(200, utils.Success([]string{}))
}

func (h *ChannelHandler) GetTypeList(c *gin.Context) {
	c.JSON(200, utils.Success([]string{}))
}

func (h *ChannelHandler) GetNetworkIdentificationList(c *gin.Context) {
	c.JSON(200, utils.Success([]string{}))
}

func (h *ChannelHandler) Update(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) Reset(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) Add(c *gin.Context) {
	c.JSON(200, utils.Fail(400, "未实现"))
}

func (h *ChannelHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	count, _ := strconv.Atoi(c.DefaultQuery("count", "10"))
	deviceID := c.Query("deviceId")
	query := c.Query("query")

	var online *bool
	if c.Query("online") != "" {
		v := c.Query("online") == "true"
		online = &v
	}

	var hasSub *bool
	if c.Query("hasSubChannel") != "" {
		v := c.Query("hasSubChannel") == "true"
		hasSub = &v
	}

	var hasStream *bool
	if c.Query("hasStream") != "" {
		v := c.Query("hasStream") == "true"
		hasStream = &v
	}

	if deviceID != "" {
		result, err := h.svcs.Channel.GetChannels(deviceID, page, count)
		if err != nil {
			c.JSON(200, utils.Fail(500, err.Error()))
			return
		}
		c.JSON(200, utils.Success(result))
	} else {
		result, err := h.svcs.Channel.GetAllChannels(query, online, hasSub, hasStream, page, count)
		if err != nil {
			c.JSON(200, utils.Fail(500, err.Error()))
			return
		}
		c.JSON(200, utils.Success(result))
	}
}

func (h *ChannelHandler) GetCivilCodeList(c *gin.Context) {
	list, err := h.svcs.Channel.GetCivilCodeList()
	if err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.Success(list))
}

func (h *ChannelHandler) GetCivilCodeUnusualList(c *gin.Context) {
	c.JSON(200, utils.Success([]string{}))
}

func (h *ChannelHandler) GetParentUnusualList(c *gin.Context) {
	c.JSON(200, utils.Success([]string{}))
}

func (h *ChannelHandler) ClearCivilCodeUnusual(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) ClearParentUnusual(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) GetParentList(c *gin.Context) {
	channels, err := h.svcs.Channel.GetParentChannels()
	if err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.Success(channels))
}

func (h *ChannelHandler) AddRegion(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) DeleteRegion(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) AddRegionByDevice(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) DeleteRegionByDevice(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) AddGroup(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) DeleteGroup(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) AddGroupByDevice(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) DeleteGroupByDevice(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) Play(c *gin.Context) {
	channelID := c.Query("channelId")
	deviceID := c.Query("deviceId")

	if channelID == "" {
		c.JSON(200, utils.Fail(400, "缺少参数channelId"))
		return
	}

	// Handle proxy channels (channelId = "proxy_N")
	if len(channelID) > 6 && channelID[:6] == "proxy_" {
		idStr := channelID[6:]
		id, _ := strconv.ParseUint(idStr, 10, 64)
		if id == 0 {
			c.JSON(200, utils.Fail(400, "无效的代理通道ID"))
			return
		}
		si, err := h.svcs.StreamProxy.Start(uint(id))
		if err != nil {
			c.JSON(200, utils.Fail(500, err.Error()))
			return
		}
		c.JSON(200, utils.Success(si))
		return
	}

	// Handle numeric channelId >= 100000 (proxy channel from channel list page)
	if numID, err := strconv.ParseUint(channelID, 10, 64); err == nil && numID >= 100000 {
		proxyID := uint(numID - 100000)
		si, err := h.svcs.StreamProxy.Start(proxyID)
		if err != nil {
			c.JSON(200, utils.Fail(500, err.Error()))
			return
		}
		c.JSON(200, utils.Success(si))
		return
	}

	// For real GB channels, if deviceId not provided, look it up
	if deviceID == "" {
		ch, err := h.svcs.Channel.GetChannelByGBDeviceID(channelID)
		if err != nil {
			c.JSON(200, utils.Fail(404, "通道不存在"))
			return
		}
		deviceID = ch.DeviceID
	}

	si, err := h.svcs.Play.StartPlay(deviceID, channelID)
	if err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.Success(si))
}

func (h *ChannelHandler) PlayStop(c *gin.Context) {
	deviceID := c.Query("deviceId")
	channelID := c.Query("channelId")
	if deviceID == "" || channelID == "" {
		c.JSON(200, utils.Fail(400, "缺少参数"))
		return
	}
	if err := h.svcs.Play.StopPlay(deviceID, channelID); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) PlaybackQuery(c *gin.Context) {
	deviceID := c.Query("deviceId")
	channelID := c.Query("channelId")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	if deviceID == "" || channelID == "" || startTime == "" || endTime == "" {
		c.JSON(200, utils.Fail(400, "缺少参数"))
		return
	}
	records, err := h.svcs.Playback.QueryRecord(deviceID, channelID, startTime, endTime)
	if err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.Success(records))
}

func (h *ChannelHandler) Playback(c *gin.Context) {
	deviceID := c.Query("deviceId")
	channelID := c.Query("channelId")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	if deviceID == "" || channelID == "" || startTime == "" || endTime == "" {
		c.JSON(200, utils.Fail(400, "缺少参数"))
		return
	}
	si, err := h.svcs.Playback.StartPlayback(deviceID, channelID, startTime, endTime)
	if err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.Success(si))
}

func (h *ChannelHandler) PlaybackStop(c *gin.Context) {
	deviceID := c.Query("deviceId")
	channelID := c.Query("channelId")
	stream := c.Query("stream")
	if err := h.svcs.Playback.StopPlayback(deviceID, channelID, stream); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) PlaybackPause(c *gin.Context) {
	stream := c.Query("stream")
	if err := h.svcs.Playback.PlaybackPause(stream); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) PlaybackResume(c *gin.Context) {
	stream := c.Query("stream")
	if err := h.svcs.Playback.PlaybackResume(stream); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) PlaybackSeek(c *gin.Context) {
	stream := c.Query("stream")
	seekTime, _ := strconv.ParseInt(c.Query("seekTime"), 10, 64)
	if err := h.svcs.Playback.PlaybackSeek(stream, seekTime); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) PlaybackSpeed(c *gin.Context) {
	stream := c.Query("stream")
	speed, _ := strconv.ParseFloat(c.Query("speed"), 64)
	if err := h.svcs.Playback.PlaybackSpeed(stream, speed); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) GetMapList(c *gin.Context) {
	tree, err := h.svcs.Channel.GetRegionTree()
	if err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.Success(tree))
}

func (h *ChannelHandler) ResetMapLevel(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) DrawThin(c *gin.Context) {
	c.JSON(200, utils.Success(""))
}

func (h *ChannelHandler) ClearThin(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) SaveThin(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *ChannelHandler) GetThinProgress(c *gin.Context) {
	c.JSON(200, utils.Success(0))
}

func (h *ChannelHandler) GetTile(c *gin.Context) {
	c.Data(200, "application/x-protobuf", []byte{})
}

func (h *ChannelHandler) GetThinTile(c *gin.Context) {
	c.Data(200, "application/x-protobuf", []byte{})
}
