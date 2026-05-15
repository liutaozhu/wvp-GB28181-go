package handler

import (
	"fmt"
	"strconv"

	"wvp-pro-go/internal/database"
	"wvp-pro-go/internal/model"
	"wvp-pro-go/internal/service"
	"wvp-pro-go/internal/utils"

	"github.com/gin-gonic/gin"
)

// StreamProxyHandler handles /api/proxy/* endpoints
type StreamProxyHandler struct {
	svcs *service.Services
}

func NewStreamProxyHandler(svcs *service.Services) *StreamProxyHandler {
	return &StreamProxyHandler{svcs: svcs}
}

func (h *StreamProxyHandler) Register(r *gin.RouterGroup) {
	r.GET("/list", h.GetList)
	r.GET("/one", h.GetOne)
	r.POST("/add", h.Add)
	r.POST("/update", h.Update)
	r.GET("/ffmpeg_cmd/list", h.GetFFmpegCmdList)
	r.DELETE("/del", h.Del)
	r.DELETE("/delete", h.Delete)
	r.GET("/start", h.Start)
	r.GET("/stop", h.Stop)
}

func (h *StreamProxyHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	count, _ := strconv.Atoi(c.DefaultQuery("count", "10"))
	result, err := h.svcs.StreamProxy.List(page, count)
	if err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.Success(result))
}

func (h *StreamProxyHandler) GetOne(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	p, err := h.svcs.StreamProxy.GetOne(uint(id))
	if err != nil {
		c.JSON(200, utils.Fail(404, "不存在"))
		return
	}
	c.JSON(200, utils.Success(p))
}

func (h *StreamProxyHandler) Add(c *gin.Context) {
	var raw map[string]interface{}
	if err := c.ShouldBindJSON(&raw); err != nil {
		c.JSON(200, utils.Fail(400, "参数错误: "+err.Error()))
		return
	}

	p := model.StreamProxy{
		Type:          getIntFromMap(raw, "type"),
		App:           getStringFromMap(raw, "app"),
		Stream:        getStringFromMap(raw, "stream"),
		Name:          getStringFromMap(raw, "name"),
		Description:   getStringFromMap(raw, "description"),
		MediaServerID: getStringFromMap(raw, "mediaServerId"),
		FFmpegCMD:     getStringFromMap(raw, "ffmpegCmd"),
		RemoveKey:     getStringFromMap(raw, "removeKey"),
		DeviceID:      getStringFromMap(raw, "deviceId"),
		Timeout:       getIntFromMap(raw, "timeout"),
		EnableAudio:   getBoolFromMap(raw, "enableAudio", "enable_audio"),
		EnableMP4:     getBoolFromMap(raw, "enableMp4", "enable_mp4"),
		Enable:        getBoolFromMap(raw, "enable"),
		EnableRemoveKey: getBoolFromMap(raw, "enableRemoveKey", "enable_remove_key"),
	}

	// Support both "url" and "src_url" field names
	p.URL = getStringFromMap(raw, "url")
	if p.URL == "" {
		p.URL = getStringFromMap(raw, "src_url")
	}
	if p.URL == "" {
		p.URL = getStringFromMap(raw, "srcUrl")
	}

	si, err := h.svcs.StreamProxy.Add(&p)
	if err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	if si != nil {
		c.JSON(200, utils.Success(si))
	} else {
		c.JSON(200, utils.Success(p))
	}
}

func getStringFromMap(m map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if v, ok := m[key]; ok {
			if s, ok := v.(string); ok {
				return s
			}
			if f, ok := v.(float64); ok {
				return strconv.Itoa(int(f))
			}
		}
	}
	return ""
}

func getIntFromMap(m map[string]interface{}, keys ...string) int {
	for _, key := range keys {
		if v, ok := m[key]; ok {
			if f, ok := v.(float64); ok {
				return int(f)
			}
			if s, ok := v.(string); ok {
				n, _ := strconv.Atoi(s)
				return n
			}
		}
	}
	return 0
}

func getBoolFromMap(m map[string]interface{}, keys ...string) bool {
	for _, key := range keys {
		if v, ok := m[key]; ok {
			switch val := v.(type) {
			case bool:
				return val
			case float64:
				return val != 0
			case string:
				return val == "true" || val == "1"
			}
		}
	}
	return false
}

func (h *StreamProxyHandler) Update(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *StreamProxyHandler) GetFFmpegCmdList(c *gin.Context) {
	c.JSON(200, utils.Success(map[string]string{}))
}

func (h *StreamProxyHandler) Del(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if err := h.svcs.StreamProxy.Delete(uint(id)); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *StreamProxyHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if err := h.svcs.StreamProxy.Delete(uint(id)); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *StreamProxyHandler) Start(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	si, err := h.svcs.StreamProxy.Start(uint(id))
	if err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.Success(si))
}

