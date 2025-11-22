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

// abortUnauthorized returns a 401 HTTP status so the frontend can trigger refresh flows.
func abortUnauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(401, utils.Response{
		Code:    401,
		Message: message,
	})
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
			abortUnauthorized(c, "未提供认证信息")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			abortUnauthorized(c, "无效的认证格式")
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
			abortUnauthorized(c, "无效的Token")
			return
		}

		if !token.Valid {
			abortUnauthorized(c, "Token已失效")
			return
		}

		// 提取claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			abortUnauthorized(c, "无效的Token声明")
			return
		}

		// 验证必要的claims
		adminID, ok := claims["admin_id"]
		if !ok {
			abortUnauthorized(c, "Token缺少必要信息")
			return
		}

		// 提取JTI
		jti, ok := claims["jti"]
		if !ok {
			abortUnauthorized(c, "Token缺少JTI")
			return
		}

		// 检查JTI是否在黑名单中
		var blacklist models.AdminTokenBlacklist
		if err := database.DB.Where("jti = ?", jti).First(&blacklist).Error; err == nil {
			abortUnauthorized(c, "Token已被撤销")
			return
		}

		// 验证管理员是否仍然存在
		var admin models.Admin
		if err := database.DB.First(&admin, uint(adminID.(float64))).Error; err != nil {
			abortUnauthorized(c, "管理员不存在")
			return
		}

		c.Set("admin_id", adminID)
		c.Set("admin_token", tokenStr)
		c.Set("jti", jti)
		c.Next()
	}
}
