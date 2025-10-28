package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nextkey/nextkey/backend/internal/database"
	"github.com/nextkey/nextkey/backend/internal/models"
	"github.com/nextkey/nextkey/backend/pkg/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(c, 401, "未提供认证信息")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Error(c, 401, "无效的认证格式")
			c.Abort()
			return
		}

		tokenStr := parts[1]

		var token models.Token
		if err := database.DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
			utils.Error(c, 401, "无效的Token")
			c.Abort()
			return
		}

		if token.IsExpired() {
			utils.Error(c, 401, "Token已过期")
			c.Abort()
			return
		}

		c.Set("token", &token)
		c.Set("card_id", token.CardID)
		c.Set("project_id", token.ProjectID)
		c.Next()
	}
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(c, 401, "未提供认证信息")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Error(c, 401, "无效的认证格式")
			c.Abort()
			return
		}

		tokenStr := parts[1]
		c.Set("admin_token", tokenStr)
		c.Next()
	}
}
