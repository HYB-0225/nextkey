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
		utils.EncryptedError(c, 400, "参数错误")
		return
	}

	if req.IP == "" {
		req.IP = c.ClientIP()
	}

	authSvc := service.NewAuthService()
	resp, err := authSvc.CardLogin(&req)
	if err != nil {
		utils.EncryptedError(c, 401, err.Error())
		return
	}

	utils.EncryptedSuccess(c, resp)
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
	cardID, exists := c.Get("card_id")
	projectID, _ := c.Get("project_id")

	// 免费模式下没有card_id,直接返回成功
	if !exists {
		utils.EncryptedSuccess(c, gin.H{"message": "心跳成功"})
		return
	}

	cardSvc := service.NewCardService()
	if err := cardSvc.Heartbeat(cardID.(uint), projectID.(uint)); err != nil {
		utils.EncryptedError(c, 500, err.Error())
		return
	}

	utils.EncryptedSuccess(c, gin.H{"message": "心跳成功"})
}
