package handler

import (
	"fmt"
	"time"

	"wvp-pro-go/internal/config"
	"wvp-pro-go/internal/service"
	"wvp-pro-go/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	svcs *service.Services
}

func NewAuthHandler(svcs *service.Services) *AuthHandler {
	return &AuthHandler{svcs: svcs}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(200, utils.Fail(400, "参数错误"))
		return
	}

	user, err := h.svcs.User.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(200, utils.Fail(401, err.Error()))
		return
	}

	// Generate JWT token
	token, err := h.generateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(200, utils.Fail(500, "生成Token失败"))
		return
	}

	c.JSON(200, utils.Success(gin.H{
		"accessToken": token,
		"username":    user.Username,
		"serverId":    config.GlobalConfig.SIP.ID,
	}))
}

func (h *AuthHandler) generateToken(userID uint, username string) (string, error) {
	cfg := config.GlobalConfig.JWT
	now := time.Now()
	claims := jwt.MapClaims{
		"userId":   userID,
		"username": username,
		"iat":      now.Unix(),
		"exp":      now.Add(time.Duration(cfg.Expire) * time.Hour).Unix(),
		"iss":      cfg.Issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(200, utils.SuccessNoData())
}

func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(200, utils.Fail(401, "未登录"))
		return
	}

	user, err := h.svcs.User.GetUserByUsername(fmt.Sprintf("%v", username))
	if err != nil {
		c.JSON(200, utils.Fail(404, "用户不存在"))
		return
	}

	c.JSON(200, utils.Success(user))
}

func (h *AuthHandler) AddUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
		RoleID   uint   `json:"roleId" form:"roleId"`
	}
	// Try JSON body first (POST with JSON), fall back to query params
	if err := c.ShouldBindJSON(&req); err != nil {
		if err2 := c.ShouldBindQuery(&req); err2 != nil {
			c.JSON(200, utils.Fail(400, "参数错误"))
			return
		}
	}

	user := &service.UserCreateRequest{
		Username: req.Username,
		Password: req.Password,
	}
	if err := h.svcs.User.CreateUserFromHandler(user); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}

	c.JSON(200, utils.SuccessNoData())
}

func (h *AuthHandler) DeleteUser(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(200, utils.Fail(400, "参数错误"))
		return
	}

	var uid uint
	fmt.Sscanf(id, "%d", &uid)
	if err := h.svcs.User.DeleteUser(uid); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}

	c.JSON(200, utils.SuccessNoData())
}

func (h *AuthHandler) QueryUsers(c *gin.Context) {
	page := 1
	count := 10
	fmt.Sscanf(c.Query("page"), "%d", &page)
	fmt.Sscanf(c.Query("count"), "%d", &count)

	users, err := h.svcs.User.GetUsers(page, count)
	if err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}

	c.JSON(200, utils.Success(users))
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"oldPassword" form:"oldPassword"`
		Password    string `json:"password" form:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		if err2 := c.ShouldBindQuery(&req); err2 != nil {
			c.JSON(200, utils.Fail(400, "参数错误"))
			return
		}
	}

	username, _ := c.Get("username")
	user, err := h.svcs.User.GetUserByUsername(fmt.Sprintf("%v", username))
	if err != nil {
		c.JSON(200, utils.Fail(404, "用户不存在"))
		return
	}

	if err := h.svcs.User.ChangePassword(user.ID, req.OldPassword, req.Password); err != nil {
		c.JSON(200, utils.Fail(500, err.Error()))
		return
	}

	c.JSON(200, utils.SuccessNoData())
}