func (h *StreamProxyHandler) Stop(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if err := h.svcs.StreamProxy.Stop(uint(id)); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

// PlatformHandler handles /api/platform/* endpoints
type PlatformHandler struct {
	svcs *service.Services
}

func NewPlatformHandler(svcs *service.Services) *PlatformHandler {
	return &PlatformHandler{svcs: svcs}
}

func (h *PlatformHandler) Register(r *gin.RouterGroup) {
	r.GET("/server_config", h.GetServerConfig)
	r.GET("/info/:id", h.GetInfo)
	r.GET("/query", h.Query)
	r.POST("/add", h.Add)
	r.POST("/update", h.Update)
	r.DELETE("/delete", h.Delete)
	r.GET("/exit/:serverGBId", h.Exit)
	r.GET("/channel/list", h.GetChannelList)
	r.POST("/channel/add", h.AddChannel)
	r.DELETE("/channel/remove", h.RemoveChannel)
	r.GET("/channel/push", h.PushChannel)
	r.POST("/channel/device/add", h.AddChannelByDevice)
	r.POST("/channel/device/remove", h.RemoveChannelByDevice)
	r.POST("/channel/custom/update", h.UpdateCustomChannel)
}

func (h *PlatformHandler) GetServerConfig(c *gin.Context) {
	c.JSON(200, utils.Success(map[string]interface{}{
		"serverId": "wvp-go-server",
	}))
}

func (h *PlatformHandler) GetInfo(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svcs.Platform.GetPlatform(uint(id))
	if err != nil {
		c.JSON(200, utils.Fail(404, "不存在"))
		return
	}
	c.JSON(200, utils.Success(p))
}

func (h *PlatformHandler) Query(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	count, _ := strconv.Atoi(c.DefaultQuery("count", "10"))
	result, err := h.svcs.Platform.ListPlatforms(page, count)
	if err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.Success(result))
}

func (h *PlatformHandler) Add(c *gin.Context) {
	c.JSON(200, utils.Fail(400, "未实现"))
}

func (h *PlatformHandler) Update(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *PlatformHandler) Delete(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *PlatformHandler) Exit(c *gin.Context) {
	c.JSON(200, utils.Success(true))
}

func (h *PlatformHandler) GetChannelList(c *gin.Context) {
	platformID, _ := strconv.ParseUint(c.Query("platformId"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	count, _ := strconv.Atoi(c.DefaultQuery("count", "10"))
	result, err := h.svcs.Platform.GetPlatformChannels(uint(platformID), page, count)
	if err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.Success(result))
}

func (h *PlatformHandler) AddChannel(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *PlatformHandler) RemoveChannel(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *PlatformHandler) PushChannel(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *PlatformHandler) AddChannelByDevice(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *PlatformHandler) RemoveChannelByDevice(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *PlatformHandler) UpdateCustomChannel(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

// RegionHandler handles /api/region/* endpoints
type RegionHandler struct {
	svcs *service.Services
}

func NewRegionHandler(svcs *service.Services) *RegionHandler {
	return &RegionHandler{svcs: svcs}
}

func (h *RegionHandler) Register(r *gin.RouterGroup) {
	r.POST("/add", h.Add)
	r.GET("/page/list", h.GetPageList)
	r.GET("/tree/list", h.GetTreeList)
	r.GET("/tree/query", h.GetTreeQuery)
	r.POST("/update", h.Update)
	r.DELETE("/delete", h.Delete)
	r.GET("/one", h.GetOne)
	r.GET("/base/child/list", h.GetChildList)
	r.GET("/path", h.GetPath)
	r.GET("/sync", h.Sync)
	r.GET("/description", h.GetDescription)
	r.GET("/addByCivilCode", h.AddByCivilCode)
}

func (h *RegionHandler) Add(c *gin.Context) {
	c.JSON(200, utils.Fail(400, "未实现"))
}

func (h *RegionHandler) GetPageList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	count, _ := strconv.Atoi(c.DefaultQuery("count", "10"))
	result, err := h.svcs.Region.GetPageList(page, count)
	if err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.Success(result))
}

func (h *RegionHandler) GetTreeList(c *gin.Context) {
	parent := c.Query("parent")
	hasChannel := c.Query("hasChannel") == "true"

	var result []map[string]interface{}

	if parent == "" && hasChannel {
		// Root level: return all channels as leaf nodes directly
		var channels []model.DeviceChannel
		database.DB.Order("gb_id ASC").Limit(200).Find(&channels)

		// Also get stream proxies as virtual channels
		var proxies []model.StreamProxy
		database.DB.Order("id ASC").Find(&proxies)

		for _, ch := range channels {
			status := ch.Status
			if status == "" {
				status = ch.GBStatus
			}
			name := ch.GBName
			if name == "" {
				name = ch.Name
			}
			result = append(result, map[string]interface{}{
				"treeId":   fmt.Sprintf("channel_%d", ch.ID),
				"deviceId": ch.GBDeviceID,
				"name":     name,
				"type":     1,
				"isLeaf":   true,
				"status":   status,
			})
		}

		for _, p := range proxies {
			status := "OFF"
			if p.Pulling {
				status = "ON"
			}
			name := p.Name
			if name == "" {
				name = p.App + "/" + p.Stream
			}
			result = append(result, map[string]interface{}{
				"treeId":   fmt.Sprintf("proxy_%d", p.ID),
				"deviceId": fmt.Sprintf("proxy_%d", p.ID),
				"name":     name,
				"type":     1,
				"isLeaf":   true,
				"status":   status,
			})
		}
	} else {
		// Return region tree nodes
		tree, err := h.svcs.Region.GetTreeList()
		if err != nil {
			c.JSON(200, utils.Fail(500, err.Error()))
			return
		}
		// Convert to frontend format
		for _, node := range tree {
			result = append(result, map[string]interface{}{
				"treeId":   fmt.Sprintf("region_%v", node["id"]),
				"deviceId": fmt.Sprintf("%v", node["id"]),
				"name":     node["name"],
				"type":     0,
				"isLeaf":   node["children"] == nil,
			})
		}
	}

	if result == nil {
		result = []map[string]interface{}{}
	}
	c.JSON(200, utils.Success(result))
}

