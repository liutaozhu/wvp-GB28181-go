package handler

import (
	"strconv"

	"wvp-pro-go/internal/service"
	"wvp-pro-go/internal/utils"

	"github.com/gin-gonic/gin"
)

type PtzHandler struct {
	svcs *service.Services
}

func NewPtzHandler(svcs *service.Services) *PtzHandler {
	return &PtzHandler{svcs: svcs}
}

func (h *PtzHandler) Register(r *gin.RouterGroup) {
	r.GET("/common/:deviceId/:channelId", h.Common)
	r.GET("/ptz/:deviceId/:channelId", h.PTZ)
	r.GET("/fi/iris/:deviceId/:channelId", h.Iris)
	r.GET("/fi/focus/:deviceId/:channelId", h.Focus)
	r.GET("/preset/query/:deviceId/:channelId", h.PresetQuery)
	r.GET("/preset/add/:deviceId/:channelId", h.PresetAdd)
	r.GET("/preset/call/:deviceId/:channelId", h.PresetCall)
	r.GET("/preset/delete/:deviceId/:channelId", h.PresetDelete)
	r.GET("/cruise/point/add/:deviceId/:channelId", h.CruisePointAdd)
	r.GET("/cruise/point/delete/:deviceId/:channelId", h.CruisePointDelete)
	r.GET("/cruise/speed/:deviceId/:channelId", h.CruiseSpeed)
	r.GET("/cruise/time/:deviceId/:channelId", h.CruiseTime)
	r.GET("/cruise/start/:deviceId/:channelId", h.CruiseStart)
	r.GET("/cruise/stop/:deviceId/:channelId", h.CruiseStop)
	r.GET("/scan/start/:deviceId/:channelId", h.ScanStart)
	r.GET("/scan/stop/:deviceId/:channelId", h.ScanStop)
	r.GET("/scan/set/left/:deviceId/:channelId", h.ScanSetLeft)
	r.GET("/scan/set/right/:deviceId/:channelId", h.ScanSetRight)
	r.GET("/scan/set/speed/:deviceId/:channelId", h.ScanSetSpeed)
	r.GET("/wiper/:deviceId/:channelId", h.Wiper)
	r.GET("/auxiliary/:deviceId/:channelId", h.Auxiliary)
}

func (h *PtzHandler) Common(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	cmdCode, _ := strconv.Atoi(c.DefaultQuery("cmdCode", "0"))
	speed, _ := strconv.Atoi(c.DefaultQuery("speed", "0"))
	if err := h.svcs.PTZ.FrontEndCommand(deviceID, channelID, cmdCode, speed, 0, 0); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PtzHandler) PTZ(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	cmdCode, _ := strconv.Atoi(c.DefaultQuery("leftRight", "0"))
	speed, _ := strconv.Atoi(c.DefaultQuery("speed", "0"))
	if err := h.svcs.PTZ.PTZControl(deviceID, channelID, cmdCode, speed); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PtzHandler) Iris(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	speed, _ := strconv.Atoi(c.DefaultQuery("speed", "0"))
	if err := h.svcs.PTZ.FrontEndCommand(deviceID, channelID, 0, speed, 0, 0); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PtzHandler) Focus(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	speed, _ := strconv.Atoi(c.DefaultQuery("speed", "0"))
	if err := h.svcs.PTZ.FrontEndCommand(deviceID, channelID, 0, speed, 0, 0); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PtzHandler) PresetQuery(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	if err := h.svcs.PTZ.QueryPresets(deviceID, channelID); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PtzHandler) PresetAdd(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	presetID, _ := strconv.Atoi(c.DefaultQuery("presetId", "0"))
	if err := h.svcs.PTZ.PresetCommand(deviceID, channelID, 1, presetID, 0); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PtzHandler) PresetCall(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	presetID, _ := strconv.Atoi(c.DefaultQuery("presetId", "0"))
	if err := h.svcs.PTZ.PresetCommand(deviceID, channelID, 2, presetID, 0); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PtzHandler) PresetDelete(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	presetID, _ := strconv.Atoi(c.DefaultQuery("presetId", "0"))
	if err := h.svcs.PTZ.PresetCommand(deviceID, channelID, 3, presetID, 0); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PtzHandler) CruisePointAdd(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	routeID, _ := strconv.Atoi(c.DefaultQuery("routeId", "0"))
	pointID, _ := strconv.Atoi(c.DefaultQuery("pointId", "0"))
	if err := h.svcs.PTZ.CruiseCommand(deviceID, channelID, 1, routeID, pointID); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PtzHandler) CruisePointDelete(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	routeID, _ := strconv.Atoi(c.DefaultQuery("routeId", "0"))
	pointID, _ := strconv.Atoi(c.DefaultQuery("pointId", "0"))
	if err := h.svcs.PTZ.CruiseCommand(deviceID, channelID, 2, routeID, pointID); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PtzHandler) CruiseSpeed(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	routeID, _ := strconv.Atoi(c.DefaultQuery("routeId", "0"))
	speed, _ := strconv.Atoi(c.DefaultQuery("speed", "0"))
	if err := h.svcs.PTZ.CruiseCommand(deviceID, channelID, 3, routeID, speed); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PtzHandler) CruiseTime(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	routeID, _ := strconv.Atoi(c.DefaultQuery("routeId", "0"))
	t, _ := strconv.Atoi(c.DefaultQuery("time", "0"))
	if err := h.svcs.PTZ.CruiseCommand(deviceID, channelID, 4, routeID, t); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PtzHandler) CruiseStart(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	routeID, _ := strconv.Atoi(c.DefaultQuery("routeId", "0"))
	if err := h.svcs.PTZ.CruiseCommand(deviceID, channelID, 5, routeID, 0); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PtzHandler) CruiseStop(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	if err := h.svcs.PTZ.CruiseCommand(deviceID, channelID, 6, 0, 0); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PtzHandler) ScanStart(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	speed, _ := strconv.Atoi(c.DefaultQuery("speed", "0"))
	if err := h.svcs.PTZ.ScanCommand(deviceID, channelID, 1, speed); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PtzHandler) ScanStop(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	if err := h.svcs.PTZ.ScanCommand(deviceID, channelID, 2, 0); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PtzHandler) ScanSetLeft(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	param, _ := strconv.Atoi(c.DefaultQuery("param", "0"))
	if err := h.svcs.PTZ.ScanCommand(deviceID, channelID, 3, param); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PtzHandler) ScanSetRight(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	param, _ := strconv.Atoi(c.DefaultQuery("param", "0"))
	if err := h.svcs.PTZ.ScanCommand(deviceID, channelID, 4, param); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PtzHandler) ScanSetSpeed(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	speed, _ := strconv.Atoi(c.DefaultQuery("speed", "0"))
	if err := h.svcs.PTZ.ScanCommand(deviceID, channelID, 5, speed); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PtzHandler) Wiper(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	on := c.DefaultQuery("on", "false") == "true"
	if err := h.svcs.PTZ.Wiper(deviceID, channelID, on); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}

func (h *PtzHandler) Auxiliary(c *gin.Context) {
	deviceID := c.Param("deviceId")
	channelID := c.Param("channelId")
	switchID, _ := strconv.Atoi(c.DefaultQuery("switchId", "0"))
	on := c.DefaultQuery("on", "false") == "true"
	if err := h.svcs.PTZ.Auxiliary(deviceID, channelID, switchID, on); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}
	c.JSON(200, utils.SuccessNoData())
}
