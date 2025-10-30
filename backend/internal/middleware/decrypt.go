package middleware

import (
	"encoding/json"
	"errors"
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

		// 获取项目级加密器
		encryptor, project, err := getEncryptorFromRequest(c, encReq.Data)
		if err != nil {
			utils.Error(c, 401, "认证失败")
			c.Abort()
			return
		}

		plaintext, err := encryptor.Decrypt(encReq.Data)
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
		c.Set("encryptor", encryptor)
		c.Set("project_id", project.ID)
		c.Next()
	}
}

func getEncryptorFromRequest(c *gin.Context, encryptedData string) (crypto.Encryptor, *models.Project, error) {
	// 尝试从token中获取project_id
	token := c.GetHeader("Authorization")
	if token != "" && len(token) > 7 && token[:7] == "Bearer " {
		tokenStr := token[7:]
		var tokenModel models.Token
		if err := database.DB.Where("token = ?", tokenStr).First(&tokenModel).Error; err == nil {
			var project models.Project
			if err := database.DB.Where("id = ?", tokenModel.ProjectID).First(&project).Error; err == nil {
				encryptor, err := crypto.NewEncryptor(project.EncryptionScheme, project.EncryptionKey)
				if err != nil {
					return nil, nil, err
				}
				return encryptor, &project, nil
			}
		}
	}

	// 尝试通过解析加密数据中的project_uuid获取项目
	// 这需要尝试所有项目的密钥（用于登录等首次请求）
	var projects []models.Project
	if err := database.DB.Find(&projects).Error; err != nil {
		return nil, nil, err
	}

	for _, project := range projects {
		encryptor, err := crypto.NewEncryptor(project.EncryptionScheme, project.EncryptionKey)
		if err != nil {
			continue
		}

		plaintext, err := encryptor.Decrypt(encryptedData)
		if err != nil {
			continue
		}

		// 尝试解析以获取project_uuid
		var tempData map[string]interface{}
		if err := json.Unmarshal([]byte(plaintext), &tempData); err != nil {
			continue
		}

		data, ok := tempData["data"].(map[string]interface{})
		if !ok {
			continue
		}

		projectUUID, ok := data["project_uuid"].(string)
		if ok && projectUUID == project.UUID {
			return encryptor, &project, nil
		}
	}

	return nil, nil, errors.New("无法找到匹配的项目加密配置")
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
