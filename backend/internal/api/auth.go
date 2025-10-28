package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nextkey/nextkey/backend/internal/middleware"
	"github.com/nextkey/nextkey/backend/internal/service"
	"github.com/nextkey/nextkey/backend/pkg/utils"
)

func CardLogin(c *gin.Context) {
	var req service.LoginRequest
	if err := middleware.GetDecryptedData(c, &req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	if req.IP == "" {
		req.IP = c.ClientIP()
	}

	authSvc := service.NewAuthService()
	resp, err := authSvc.CardLogin(&req)
	if err != nil {
		utils.Error(c, 401, err.Error())
		return
	}

	utils.Success(c, resp)
}

func AdminLogin(c *gin.Context) {
	var req service.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	authSvc := service.NewAuthService()
	resp, err := authSvc.AdminLogin(&req)
	if err != nil {
		utils.Error(c, 401, err.Error())
		return
	}

	utils.Success(c, resp)
}

func Heartbeat(c *gin.Context) {
	cardID, _ := c.Get("card_id")
	projectID, _ := c.Get("project_id")

	cardSvc := service.NewCardService()
	if err := cardSvc.Heartbeat(cardID.(uint), projectID.(uint)); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "心跳成功"})
}