func (h *RegionHandler) GetTreeQuery(c *gin.Context) {
	c.JSON(200, utils.Fail(400, "未实现"))
}

func (h *RegionHandler) Update(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *RegionHandler) Delete(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *RegionHandler) GetOne(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	r, err := h.svcs.Region.GetOne(uint(id))
	if err != nil {
		c.JSON(200, utils.Fail(404, "不存在"))
		return
	}
	c.JSON(200, utils.Success(r))
}

func (h *RegionHandler) GetChildList(c *gin.Context) {
	parentID, _ := strconv.ParseUint(c.Query("parentId"), 10, 64)
	children, err := h.svcs.Region.GetChildren(uint(parentID))
	if err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.Success(children))
}

func (h *RegionHandler) GetPath(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	path, err := h.svcs.Region.GetPath(uint(id))
	if err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.Success(path))
}

func (h *RegionHandler) Sync(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *RegionHandler) GetDescription(c *gin.Context) {
	c.JSON(200, utils.Success(""))
}

func (h *RegionHandler) AddByCivilCode(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

// GroupHandler handles /api/group/* endpoints
type GroupHandler struct {
	svcs *service.Services
}

func NewGroupHandler(svcs *service.Services) *GroupHandler {
	return &GroupHandler{svcs: svcs}
}

func (h *GroupHandler) Register(r *gin.RouterGroup) {
	r.POST("/add", h.Add)
	r.GET("/tree/list", h.GetTreeList)
	r.GET("/tree/query", h.GetTreeQuery)
	r.POST("/update", h.Update)
	r.DELETE("/delete", h.Delete)
	r.GET("/path", h.GetPath)
}

func (h *GroupHandler) Add(c *gin.Context) {
	c.JSON(200, utils.Fail(400, "未实现"))
}

func (h *GroupHandler) GetTreeList(c *gin.Context) {
	hasChannel := c.Query("hasChannel") == "true"

	var result []map[string]interface{}

	if hasChannel {
		// Return all channels as leaf nodes
		var channels []model.DeviceChannel
		database.DB.Order("gb_id ASC").Limit(200).Find(&channels)

		var proxies []model.StreamProxy
		database.DB.Order("id ASC").Find(&proxies)

		for _, ch := range channels {
			status := ch.Status
			if status == "" {
				status = ch.GBStatus
			}
			name := ch.GBName
			if name == "" {
				name = ch.Name
			}
			result = append(result, map[string]interface{}{
				"treeId":   fmt.Sprintf("channel_%d", ch.ID),
				"deviceId": ch.GBDeviceID,
				"name":     name,
				"type":     1,
				"isLeaf":   true,
				"status":   status,
			})
		}

		for _, p := range proxies {
			status := "OFF"
			if p.Pulling {
				status = "ON"
			}
			name := p.Name
			if name == "" {
				name = p.App + "/" + p.Stream
			}
			result = append(result, map[string]interface{}{
				"treeId":   fmt.Sprintf("proxy_%d", p.ID),
				"deviceId": fmt.Sprintf("proxy_%d", p.ID),
				"name":     name,
				"type":     1,
				"isLeaf":   true,
				"status":   status,
			})
		}
	} else {
		tree, err := h.svcs.Group.GetTreeList()
		if err != nil {
			c.JSON(200, utils.Fail(500, err.Error()))
			return
		}
		for _, node := range tree {
			result = append(result, map[string]interface{}{
				"treeId":   fmt.Sprintf("group_%v", node["id"]),
				"deviceId": fmt.Sprintf("%v", node["id"]),
				"name":     node["name"],
				"type":     0,
				"isLeaf":   node["children"] == nil,
			})
		}
	}

	if result == nil {
		result = []map[string]interface{}{}
	}
	c.JSON(200, utils.Success(result))
}

func (h *GroupHandler) GetTreeQuery(c *gin.Context) {
	c.JSON(200, utils.Fail(400, "未实现"))
}

func (h *GroupHandler) Update(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *GroupHandler) Delete(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *GroupHandler) GetPath(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	path, err := h.svcs.Group.GetPath(uint(id))
	if err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.Success(path))
}
