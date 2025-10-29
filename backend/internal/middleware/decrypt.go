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

type InternalRequest struct {
	Nonce     string          `json:"nonce"`
	Timestamp int64           `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
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

		// 使用数据库唯一约束防止重放攻击和竞态条件
		nonce := models.Nonce{Nonce: encReq.Nonce}
		if err := database.DB.Create(&nonce).Error; err != nil {
			// 违反唯一约束表示nonce已存在(重放攻击)
			utils.Error(c, 401, "检测到重放攻击")
			c.Abort()
			return
		}

		plaintext, err := crypto.Decrypt(encReq.Data)
		if err != nil {
			utils.Error(c, 400, "解密失败")
			c.Abort()
			return
		}

		// 解密后验证内层nonce和timestamp
		var internalReq InternalRequest
		if err := json.Unmarshal([]byte(plaintext), &internalReq); err != nil {
			utils.Error(c, 400, "内部数据格式错误")
			c.Abort()
			return
		}

		// 验证内外层nonce一致性
		if internalReq.Nonce != encReq.Nonce {
			utils.Error(c, 401, "Nonce验证失败")
			c.Abort()
			return
		}

		// 验证内外层timestamp一致性
		if internalReq.Timestamp != encReq.Timestamp {
			utils.Error(c, 401, "Timestamp验证失败")
			c.Abort()
			return
		}

		c.Set("decrypted_data", string(internalReq.Data))
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
