package middleware

import (
	"encoding/json"
	"io"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nextkey/nextkey/backend/internal/crypto"
	"github.com/nextkey/nextkey/backend/internal/database"
	"github.com/nextkey/nextkey/backend/internal/models"
	"github.com/nextkey/nextkey/backend/pkg/utils"
)

type EncryptedRequest struct {
	Timestamp int64  `json:"timestamp"`
	Nonce     string `json:"nonce"`
	Data      string `json:"data"`
}

func DecryptMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			utils.Error(c, 400, "读取请求失败")
			c.Abort()
			return
		}

		var encReq EncryptedRequest
		if err := json.Unmarshal(body, &encReq); err != nil {
			utils.Error(c, 400, "无效的请求格式")
			c.Abort()
			return
		}

		now := time.Now().Unix()
		timeDiff := math.Abs(float64(now - encReq.Timestamp))
		if timeDiff > 300 {
			utils.Error(c, 401, "请求已过期")
			c.Abort()
			return
		}

		var existingNonce models.Nonce
		if err := database.DB.Where("nonce = ?", encReq.Nonce).First(&existingNonce).Error; err == nil {
			utils.Error(c, 401, "检测到重放攻击")
			c.Abort()
			return
		}

		database.DB.Create(&models.Nonce{Nonce: encReq.Nonce})

		plaintext, err := crypto.Decrypt(encReq.Data)
		if err != nil {
			utils.Error(c, 400, "解密失败")
			c.Abort()
			return
		}

		c.Set("decrypted_data", plaintext)
		c.Set("request_nonce", encReq.Nonce)
		c.Next()
	}
}

func GetDecryptedData(c *gin.Context, v interface{}) error {
	data, exists := c.Get("decrypted_data")
	if !exists {
		return nil
	}

	plaintext, ok := data.(string)
	if !ok {
		return nil
	}

	return json.Unmarshal([]byte(plaintext), v)
}
