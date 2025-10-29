package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nextkey/nextkey/backend/internal/database"
	"github.com/nextkey/nextkey/backend/internal/models"
	"github.com/nextkey/nextkey/backend/pkg/utils"
)

var jwtSecret []byte

func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

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
		if token.CardID != nil {
			c.Set("card_id", *token.CardID)
		}
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

		// 验证JWT token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("无效的签名方法")
			}
			return jwtSecret, nil
		})

		if err != nil {
			utils.Error(c, 401, "无效的Token")
			c.Abort()
			return
		}

		if !token.Valid {
			utils.Error(c, 401, "Token已失效")
			c.Abort()
			return
		}

		// 提取claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.Error(c, 401, "无效的Token声明")
			c.Abort()
			return
		}

		// 验证必要的claims
		adminID, ok := claims["admin_id"]
		if !ok {
			utils.Error(c, 401, "Token缺少必要信息")
			c.Abort()
			return
		}

		c.Set("admin_id", adminID)
		c.Set("admin_token", tokenStr)
		c.Next()
	}
}
