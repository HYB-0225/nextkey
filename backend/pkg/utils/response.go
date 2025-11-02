package utils

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nextkey/nextkey/backend/internal/crypto"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type EncryptedResponse struct {
	Nonce string `json:"nonce"`
	Data  string `json:"data"`
}

type InternalResponse struct {
	Nonce     string      `json:"nonce"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(200, Response{
		Code:    code,
		Message: message,
	})
}

func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func EncryptedSuccess(c *gin.Context, data interface{}) {
	resp := Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}

	// 从上下文获取nonce
	nonce, _ := c.Get("request_nonce")
	nonceStr, _ := nonce.(string)

	// 包装内层数据，嵌入nonce和服务器时间戳
	internalResp := InternalResponse{
		Nonce:     nonceStr,
		Timestamp: time.Now().Unix(),
		Data:      resp,
	}

	jsonData, err := json.Marshal(internalResp)
	if err != nil {
		Error(c, 500, "序列化失败")
		return
	}

	// 从上下文获取加密器
	encryptorVal, exists := c.Get("encryptor")
	if !exists {
		Error(c, 500, "加密器未初始化")
		return
	}

	encryptor, ok := encryptorVal.(crypto.Encryptor)
	if !ok {
		Error(c, 500, "加密器类型错误")
		return
	}

	encryptedData, err := encryptor.Encrypt(string(jsonData))
	if err != nil {
		Error(c, 500, "加密失败")
		return
	}

	c.JSON(200, EncryptedResponse{
		Nonce: nonceStr,
		Data:  encryptedData,
	})
}

func EncryptedError(c *gin.Context, code int, message string) {
	resp := Response{
		Code:    code,
		Message: message,
	}

	// 从上下文获取nonce
	nonce, _ := c.Get("request_nonce")
	nonceStr, _ := nonce.(string)

	// 包装内层数据，嵌入nonce和服务器时间戳
	internalResp := InternalResponse{
		Nonce:     nonceStr,
		Timestamp: time.Now().Unix(),
		Data:      resp,
	}

	jsonData, err := json.Marshal(internalResp)
	if err != nil {
		Error(c, 500, "序列化失败")
		return
	}

	// 从上下文获取加密器
	encryptorVal, exists := c.Get("encryptor")
	if !exists {
		Error(c, 500, "加密器未初始化")
		return
	}

	encryptor, ok := encryptorVal.(crypto.Encryptor)
	if !ok {
		Error(c, 500, "加密器类型错误")
		return
	}

	encryptedData, err := encryptor.Encrypt(string(jsonData))
	if err != nil {
		Error(c, 500, "加密失败")
		return
	}

	c.JSON(200, EncryptedResponse{
		Nonce: nonceStr,
		Data:  encryptedData,
	})
}
