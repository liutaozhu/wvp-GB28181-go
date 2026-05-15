package handler

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"wvp-pro-go/internal/database"
	"wvp-pro-go/internal/model"
	"wvp-pro-go/internal/service"
	"wvp-pro-go/internal/utils"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

// ServerHandler handles /api/server/* endpoints
type ServerHandler struct {
	svcs *service.Services
}

func NewServerHandler(svcs *service.Services) *ServerHandler {
	return &ServerHandler{svcs: svcs}
}

func (h *ServerHandler) Register(r *gin.RouterGroup) {
	r.GET("/media_server/online/list", h.GetOnlineMediaServerList)
	r.GET("/media_server/list", h.GetMediaServerList)
	r.GET("/media_server/one/:id", h.GetMediaServer)
	r.GET("/media_server/check", h.CheckMediaServer)
	r.GET("/media_server/record/check", h.CheckMediaServerRecord)
	r.POST("/media_server/save", h.SaveMediaServer)
	r.DELETE("/media_server/delete", h.DeleteMediaServer)
	r.GET("/media_server/media_info", h.GetMediaInfo)
	r.GET("/media_server/load", h.GetMediaServerLoad)
	r.GET("/info", h.GetInfo)
	r.GET("/system/configInfo", h.GetSystemConfig)
	r.GET("/system/info", h.GetSystemInfo)
	r.GET("/resource/info", h.GetResourceInfo)
	r.GET("/map/config", h.GetMapConfig)
	r.GET("/map/model-icon/list", h.GetModelList)
}

func (h *ServerHandler) GetOnlineMediaServerList(c *gin.Context) {
	var servers []model.MediaServer
	database.DB.Where("status = ?", true).Find(&servers)
	c.JSON(200, utils.Success(servers))
}

func (h *ServerHandler) GetMediaServerList(c *gin.Context) {
	var servers []model.MediaServer
	database.DB.Find(&servers)
	c.JSON(200, utils.Success(servers))
}

func (h *ServerHandler) GetMediaServer(c *gin.Context) {
	id := c.Param("id")
	var ms model.MediaServer
	if err := database.DB.Where("id = ?", id).First(&ms).Error; err != nil {
		c.JSON(200, utils.Fail(404, "not found"))
		return
	}
	c.JSON(200, utils.Success(ms))
}

func (h *ServerHandler) CheckMediaServer(c *gin.Context) {
	c.JSON(200, utils.Success(true))
}

func (h *ServerHandler) CheckMediaServerRecord(c *gin.Context) {
	c.JSON(200, utils.Success(true))
}

func (h *ServerHandler) SaveMediaServer(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *ServerHandler) DeleteMediaServer(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *ServerHandler) GetMediaInfo(c *gin.Context) {
	c.JSON(200, utils.Success(map[string]interface{}{
		"online": true,
	}))
}

func (h *ServerHandler) GetMediaServerLoad(c *gin.Context) {
	c.JSON(200, utils.Success([]map[string]interface{}{}))
}

func (h *ServerHandler) GetInfo(c *gin.Context) {
	hostname, _ := os.Hostname()
	uptime := time.Since(startTime)

	info := map[string]map[string]string{
		"系统信息": {
			"操作系统":  runtime.GOOS + "/" + runtime.GOARCH,
			"主机名":   hostname,
			"Go版本":  runtime.Version(),
			"CPU核数": fmt.Sprintf("%d", runtime.NumCPU()),
			"运行时间":  formatDuration(uptime),
		},
		"服务信息": {
			"版本":    "2.7.4",
			"启动时间":  startTime.Format("2006-01-02 15:04:05"),
			"Goroutines": fmt.Sprintf("%d", runtime.NumGoroutine()),
		},
		"数据统计": {
			"在线设备数": fmt.Sprintf("%d", countDevices()),
			"通道总数":  fmt.Sprintf("%d", countChannels()),
			"拉流代理数": fmt.Sprintf("%d", countProxies()),
		},
	}
	c.JSON(200, utils.Success(info))
}

func (h *ServerHandler) GetSystemConfig(c *gin.Context) {
	c.JSON(200, utils.Success(map[string]interface{}{
		"serverID": "wvp-pro-go",
		"version":  "2.7.4",
	}))
}

func (h *ServerHandler) GetSystemInfo(c *gin.Context) {
	hostname, _ := os.Hostname()
	c.JSON(200, utils.Success(map[string]interface{}{
		"os":        runtime.GOOS,
		"arch":      runtime.GOARCH,
		"hostname":  hostname,
		"goVersion": runtime.Version(),
		"cpuNum":    runtime.NumCPU(),
		"uptime":    time.Since(startTime).Seconds(),
	}))
}

func (h *ServerHandler) GetResourceInfo(c *gin.Context) {
	c.JSON(200, utils.Success(map[string]interface{}{
		"device":      countDevices(),
		"channel":     countChannels(),
		"push":        0,
		"proxy":       countProxies(),
		"gbSend":      0,
		"gbReceive":   0,
	}))
}

func (h *ServerHandler) GetMapConfig(c *gin.Context) {
	c.JSON(200, utils.Success(map[string]interface{}{
		"center": map[string]interface{}{
			"lng": 116.4,
			"lat": 39.9,
		},
		"zoom": 10,
	}))
}

func (h *ServerHandler) GetModelList(c *gin.Context) {
	c.JSON(200, utils.Success([]interface{}{}))
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	if days > 0 {
		return fmt.Sprintf("%d天%d小时%d分钟", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%d小时%d分钟", hours, minutes)
	}
	return fmt.Sprintf("%d分钟", minutes)
}

func countDevices() int64 {
	var count int64
	database.DB.Model(&model.Device{}).Count(&count)
	return count
}

func countChannels() int64 {
	var count int64
	database.DB.Model(&model.DeviceChannel{}).Count(&count)
	return count
}

func countProxies() int64 {
	var count int64
	database.DB.Model(&model.StreamProxy{}).Count(&count)
	return count
}
